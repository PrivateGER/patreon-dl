package main

import (
	"fmt"
	"net/http"
)

var sessionID string = ""

func main()  {
	http.HandleFunc("/gadget", ServeGadget)
	http.HandleFunc("/user", UserInfo)
	http.HandleFunc("/download", DownloadURLCollector)
	http.HandleFunc("/done", JSFinished)

	fmt.Println("patreon-dl v0.0.1 - Patreon Image Downloader\nPlease open https://patreon.com/creatorname/posts, open the developer console (F12), paste the following into the console and run it with ENTER:")
	fmt.Println(`(async()=>{eval(await(await fetch("http://localhost:9849/gadget")).text());})();`)

	fmt.Println("\nWaiting for the browser to send data on port 9849...")
	err := http.ListenAndServe(":9849", nil)
	if err != nil {
		return
	}
}

