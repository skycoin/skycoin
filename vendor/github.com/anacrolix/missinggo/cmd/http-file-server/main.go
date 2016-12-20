package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	const addr = ":8080"
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("serving %q at %s", dir, addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(dir))))
}
