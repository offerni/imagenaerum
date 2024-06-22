package main

import (
	"github.com/offerni/imagenaerum/img"
)

func main() {
	if err := img.Blur("./files/raw/imgtest.jpg", 2); err != nil {
		panic(err)
	}
}
