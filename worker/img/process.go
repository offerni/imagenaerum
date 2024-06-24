package img

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/offerni/imagenaerum/worker/rabbitmq"
	"github.com/offerni/imagenaerum/worker/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (svc *Service) Process() error {
	msgsChannel := make(chan amqp.Delivery)
	errorsChannel := make(chan error)

	err := svc.RabbitmqService.Consume(rabbitmq.ConsumeOpts{
		Ch:           svc.RabbitmqService.Channel,
		Consumer:     "img-consumer",
		ExchangeName: "file_exchange",
		Out:          msgsChannel,
		QueueName:    "files_to_process",
		RoutingKey:   "to_process",
	})
	if err != nil {
		return err
	}

	defer svc.RabbitmqService.Close()

	var wg sync.WaitGroup

	for msg := range msgsChannel {
		wg.Add(1)
		go func(m amqp.Delivery) {
			defer wg.Done()
			m.Ack(false)

			var resp ProcessResponse
			if err := json.Unmarshal(m.Body, &resp); err != nil {
				fmt.Printf("error unmarshalling %s", err.Error())
				errorsChannel <- err
				return
			}

			filePath := filepath.Join(utils.RawPath, resp.File)
			if err := svc.processFile(processFileOpts{
				Blur:       resp.Blur,
				CropAnchor: resp.CropAnchor,
				File:       filePath,
				Grayscale:  resp.Grayscale,
				Invert:     resp.Invert,
				Resize:     resp.Resize,
			}); err != nil {
				fmt.Printf("error processing %s", err.Error())
				errorsChannel <- err
				return
			}
		}(msg)
	}
	wg.Wait()
	close(errorsChannel)

	// check for errors in the channel
	for err := range errorsChannel {
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
	}

	return nil
}

func (svc Service) processFile(opts processFileOpts) error {
	src, err := imaging.Open(opts.File)
	if err != nil {
		return err
	}

	fileExtension := filepath.Ext(opts.File)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	fileProcessedPath := fmt.Sprintf("%s/%s", utils.ProcessedPath, fileName)

	if opts.Blur != nil {
		blurValue, err := strconv.ParseFloat(*opts.Blur, 64)
		if err != nil {
			return err
		}
		src = imaging.Blur(src, blurValue)
	}

	if opts.CropAnchor != nil {
		convertedValues, err := utils.ConvertToIntSlice(*opts.CropAnchor, ",")
		if err != nil {
			return err
		}
		width := convertedValues[0]
		height := convertedValues[1]

		src = imaging.CropAnchor(src, width, height, imaging.Center)
	}

	if opts.Resize != nil {
		convertedValues, err := utils.ConvertToIntSlice(*opts.Resize, ",")
		if err != nil {
			return err
		}

		width := convertedValues[0]
		height := convertedValues[1]

		src = imaging.Resize(src, width, height, imaging.Lanczos)
	}

	if opts.Grayscale != nil {
		src = imaging.Grayscale(src)
	}

	if opts.Invert != nil {
		src = imaging.Invert(src)
	}

	if err := imaging.Save(src, fileProcessedPath); err != nil {
		return err
	}

	// removing files at the end of the processing
	if err := os.Remove(opts.File); err != nil {
		return err
	}

	req, err := json.Marshal(ProcessedRequest{
		File: fileName,
	})
	if err != nil {
		return err
	}

	var mu sync.Mutex

	// new channel to prevent threads sharing the same channel on publish
	ch, err := svc.RabbitmqService.NewChannel()
	if err != nil {
		return fmt.Errorf("error creating channel %v", err)
	}

	defer ch.Close()

	mu.Lock()
	svc.RabbitmqService.Publish(rabbitmq.PublishOpts{
		Body:         req,
		Ch:           ch,
		ExchangeName: "file_exchange",
		QueueName:    "processed_files",
		RoutingKey:   "processed",
	})
	mu.Unlock()

	return nil
}

type processFileOpts struct {
	Blur       *string
	CropAnchor *string
	File       string
	Grayscale  *string
	Invert     *string
	Resize     *string
}

type ProcessResponse struct {
	Blur       *string `json:"blur"`
	CropAnchor *string `json:"crop_anchor"`
	File       string  `json:"file"`
	Grayscale  *string `json:"grayscale"`
	Invert     *string `json:"invert"`
	Resize     *string `json:"resize"`
}

type ProcessedRequest struct {
	File string `json:"file"`
}
