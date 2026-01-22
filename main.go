package main

import (
	"fmt"
	"net/http"
)

func main() {
	RegisterHandlers()
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
