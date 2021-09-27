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

func DownloadJobHandler(safeDownloadList *SafeDownloadList) {
	safeDownloadList.mu.Lock()
	downloadQueue = make(chan []string, len(safeDownloadList.list))
	downloadResults = make(chan string, len(safeDownloadList.list))
	safeDownloadList.mu.Unlock()

	filepath := filepath.Join(".", downloadDir)
	_ = os.MkdirAll(filepath, os.ModePerm)

	for w := 1; w <= 5; w++ {
		go DownloadWorker(w, downloadQueue, downloadResults)
	}

	safeDownloadList.mu.Lock()
	startedAt := time.Now()
	fmt.Println("Queuing", strconv.Itoa(cap(downloadQueue)), "downloads...")
	for _, dlURL := range safeDownloadList.list {
		downloadQueue <- dlURL
	}
	close(downloadQueue)
	safeDownloadList.mu.Unlock()

	for index := 1; index <= cap(downloadQueue); index++ {
		fmt.Println(<-downloadResults)
	}
	endedAt := time.Now()

	fmt.Printf(
		"Finished downloading in %d second(s)!\nYou can download another set of images using the below JS snippet or exit patreon-dl:\n(async()=>{eval(await(await fetch(\"http://localhost:9849/gadget\")).text());})();",
		int(endedAt.Sub(startedAt).Seconds()),
	)

	return
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
			_ = resp.Body.Close()
			results <- err.Error()
			continue
		}

		// Create the file
		out, err := os.Create(filepath.Join(".", downloadDir, path.Base(j[1])))
		if err != nil {
			_ = resp.Body.Close()
			results <- err.Error()
			continue
		}

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			_ = out.Close()
			_ = resp.Body.Close()
			results <- err.Error()
			continue
		}

		results <- path.Base(j[1]) + " downloaded"
		_ = resp.Body.Close()
		_ = out.Close()
	}
}
