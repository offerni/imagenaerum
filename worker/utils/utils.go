package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const ProcessedPath string = "../_files/processed"
const RawPath string = "../_files/raw"

var directories = []string{
	ProcessedPath,
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

func ToPointer[T any](value T) *T {
	return &value
}

func ConvertToIntSlice(input string, separator string) ([]int, error) {
	values := strings.Split(input, separator)
	result := make([]int, len(values))

	for i, val := range values {
		parsedVal, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		result[i] = parsedVal
	}

	return result, nil
}
