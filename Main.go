package main

import (
	"fmt"
	"log"
	"net/http"
)

var version string = "DEVELOPMENT"

func main()  {
	http.HandleFunc("/gadget", ServeGadget)
	http.HandleFunc("/user", UserInfo)
	http.HandleFunc("/download", DownloadURLCollector)

	fmt.Println("patreon-dl - Patreon Image Downloader | release: " + version)
	fmt.Println("Send features/bugreports/suggestions to https://github.com/PrivateGER/patreon-dl/issues")
	fmt.Println("Please open https://patreon.com/creatorname/posts, open the developer console (F12), paste the following into the console and run it with ENTER:")
	fmt.Println(`(async()=>{eval(await(await fetch("http://localhost:9849/gadget")).text());})();`)

	fmt.Println("\nHint: highlight the text and use CTRL+SHIFT+C to copy it.\nWaiting for the browser to send data on port 9849...")
	err := http.ListenAndServe(":9849", nil)
	if err != nil {
		log.Println(err)
	}
}
