package utils

import (
	"os"
)

// creates a file repos.yaml and writes the output byte slice to it
func WriteJson(reposJsonOutput []byte) error {
	f, err := os.Create("repos.json")

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
func WriteYaml(output []byte) error {
	f, err := os.Create("repos.yaml")

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
