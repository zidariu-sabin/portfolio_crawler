package utils

import (
	"os"
)

// creates a file repos.yaml and writes the output byte slice to it
func WriteJson(fileName string, reposJsonOutput []byte) error {
	f, err := os.Create(fileName)

	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(reposJsonOutput)

	if err != nil {
		return err
	}

	return nil
}

// creates a file repos.yaml and writes the output byte slice to it
func WriteYaml(fileName string, output []byte) error {
	f, err := os.Create(fileName)

	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(output)

	if err != nil {
		return err
	}

	return nil
}
