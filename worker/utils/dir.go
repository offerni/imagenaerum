package utils

import (
	"fmt"
	"log"
	"os"
)

const ConvertedPath string = "../files/converted"
const RawPath string = "../files/raw"

var directories = []string{
	ConvertedPath,
	RawPath,
}

func ensureDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dirName, err)
		}
		log.Printf("Directory %s not found, created.\n", dirName)
	} else if err != nil {
		return fmt.Errorf("error checking directory %s: %w", dirName, err)
	}

	return nil
}

func EnsureDirectories() {
	for _, dir := range directories {
		ensureDir(dir)
	}
}
