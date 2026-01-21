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

	deviceKey := q.Get("device")
	if deviceKey == "" {
		deviceKey = "iphone-15"
	}

	device, ok := Devices[deviceKey]
	if !ok {
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

	tzOffset, _ := strconv.Atoi(q.Get("timezone"))
	loc := time.FixedZone("user", tzOffset*3600)
	now := time.Now().In(loc)

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
