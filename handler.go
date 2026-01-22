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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func wallpaperHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// 1. Device
	deviceKey := q.Get("device")
	device, ok := Devices[deviceKey]
	if !ok {
		device = Devices["iphone-15"]
	}

	// 2. Language
	lang := q.Get("lang")
	if lang == "" {
		lang = "en"
	}

	// 3. Weekends
	weekends := q.Get("weekends")
	if weekends == "" {
		weekends = "off"
	}

	// 4. Timezone
	tz, _ := strconv.Atoi(q.Get("timezone"))
	loc := time.FixedZone("user", tz*3600)
	now := time.Now().In(loc)

	// 5. Render
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
