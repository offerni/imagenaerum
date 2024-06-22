package img

import "github.com/disintegration/imaging"

const ConvertedPath string = "./files/converted"

func Blur(path string, sigma float64) error {
	src, err := imaging.Open(path)
	if err != nil {
		return err
	}

	img := imaging.Blur(src, sigma)
	if err := imaging.Save(img, "./files/converted/blurred.jpg"); err != nil {
		return err
	}

	return nil
}
