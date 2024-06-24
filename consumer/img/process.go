package img

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/offerni/imagenaerum/consumer/rabbitmq"
	"github.com/offerni/imagenaerum/worker/utils"
)

type ProcessOpts struct {
	Files  []*multipart.FileHeader
	Params string
}

const maxSizeMB = 10
const maxSizeBytes = maxSizeMB * 1024 * 1024

func (svc Service) Process(opts ProcessOpts) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	var wg sync.WaitGroup

	var mu sync.Mutex

	errCh := make(chan error, len(opts.Files)) // Buffer size to avoid blocking

	for _, file := range opts.Files {
		wg.Add(1)
		go func(f *multipart.FileHeader) {
			defer wg.Done()

			processedFile, err := svc.processAndStoreFile(f)
			if err != nil {
				errCh <- fmt.Errorf("error svc.processAndStoreFile %v", err)
				return
			}

			req, err := json.Marshal(ProcessRequest{
				File: processedFile.Filename,
				Params: map[string]string{
					"sigma": opts.Params,
				},
			})
			if err != nil {
				errCh <- fmt.Errorf("error json.Marshal(ProcessRequest %v", err)
				return
			}

			ch, err := svc.RabbitMQService.NewChannel()
			if err != nil {
				errCh <- fmt.Errorf("error creating channel %v", err)
				return
			}
			defer func() {
				if err := ch.Close(); err != nil {
					log.Printf("error closing channel: %v", err)
				}
			}()

			// adding mutex here publishing is not thread-safe
			mu.Lock()
			err = svc.RabbitMQService.Publish(rabbitmq.PublishOpts{
				Ch:           ch,
				QueueName:    "files_to_process",
				ExchangeName: "file_exchange",
				RoutingKey:   "to_process",
				Body:         req,
			})
			if err != nil {
				errCh <- fmt.Errorf("error publishing message: %v", err)
				return
			}
			mu.Unlock()
		}(file)
	}

	// Wait until all threads are done and close the error channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Check for errors in the channel
	for err := range errCh {
		if err != nil {
			log.Printf("Errors Processing file: %s", err.Error())
		}
	}

	return nil
}

func (svc Service) processAndStoreFile(file *multipart.FileHeader) (*multipart.FileHeader, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// renaming file to uuid
	file.Filename = fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))

	fileRawPath := fmt.Sprintf("%s/%s", utils.RawPath, file.Filename)
	dst, err := os.Create(fileRawPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, f); err != nil {
		return nil, err
	}

	return file, nil
}

func (opts ProcessOpts) Validate() error {
	if len(opts.Files) == 0 {
		return ErrNoFiles
	}

	if opts.Params == "" {
		return ErrNoParam
	}

	for _, file := range opts.Files {
		if file.Size > maxSizeBytes {
			return ErrInvalidSize
		}
	}

	return nil
}

type ProcessRequest struct {
	File   string            `json:"file"`
	Params map[string]string `json:"params"`
}
