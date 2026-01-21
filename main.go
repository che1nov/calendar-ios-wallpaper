package main

import (
	"net/http"
)

func main() {
	RegisterHandlers()

	http.Handle("/web/",
		http.StripPrefix("/web/",
			http.FileServer(http.Dir("web"))))

	http.ListenAndServe(":8080", nil)
}
