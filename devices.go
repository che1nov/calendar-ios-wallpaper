package main

type DeviceProfile struct {
	Key        string
	Width      int
	Height     int
	TopSafe    int
	BottomSafe int

	MonthFont  float64
	FooterFont float64
	DotRadius  int
	DotSpacing int
}

var Devices = map[string]DeviceProfile{
	// ===== SMALL =====
	"iphone-se": {
		Key:        "iphone-se",
		Width:      750,
		Height:     1334,
		TopSafe:    110,
		BottomSafe: 90,

		MonthFont:  26,
		FooterFont: 20,
		DotRadius:  6,
		DotSpacing: 22,
	},

	// ===== NOTCH =====
	"iphone-11": {
		Key:        "iphone-11",
		Width:      828,
		Height:     1792,
		TopSafe:    260,
		BottomSafe: 150,

		MonthFont:  32,
		FooterFont: 26,
		DotRadius:  8,
		DotSpacing: 28,
	},

	// ===== DYNAMIC ISLAND =====
	"iphone-14": {
		Key:        "iphone-14",
		Width:      1170,
		Height:     2532,
		TopSafe:    360,
		BottomSafe: 200,

		MonthFont:  38,
		FooterFont: 32,
		DotRadius:  9,
		DotSpacing: 32,
	},

	"iphone-15": {
		Key:        "iphone-15",
		Width:      1179,
		Height:     2556,
		TopSafe:    380,
		BottomSafe: 210,

		MonthFont:  40,
		FooterFont: 34,
		DotRadius:  9,
		DotSpacing: 34,
	},

	// ===== PRO MAX =====
	"iphone-15-pro-max": {
		Key:        "iphone-15-pro-max",
		Width:      1290,
		Height:     2796,
		TopSafe:    420,
		BottomSafe: 240,

		MonthFont:  44,
		FooterFont: 38,
		DotRadius:  10,
		DotSpacing: 38,
	},
}
