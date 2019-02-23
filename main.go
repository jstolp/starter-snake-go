// PofAdder-GO - a batttlesnake.io AI
// defines routes, runs on port 9000
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/start", Start)
	http.HandleFunc("/move", Move)
	http.HandleFunc("/end", End)
	http.HandleFunc("/ping", Ping)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Add filename into logging messages, and MICO SECS
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.Printf("Running jSnake Server on port %s...\n copy&paste http://localhost:%s", port, port)
	http.ListenAndServe(":"+port, LoggingHandler(http.DefaultServeMux))
}
