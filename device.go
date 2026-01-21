package main

type DeviceProfile struct {
	Key         string
	Width       int
	Height      int
	TopInset    int // island / notch
	ClockInset  int // ЧАСЫ + системный текст
	BottomInset int // кнопки
}

var Devices = map[string]DeviceProfile{
	"iphone-se": {
		Key:         "iphone-se",
		Width:       750,
		Height:      1334,
		TopInset:    80,
		ClockInset:  260,
		BottomInset: 160,
	},
	"iphone-11": {
		Key:         "iphone-11",
		Width:       828,
		Height:      1792,
		TopInset:    120,
		ClockInset:  360,
		BottomInset: 200,
	},
	"iphone-14": {
		Key:         "iphone-14",
		Width:       1170,
		Height:      2532,
		TopInset:    160,
		ClockInset:  420,
		BottomInset: 240,
	},
	"iphone-15": {
		Key:         "iphone-15",
		Width:       1179,
		Height:      2556,
		TopInset:    180,
		ClockInset:  440,
		BottomInset: 260,
	},
}
