package main

// DeviceProfile описывает параметры экрана iPhone
// Все зоны считаются ОТ ВЫСОТЫ экрана
type DeviceProfile struct {
	Key    string
	Name   string
	Width  int
	Height int

	// Процент высоты экрана (0.0 – 1.0)
	ClockZoneRatio   float64 // где заканчивается зона часов / даты
	ButtonsZoneRatio float64 // где начинается зона кнопок

	BottomInset int // для будущего (home indicator)
}

// =======================
// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
// =======================

func (d DeviceProfile) ClockBottom() int {
	return int(float64(d.Height) * d.ClockZoneRatio)
}

func (d DeviceProfile) ButtonsTop() int {
	return int(float64(d.Height) * d.ButtonsZoneRatio)
}

func (d DeviceProfile) Scale() float64 {
	return float64(d.Width) / float64(BaseWidth)
}

// =======================
// БАЗОВЫЕ ПРОФИЛИ (ШАБЛОНЫ)
// =======================

var (
	profileSE23 = DeviceProfile{
		Name:             "iPhone SE (2 / 3)",
		Width:            750,
		Height:           1334,
		ClockZoneRatio:   0.26,
		ButtonsZoneRatio: 0.84,
	}

	profile12 = DeviceProfile{
		Name:             "iPhone 12 / 13 / 14",
		Width:            1170,
		Height:           2532,
		ClockZoneRatio:   0.30,
		ButtonsZoneRatio: 0.82,
		BottomInset:      34,
	}

	profile15 = DeviceProfile{
		Name:             "iPhone 15 / 16",
		Width:            1179,
		Height:           2556,
		ClockZoneRatio:   0.31,
		ButtonsZoneRatio: 0.81,
		BottomInset:      34,
	}

	profilePro = DeviceProfile{
		Name:             "iPhone Pro (Dynamic Island)",
		Width:            1179,
		Height:           2556,
		ClockZoneRatio:   0.32,
		ButtonsZoneRatio: 0.80,
		BottomInset:      34,
	}

	profileProMax = DeviceProfile{
		Name:             "iPhone Pro Max",
		Width:            1290,
		Height:           2796,
		ClockZoneRatio:   0.32,
		ButtonsZoneRatio: 0.80,
		BottomInset:      34,
	}
)

// =======================
// ВСЕ МОДЕЛИ (1:1 с UI)
// =======================

var Devices = map[string]DeviceProfile{

	// ===== Classic (Home Button) =====
	"iphone-se-1": {
		Key: "iphone-se-1", Name: "iPhone SE (1st gen)",
		Width: 640, Height: 1136,
		ClockZoneRatio: 0.25, ButtonsZoneRatio: 0.84,
	},
	"iphone-se-2": withKey("iphone-se-2", profileSE23),
	"iphone-se-3": withKey("iphone-se-3", profileSE23),
	"iphone-16e":  withKey("iphone-16e", profile15), // future

	// ===== Notch =====
	"iphone-x": {
		Key: "iphone-x", Name: "iPhone X / XS / 11 Pro",
		Width: 1125, Height: 2436,
		ClockZoneRatio: 0.30, ButtonsZoneRatio: 0.82,
		BottomInset: 34,
	},
	"iphone-xr": {
		Key: "iphone-xr", Name: "iPhone XR / 11",
		Width: 828, Height: 1792,
		ClockZoneRatio: 0.29, ButtonsZoneRatio: 0.83,
		BottomInset: 34,
	},
	"iphone-xs-max": {
		Key: "iphone-xs-max", Name: "iPhone XS Max / 11 Pro Max",
		Width: 1242, Height: 2688,
		ClockZoneRatio: 0.30, ButtonsZoneRatio: 0.82,
		BottomInset: 34,
	},

	// ===== Mini =====
	"iphone-12-mini": {
		Key: "iphone-12-mini", Name: "iPhone 12 mini",
		Width: 1080, Height: 2340,
		ClockZoneRatio: 0.30, ButtonsZoneRatio: 0.82,
		BottomInset: 34,
	},
	"iphone-13-mini": {
		Key: "iphone-13-mini", Name: "iPhone 13 mini",
		Width: 1080, Height: 2340,
		ClockZoneRatio: 0.30, ButtonsZoneRatio: 0.82,
		BottomInset: 34,
	},

	// ===== Standard =====
	"iphone-12": withKey("iphone-12", profile12),
	"iphone-13": withKey("iphone-13", profile12),
	"iphone-14": withKey("iphone-14", profile12),
	"iphone-15": withKey("iphone-15", profile15),
	"iphone-16": withKey("iphone-16", profile15),
	"iphone-17": withKey("iphone-17", profile15), // future

	// ===== Plus =====
	"iphone-14-plus": {
		Key: "iphone-14-plus", Name: "iPhone 14 Plus",
		Width: 1284, Height: 2778,
		ClockZoneRatio: 0.31, ButtonsZoneRatio: 0.81,
		BottomInset: 34,
	},
	"iphone-15-plus": {
		Key: "iphone-15-plus", Name: "iPhone 15 Plus",
		Width: 1290, Height: 2796,
		ClockZoneRatio: 0.31, ButtonsZoneRatio: 0.81,
		BottomInset: 34,
	},
	"iphone-16-plus": {
		Key: "iphone-16-plus", Name: "iPhone 16 Plus",
		Width: 1290, Height: 2796,
		ClockZoneRatio: 0.31, ButtonsZoneRatio: 0.81,
		BottomInset: 34,
	},

	// ===== Pro =====
	"iphone-12-pro": withKey("iphone-12-pro", profile12),
	"iphone-13-pro": withKey("iphone-13-pro", profile12),
	"iphone-14-pro": withKey("iphone-14-pro", profilePro),
	"iphone-15-pro": withKey("iphone-15-pro", profilePro),
	"iphone-16-pro": {
		Key: "iphone-16-pro", Name: "iPhone 16 Pro",
		Width: 1206, Height: 2622,
		ClockZoneRatio: 0.32, ButtonsZoneRatio: 0.80,
		BottomInset: 34,
	},
	"iphone-17-pro": withKey("iphone-17-pro", profilePro), // future

	// ===== Pro Max =====
	"iphone-12-pro-max": withKey("iphone-12-pro-max", profileProMax),
	"iphone-13-pro-max": withKey("iphone-13-pro-max", profileProMax),
	"iphone-14-pro-max": withKey("iphone-14-pro-max", profileProMax),
	"iphone-15-pro-max": withKey("iphone-15-pro-max", profileProMax),
	"iphone-16-pro-max": {
		Key: "iphone-16-pro-max", Name: "iPhone 16 Pro Max",
		Width: 1320, Height: 2868,
		ClockZoneRatio: 0.32, ButtonsZoneRatio: 0.80,
		BottomInset: 34,
	},
	"iphone-17-pro-max": withKey("iphone-17-pro-max", profileProMax), // future

	// ===== Air =====
	"iphone-air": withKey("iphone-air", profile15), // концепт
}

// =======================
// УТИЛИТА ДЛЯ АЛИАСОВ
// =======================

func withKey(key string, base DeviceProfile) DeviceProfile {
	base.Key = key
	return base
}
