package main

type DeviceProfile struct {
	Key string

	Width  int
	Height int

	TopSafe    int // часы / island
	BottomSafe int // кнопки

	FooterOffset int // отступ футера от кнопок

	MonthFont  float64
	FooterFont float64

	DotRadius  int
	DotSpacing int
}

var Devices = map[string]DeviceProfile{
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

	"iphone-11": {
		Key:          "iphone-11",
		Width:        828,
		Height:       1792,
		TopSafe:      260,
		BottomSafe:   150,
		FooterOffset: 60,

		MonthFont:  32,
		FooterFont: 26,
		DotRadius:  8,
		DotSpacing: 28,
	},

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
