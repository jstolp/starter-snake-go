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
	http.HandleFunc("/info", Info)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Add filename into logging messages
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Running jSnake Server on port %s...\n", port)
	http.ListenAndServe(":"+port, LoggingHandler(http.DefaultServeMux))
}
