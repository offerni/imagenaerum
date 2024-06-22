package img

import (
	"fmt"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const convertedPath string = "./files/converted"

func (svc *Service) Blur(path string, sigma float64) error {
	src, err := imaging.Open(path)
	if err != nil {
		return err
	}

	img := imaging.Blur(src, sigma)
	fileName := uuid.New().String()
	if err := imaging.Save(img, fmt.Sprintf("%s/%s.jpg", convertedPath, fileName)); err != nil {
		return err
	}

	return nil
}
