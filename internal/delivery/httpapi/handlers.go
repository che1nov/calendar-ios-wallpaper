package httpapi

import (
	"image/png"
	"net/http"
	"strconv"

	"calendar-wallpaper/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service usecase.Service
}

func RegisterHandlers(router chi.Router, h Handler) {
	router.Get("/", h.indexHandler)
	router.Get("/wallpaper", h.wallpaperHandler)
	router.Handle("/images/*",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("web/images")),
		),
	)
}

func (h Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func (h Handler) wallpaperHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	tz, _ := strconv.Atoi(q.Get("timezone"))
	size, _ := strconv.Atoi(q.Get("size"))

	params := usecase.RenderParams{
		DeviceKey:   q.Get("device"),
		Lang:        q.Get("lang"),
		Weekends:    q.Get("weekends"),
		DayStyle:    q.Get("style"),
		Timezone:    tz,
		SizePercent: size,
		BgStyle:     q.Get("bg"),
		BgColor:     q.Get("color"),
	}

	img, err := h.Service.RenderWallpaper(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	_ = png.Encode(w, img)
}
