package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type PatreonUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DownloadList [][]string
type SafeDownloadList struct {
	list DownloadList
	mu   sync.Mutex
}


func UserInfo(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var user PatreonUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Couldn't parse /user body.")
		fmt.Println(err)
	}

	fmt.Println("Downloading images from " + user.Name + " with Patreon-ID " + strconv.Itoa(user.ID) + ".\nPlease wait for the download links to be collected...")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = fmt.Fprintf(w, "OK")
	if err != nil {
		log.Println(err)
		return
	}
}

func DownloadURLCollector(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var safeDownloadList SafeDownloadList
	safeDownloadList.mu.Lock()
	defer safeDownloadList.mu.Unlock()

	err = json.Unmarshal(body, &safeDownloadList.list)
	if err != nil {
		fmt.Println("Couldn't parse /download body.")
		fmt.Println(err)
	}

	fmt.Println("Download links received! Starting download...")

	go DownloadJobHandler(&safeDownloadList)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = fmt.Fprintf(w, "OK")
	if err != nil {
		log.Println(err)
		return 
	}
}

//go:embed client.js
var jsGadget string
func ServeGadget(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err := fmt.Fprintf(w, jsGadget)
	if err != nil {
		log.Println(err)
		return 
	}
}
