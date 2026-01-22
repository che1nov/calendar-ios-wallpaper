package main

// DeviceProfile описывает геометрию lock screen iPhone.
// Все зоны считаются ОТ ВЫСОТЫ экрана через ratio.
type DeviceProfile struct {
	Key    string
	Name   string
	Width  int
	Height int

	// Проценты высоты экрана (0.0 – 1.0)
	ClockZoneRatio   float64 // где заканчиваются часы / дата
	ButtonsZoneRatio float64 // где начинается зона фонарик / камера
}

// =======================
// SAFE ZONE HELPERS
// =======================

// Нижняя граница зоны часов / даты
func (d DeviceProfile) ClockBottom() int {
	return int(float64(d.Height) * d.ClockZoneRatio)
}

// Верхняя граница зоны кнопок (фонарик / камера)
func (d DeviceProfile) ButtonsTop() int {
	return int(float64(d.Height) * d.ButtonsZoneRatio)
}

// =======================
// ЭТАЛОННЫЕ ПРОФИЛИ
// =======================

var Devices = map[string]DeviceProfile{

	// ===== SE / Home Button =====
	"iphone-se": {
		Key:    "iphone-se",
		Name:   "iPhone SE (2 / 3)",
		Width:  750,
		Height: 1334,

		ClockZoneRatio:   0.26,
		ButtonsZoneRatio: 0.84,
	},

	// ===== XR / 11 =====
	"iphone-11": {
		Key:    "iphone-11",
		Name:   "iPhone 11 / XR",
		Width:  828,
		Height: 1792,

		ClockZoneRatio:   0.28,
		ButtonsZoneRatio: 0.83,
	},

	// ===== 12 / 13 / 14 =====
	"iphone-12": {
		Key:    "iphone-12",
		Name:   "iPhone 12 / 13 / 14",
		Width:  1170,
		Height: 2532,

		ClockZoneRatio:   0.30,
		ButtonsZoneRatio: 0.82,
	},

	// ===== Pro / Dynamic Island =====
	"iphone-15": {
		Key:    "iphone-15",
		Name:   "iPhone 14 Pro / 15 / 15 Pro",
		Width:  1179,
		Height: 2556,

		ClockZoneRatio:   0.30,
		ButtonsZoneRatio: 0.75,
	},

	// ===== Pro Max =====
	"iphone-15-pro-max": {
		Key:    "iphone-15-pro-max",
		Name:   "iPhone Pro Max",
		Width:  1290,
		Height: 2796,

		ClockZoneRatio:   0.32,
		ButtonsZoneRatio: 0.80,
	},
}
