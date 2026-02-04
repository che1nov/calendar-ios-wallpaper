package main

import (
	"fmt"
	"net/http"
	"time"

	"calendar-wallpaper/internal/delivery/httpapi"
	"calendar-wallpaper/internal/domain"
	"calendar-wallpaper/internal/rendering"
	"calendar-wallpaper/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func main() {
	service := usecase.Service{
		Clock:    usecase.SystemClock{},
		Renderer: rendering.Renderer{},
		Theme:    domain.IOSTheme(),
	}

	router := chi.NewRouter()
	httpapi.RegisterHandlers(router, httpapi.Handler{Service: service})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	fmt.Println("Listening on :8080")
	_ = server.ListenAndServe()
}
