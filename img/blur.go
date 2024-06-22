package img

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/offerni/imagenaerum/utils"
)

func (svc *Service) Blur(files []*multipart.FileHeader, sigma float64) error {
	var wg sync.WaitGroup

	errCh := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			fmt.Printf("Processing file: %s at: %v\n", file.Filename, time.Now().Format("2006-01-02 15:04:05.000"))
			if err := svc.processFileBlur(file, sigma); err != nil {
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

func (svc Service) processFileBlur(file *multipart.FileHeader, sigma float64) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	fileRawPath := fmt.Sprintf("%s/%s", utils.RawPath, file.Filename)
	dst, err := os.Create(fileRawPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, f); err != nil {
		return err
	}

	src, err := imaging.Open(fileRawPath)
	if err != nil {
		return err
	}

	fileName := uuid.New().String()
	fileCOnvertedPath := fmt.Sprintf("%s/%s.jpg", utils.ConvertedPath, fileName)

	img := imaging.Blur(src, sigma)
	if err := imaging.Save(img, fileCOnvertedPath); err != nil {
		return err
	}

	// removing files at the end of the processing
	if err := os.Remove(fileRawPath); err != nil {
		return err
	}

	return nil
}
