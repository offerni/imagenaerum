package img

import (
	"encoding/json"
	"fmt"
	"io"
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

const maxSizeMB = 5
const maxSizeBytes = maxSizeMB * 1024 * 1024

func (svc Service) Process(opts ProcessOpts) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error)
	for _, file := range opts.Files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			file, err := svc.processAndStoreFile(file)
			if err != nil {
				errCh <- err
			}

			req, err := json.Marshal(ProcessRequest{
				File: file.Filename,
				Params: map[string]string{
					"sigma": opts.Params,
				},
			})
			if err != nil {
				errCh <- err
			}

			err = svc.RabbitMQService.Publish(rabbitmq.PublishOpts{
				Ch:           svc.RabbitMQService.Channel,
				QueueName:    "files_to_process",
				ExchangeName: "file_exchange",
				RoutingKey:   "to_process",
				Body:         req,
			})
			if err != nil {
				errCh <- err
			}

		}(file)
	}

	// Wait until all threads are done and close the error channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// check for errors in the channel
	for err := range errCh {
		if err != nil {
			return err
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
