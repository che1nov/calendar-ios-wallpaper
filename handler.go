package main

import (
	"image/png"
	"net/http"
	"strconv"
	"time"
)

func RegisterHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	http.HandleFunc("/wallpaper", wallpaperHandler)
}

func wallpaperHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	device := Devices[q.Get("device")]
	if device.Key == "" {
		device = Devices["iphone-15"]
	}

	lang := q.Get("lang")
	if lang == "" {
		lang = "en"
	}

	weekends := q.Get("weekends")
	if weekends == "" {
		weekends = "off"
	}

	tz, _ := strconv.Atoi(q.Get("timezone"))
	now := time.Now().In(time.FixedZone("user", tz*3600))

	img := RenderCalendar(now, device, IOSTheme(), "months", lang, weekends)

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
}
