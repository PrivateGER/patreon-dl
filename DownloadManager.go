package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
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

type result struct {
	err error
	path string
}

var downloadQueue chan []string
var downloadResults chan *result

func DownloadJobHandler(safeDownloadList *SafeDownloadList) {
	safeDownloadList.mu.Lock()
	downloadQueue = make(chan []string, len(safeDownloadList.list))
	downloadResults = make(chan *result, len(safeDownloadList.list))
	safeDownloadList.mu.Unlock()

	fp := filepath.Join(".", downloadDir)
	pb := progressbar.Default(int64(cap(downloadQueue)))
	_ = os.MkdirAll(fp, os.ModePerm)

	for w := 1; w <= 5; w++ {
		go DownloadWorker(w, downloadQueue, downloadResults, pb)
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
		res := <-downloadResults
		if res.err != nil {
			fmt.Printf("Got error: %v", res.err)
			continue
		}
		//fmt.Println(res.path) // enable to log each downloaded file
	}
	endedAt := time.Now()

	fmt.Printf(
		"Finished downloading in %d second(s)!\nYou can download another set of images using the below JS snippet or exit patreon-dl:\n(async()=>{eval(await(await fetch(\"http://localhost:9849/gadget\")).text());})();",
		int(endedAt.Sub(startedAt).Seconds()),
	)

	return
}

func DownloadWorker(id int, jobs <-chan []string, resultCh chan<- *result, bar *progressbar.ProgressBar) {
	fmt.Println("Download worker " + strconv.Itoa(id) + " started")
	processJob := func(job, strPath string) *result {
		if job == "" {
			err := fmt.Errorf("you do not have access to %s, membership tier too low", path.Base(strPath))
			return &result{err: err}
		}

		// Skip Dropbox logo embed
		if strings.HasSuffix(job, "content-folder_dropbox-large.png") {
			err := fmt.Errorf( "cannot download Dropbox folders, do that manually")
			return &result{err: err}
		}

		resp, err := http.Get(job)
		if err != nil {
			return &result{err: err}
		}
		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(filepath.Join(".", downloadDir, path.Base(strPath)))
		if err != nil {
			return &result{err: err}
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return &result{err: err}
		}

		bar.Add(1)
		pathStr := path.Base(strPath) + " downloaded"
		return  &result{path: pathStr}
	}

	for j := range jobs {
		resultCh <-processJob(j[0], j[1])
	}
}
