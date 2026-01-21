package main

import (
	"log"
	"net/http"
)

func main() {
	RegisterHandlers()
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
