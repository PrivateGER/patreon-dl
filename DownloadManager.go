package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var downloadDir string = "images"

var downloadQueue chan []string
var downloadResults chan string

func DownloadJobHandler(downloadList [][]string) {
	downloadQueue = make(chan []string, len(downloadList))
	downloadResults = make(chan string, len(downloadList))

	path := filepath.Join(".", downloadDir)
	_ = os.MkdirAll(path, os.ModePerm)

	for w := 1; w <= 5; w++ {
		go DownloadWorker(w, downloadQueue, downloadResults)
	}

	startedAt := time.Now()
	fmt.Println("Queuing", strconv.Itoa(len(downloadList)), "downloads...")
	for _, dlURL := range downloadList {
		downloadQueue <- dlURL
	}
	close(downloadQueue)

	for index := 1; index <= len(downloadList); index++ {
		fmt.Println(<-downloadResults)
	}
	endedAt := time.Now()

	fmt.Printf("Finished downloading in %d second(s)! Exiting...", int(endedAt.Sub(startedAt).Seconds()))
	os.Exit(0)
}

func DownloadWorker(id int, jobs <-chan []string, results chan<- string) {
	fmt.Println("Download worker " + strconv.Itoa(id) + " started")
	for j := range jobs {
		if j[0] == "" {
			results <- "You do not have access to " + path.Base(j[1] + ", membership tier too low")
			continue
		}

		// Skip Dropbox logo embed
		if strings.HasSuffix(j[0], "content-folder_dropbox-large.png") {
			results <- "Cannot download Dropbox folders, do that manually!"
			continue
		}

		resp, err := http.Get(j[0])
		if err != nil {
			results <- err.Error()
			continue
		}

		// Create the file
		out, err := os.Create(filepath.Join(".", downloadDir, path.Base(j[1])))
		if err != nil {
			results <- err.Error()
			continue
		}

		// Write the body to file
		_, err = io.Copy(out, resp.Body)

		if err != nil {
			results <- err.Error()
		}
		results <- path.Base(j[1]) + " downloaded"
		_ = out.Close()
		_ = resp.Body.Close()
	}
}
