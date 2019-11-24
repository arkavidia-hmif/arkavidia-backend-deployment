package main

import (
	"log"
	"net/http"

	"github.com/arkavidia-hmif/deployment/handler"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
