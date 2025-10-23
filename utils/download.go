package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"portfolio_crawler/globals"
	"strings"
	"sync"
	"time"
)

// ValidateDest checks that the destination path's directory exists and is writable.
// If the directory doesn't exist it attempts to create it.
func ValidateDest(dest string) (string, error) {

	finalDest := ""

	// expand ~ to home directory based on OS
	if dest == "~" || strings.HasPrefix(dest, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return finalDest, fmt.Errorf("cannot determine user home directory: %w", err)
		}
		dest = filepath.Join(home, strings.TrimPrefix(dest, "~/"))
	}
	finalDest = filepath.Clean(dest)

	info, err := os.Stat(dest)
	if os.IsNotExist(err) {
		// try to create the directory
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return finalDest, fmt.Errorf("creating directory failed %s: %w", dest, err)
		}
	} else if err != nil {
		return finalDest, fmt.Errorf("stat check failed %s: %w", dest, err)
	} else if !info.IsDir() {
		return finalDest, fmt.Errorf("%s is not a directory", dest)
	}

	// check writability by creating a temp file
	f, err := os.CreateTemp(dest, ".permcheck-*")
	if err != nil {
		return finalDest, fmt.Errorf("directory not writable: %w", err)
	}
	f.Close()
	_ = os.Remove(f.Name())

	return finalDest, nil
}

func CreateMetaData(name string) ([]byte, error) {
	//have to find a way to link the meta data to the specific file details
	//use reflect library to on how to dynamically acess types defiend in the metaData stuct
	//ReposMetaData
}

func DownloadFile(url, newFilePath string, fileName string) error {

	//request timeout after 10 seconds
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "portfolio-crawler")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("downloading %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response %d for %s", resp.StatusCode, url)
	}

	//create the file
	out, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	//
	//should add meta data of the file and create a new one
	//
	return err
}

// DownloadMany downloads multiple files concurrently.
// Retries each failed download up to `retries` times.
func DownloadMany(urls map[string]string, retries int) error {

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // limit concurrency to 5

	for name, url := range urls {
		wg.Add(1)
		dest := globals.DestinationDir + "/" + name + ".md"
		go func(url, dest string) {
			defer wg.Done()
			sem <- struct{}{} // acquire slot
			defer func() { <-sem }()

			for attempt := 1; attempt <= retries; attempt++ {
				err := DownloadFile(url, dest, name)
				if err != nil {
					fmt.Printf("Error: (%d/%d) Failed to download %s: %v. Retrying...\n", attempt, retries, url, err)
					time.Sleep(2 * time.Second)
					continue
				}
				fmt.Printf("Downloaded: %s >> %s\n", url, dest)
				return
			}
			fmt.Printf("Exit: Giving up on %s after %d attempts\n", url, retries)
		}(url, dest)
	}

	wg.Wait()

	return nil
}
