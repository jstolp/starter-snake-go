// PofAdder-GO - a batttlesnake.io AI
// defines routes, runs on port 9000
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
    http.HandleFunc("/", RootInfo)
	http.HandleFunc("/start", Start)
	http.HandleFunc("/move", Move)
	http.HandleFunc("/end", End)
	// http.HandleFunc("/ping", Ping) - Ping is deprecated in version v1 of the bs api
	http.HandleFunc("/favicon.ico", favicon)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Add filename into logging messages, and MICO SECS
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.Printf("Running on port %s...\n http://localhost:%s", port, port)
	http.ListenAndServe(":"+port, LoggingHandler(http.DefaultServeMux))
}
