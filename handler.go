package main

import (
	"image/png"
	"net/http"
	"strconv"
	"time"
)

func registerHandlers() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/wallpaper", wallpaperHandler)
}

// ===== HANDLERS =====

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func wallpaperHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// ---- Device (manual or fallback) ----
	deviceKey := q.Get("device")
	device, ok := Devices[deviceKey]
	if !ok {
		device = Devices["iphone-15"]
	}

	// ---- Language ----
	lang := q.Get("lang")
	if lang == "" {
		lang = "en"
	}

	// ---- Timezone ----
	tzOffset, _ := strconv.Atoi(q.Get("timezone"))
	loc := time.FixedZone("user", tzOffset*3600)
	now := time.Now().In(loc)

	// ---- Weekends mode ----
	weekends := q.Get("weekends")
	if weekends == "" {
		weekends = "off"
	}

	// ---- Render ----
	img := RenderCalendar(
		now,
		device,
		IOSTheme(),
		"months",
		lang,
		weekends,
	)

	w.Header().Set("Content-Type", "image/png")
	_ = png.Encode(w, img)
}
