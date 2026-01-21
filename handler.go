package main

import (
	"image/png"
	"net/http"
	"strconv"
	"time"
)

func RegisterHandlers() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/wallpaper", wallpaperHandler)
}

/* ================= INDEX ================= */

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

/* ================= WALLPAPER ================= */

func wallpaperHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// ---- DEVICE ----
	deviceKey := q.Get("device")
	device, ok := Devices[deviceKey]
	if !ok {
		device = detectDeviceFromUA(r.UserAgent())
	}
	if device.Key == "" {
		device = Devices["iphone-15"]
	}

	// ---- LANGUAGE ----
	lang := q.Get("lang")
	if lang == "" {
		lang = "en"
	}

	// ---- WEEKENDS ----
	weekends := q.Get("weekends")
	if weekends == "" {
		weekends = "off"
	}

	// ---- TIMEZONE ----
	tzOffset := 0
	if tzStr := q.Get("timezone"); tzStr != "" {
		if v, err := strconv.Atoi(tzStr); err == nil {
			tzOffset = v
		}
	}

	loc := time.FixedZone("user", tzOffset*3600)
	now := time.Now().In(loc)

	// ---- RENDER ----
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
