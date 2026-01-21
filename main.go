package main

import (
	"image/png"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	http.HandleFunc("/wallpaper", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		deviceKey := q.Get("device")
		device, ok := Devices[deviceKey]
		if !ok {
			device = Devices["iphone-15"]
		}

		lang := q.Get("lang")
		if lang == "" {
			lang = "en"
		}

		tz, _ := strconv.Atoi(q.Get("timezone"))
		loc := time.FixedZone("user", tz*3600)
		now := time.Now().In(loc)

		weekends := q.Get("weekends")
		if weekends == "" {
			weekends = "off"
		}

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
	})

	http.ListenAndServe(":8080", nil)
}
