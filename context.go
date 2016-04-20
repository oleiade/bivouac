package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var bivouacFile string = ".bivouac"

func findBivouacFile() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var evaluatedPath string = pwd

	for {
		if evaluatedPath == "/" {
			break
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s", evaluatedPath, bivouacFile)); err == nil {
			return fmt.Sprintf("%s/%s", evaluatedPath, bivouacFile), nil
		}

		evaluatedPath = filepath.Dir(evaluatedPath)
	}

	return "", fmt.Errorf("unable to locate your projects .bivouac file")
}
