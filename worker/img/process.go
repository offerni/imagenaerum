package img

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	var wg sync.WaitGroup
	go func() {
		for msg := range msgsChannel {
			wg.Add(1)
			go func(m amqp.Delivery) {
				defer wg.Done()
				m.Ack(false)

				var resp ProcessResponse
				if err := json.Unmarshal(m.Body, &resp); err != nil {
					errorsChannel <- err
					return
				}

				filePath := filepath.Join(utils.RawPath, resp.File)
				if err := svc.processFile(filePath, 5); err != nil {
					errorsChannel <- err
					return
				}
			}(msg)
		}
		wg.Wait()
		close(errorsChannel)
	}()

	// check for errors in the channel
	for err := range errorsChannel {
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
	}

	return nil
}

func (svc Service) processFile(fileRawPath string, sigma float64) error {

	src, err := imaging.Open(fileRawPath)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.jpg", uuid.New().String())
	fileProcessedPath := fmt.Sprintf("%s/%s", utils.ProcessedPath, fileName)

	img := imaging.Blur(src, sigma)
	if err := imaging.Save(img, fileProcessedPath); err != nil {
		return err
	}

	// removing files at the end of the processing
	if err := os.Remove(fileRawPath); err != nil {
		return err
	}

	req, err := json.Marshal(ProcessedRequest{
		File: fileName,
	})
	if err != nil {
		return err
	}
	svc.RabbitmqService.Publish(rabbitmq.PublishOpts{
		Body:         req,
		Ch:           svc.RabbitmqService.Channel,
		ExchangeName: "file_exchange",
		QueueName:    "processed_files",
		RoutingKey:   "processed",
	})

	return nil
}

type ProcessResponse struct {
	File   string            `json:"file"`
	Params map[string]string `json:"params"`
}

type ProcessedRequest struct {
	File string `json:"file"`
}
