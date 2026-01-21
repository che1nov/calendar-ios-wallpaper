package main

import (
	"fmt"
	"net/http"
)

func main() {
	registerHandlers()

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
