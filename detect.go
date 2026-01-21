package main

import (
	"net/http"
	"strings"
)

// DetectDeviceProfile
// 1) если device передан явно → используем его
// 2) иначе пытаемся определить по User-Agent
// 3) fallback → iphone-15
func DetectDeviceProfile(r *http.Request) DeviceProfile {
	ua := strings.ToLower(r.UserAgent())

	// --- НЕ iPhone ---
	if !strings.Contains(ua, "iphone") {
		return Devices["iphone-15"]
	}

	// --- iPhone SE ---
	if strings.Contains(ua, "se") {
		return Devices["iphone-se"]
	}

	// --- Mini ---
	if strings.Contains(ua, "mini") {
		if d, ok := Devices["iphone-13-mini"]; ok {
			return d
		}
	}

	// --- Pro Max ---
	if strings.Contains(ua, "pro max") {
		if d, ok := Devices["iphone-15-pro-max"]; ok {
			return d
		}
	}

	// --- Pro ---
	if strings.Contains(ua, "pro") {
		return Devices["iphone-15"]
	}

	// --- Default ---
	return Devices["iphone-15"]
}

func detectDeviceFromUA(ua string) DeviceProfile {
	ua = strings.ToLower(ua)

	switch {
	case strings.Contains(ua, "iphone se"):
		return Devices["iphone-se"]
	case strings.Contains(ua, "iphone"):
		return Devices["iphone-15"]
	default:
		return Devices["iphone-15"]
	}
}
