package img

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/offerni/imagenaerum/utils"
)

func (svc *Service) Blur(files []*multipart.FileHeader, sigma float64) error {
	for _, file := range files {
		svc.processFileBlur(file, sigma)
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
