package main

type DeviceProfile struct {
	Name   string
	Width  int
	Height int
}

var Devices = map[string]DeviceProfile{
	"iphone-15": {
		Name:   "iPhone 15 / 15 Pro",
		Width:  1179,
		Height: 2556,
	},
	"iphone-15-plus": {
		Name:   "iPhone 15 Plus",
		Width:  1290,
		Height: 2796,
	},
	"iphone-14": {
		Name:   "iPhone 14 / 13 / 12",
		Width:  1170,
		Height: 2532,
	},
	"iphone-mini": {
		Name:   "iPhone 13 mini",
		Width:  1080,
		Height: 2340,
	},
	"iphone-se": {
		Name:   "iPhone SE",
		Width:  750,
		Height: 1334,
	},
}
