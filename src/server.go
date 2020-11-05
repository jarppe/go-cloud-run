package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetPrefix("go-cloud-run: ")

	log.Printf("Server starting...")

	fileServer := http.FileServer(http.Dir("./assets"))
	http.Handle("/", fileServer)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello!"))
	})

	host, hostSet := os.LookupEnv("HOST")
	if !hostSet {
		host = "0.0.0.0"
	}
	port, portSet := os.LookupEnv("PORT")
	if !portSet {
		port = "8080"
	}

	log.Printf("Server listening at %s:%s", host, port)
	if err := http.ListenAndServe(host + ":" + port, nil); err != nil {
		log.Fatal(err)
	}
}
