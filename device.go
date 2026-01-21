package main

type DeviceProfile struct {
	Key string

	Width  int
	Height int

	TopSafe    int
	BottomSafe int

	FooterOffset int

	MonthFont  float64
	FooterFont float64

	DotRadius  int
	DotSpacing int
}

var Devices = map[string]DeviceProfile{
	// ===== SE =====
	"iphone-se": {
		Key:          "iphone-se",
		Width:        750,
		Height:       1334,
		TopSafe:      120,
		BottomSafe:   90,
		FooterOffset: 40,

		MonthFont:  26,
		FooterFont: 20,
		DotRadius:  6,
		DotSpacing: 22,
	},

	// ===== MINI =====
	"iphone-13-mini": {
		Key:          "iphone-13-mini",
		Width:        1080,
		Height:       2340,
		TopSafe:      300,
		BottomSafe:   200,
		FooterOffset: 60,

		MonthFont:  32,
		FooterFont: 26,
		DotRadius:  8,
		DotSpacing: 28,
	},

	// ===== STANDARD =====
	"iphone-15": {
		Key:          "iphone-15",
		Width:        1179,
		Height:       2556,
		TopSafe:      380,
		BottomSafe:   210,
		FooterOffset: 70,

		MonthFont:  40,
		FooterFont: 34,
		DotRadius:  9,
		DotSpacing: 34,
	},

	// ===== PRO MAX =====
	"iphone-15-pro-max": {
		Key:          "iphone-15-pro-max",
		Width:        1290,
		Height:       2796,
		TopSafe:      420,
		BottomSafe:   240,
		FooterOffset: 80,

		MonthFont:  44,
		FooterFont: 38,
		DotRadius:  10,
		DotSpacing: 38,
	},
}
