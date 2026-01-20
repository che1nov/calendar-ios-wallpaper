package main

import (
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/wallpaper", wallpaperHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func wallpaperHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	mode := q.Get("mode")
	if mode == "" {
		mode = "months"
	}

	themeName := q.Get("theme")
	if themeName == "" {
		themeName = "ios"
	}

	lang := q.Get("lang")
	if lang == "" {
		lang = "en"
	}

	tz, _ := strconv.Atoi(q.Get("timezone"))
	loc := time.FixedZone("user", tz*3600)
	now := time.Now().In(loc)

	theme := Themes[themeName]

	img := RenderCalendarWithLang(now, theme, mode, lang)

	w.Header().Set("Content-Type", "image/png")
	_ = png.Encode(w, img)
}
