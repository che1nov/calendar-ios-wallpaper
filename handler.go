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

	http.Handle("/images/",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("web/images")),
		),
	)

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

	// 4. Day style (НОВОЕ)
	styleParam := q.Get("style")
	dayStyle := DayDots // default

	bg := BackgroundStyle(q.Get("bg"))
	if bg == "" {
		bg = BgIOS
	}

	switch styleParam {
	case "bars":
		dayStyle = DayBars
	case "numbers":
		dayStyle = DayNumbers
	}

	// 5. Timezone
	tz, _ := strconv.Atoi(q.Get("timezone"))
	loc := time.FixedZone("user", tz*3600)
	now := time.Now().In(loc)

	//6. Size
	sizeParam := q.Get("size")
	uiScale := 1.0

	if sizeParam != "" {
		if v, err := strconv.Atoi(sizeParam); err == nil {
			uiScale = float64(v) / 100.0
		}
	}

	bgStyle := BackgroundStyle(q.Get("bg"))
	if bgStyle == "" {
		bgStyle = BgIOS
	}

	bgColor := q.Get("color")
	if bgColor == "" {
		bgColor = "black"
	}

	// 7. Render
	img := RenderCalendar(
		now,
		device,
		IOSTheme(),
		"months",
		lang,
		weekends,
		dayStyle,
		uiScale,
		bgStyle,
		bgColor,
	)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	_ = png.Encode(w, img)

}
