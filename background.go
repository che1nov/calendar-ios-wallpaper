package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

type BackgroundStyle string
type BackgroundColor string

const (
	BgPlain    BackgroundStyle = "plain"
	BgGradient BackgroundStyle = "gradient"
	BgNoise    BackgroundStyle = "noise"
	BgIOS      BackgroundStyle = "ios"

	BgBlack  BackgroundColor = "black"
	BgBlue   BackgroundColor = "blue"
	BgPurple BackgroundColor = "purple"
	BgGreen  BackgroundColor = "green"
	BgRed    BackgroundColor = "red"
)

func backgroundBaseColor(c BackgroundColor) color.RGBA {
	switch c {
	case BgBlue:
		return color.RGBA{10, 20, 40, 255}
	case BgPurple:
		return color.RGBA{25, 10, 40, 255}
	case BgGreen:
		return color.RGBA{10, 40, 20, 255}
	case BgRed:
		return color.RGBA{40, 10, 10, 255}
	default:
		return color.RGBA{0, 0, 0, 255}
	}
}

func drawBackground(img *image.RGBA, device DeviceProfile, style BackgroundStyle, bgColor BackgroundColor) {
	base := backgroundBaseColor(bgColor)

	switch style {
	case BgPlain:
		fillSolid(img, base)
	case BgGradient:
		drawGradientWithBase(img, base)
	case BgNoise:
		drawNoiseWithBase(img, base)
	case BgIOS:
		drawPremiumBackgroundWithBase(img, device, base)
	default:
		drawPremiumBackgroundWithBase(img, device, base)
	}
}

/* =========================
   BACKGROUND STYLES
========================= */

func fillSolid(img *image.RGBA, base color.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, base)
		}
	}
}

func drawGradientWithBase(img *image.RGBA, base color.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	top := lighten(base, 1.2)
	bottom := darken(base, 0.4)

	for y := 0; y < h; y++ {
		t := float64(y) / float64(h)
		r := uint8(float64(top.R)*(1-t) + float64(bottom.R)*t)
		g := uint8(float64(top.G)*(1-t) + float64(bottom.G)*t)
		b := uint8(float64(top.B)*(1-t) + float64(bottom.B)*t)

		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	addVignette(img, 0.45)
}

func drawNoiseWithBase(img *image.RGBA, base color.RGBA) {
	fillSolid(img, base)
	rand.Seed(time.Now().UnixNano())

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			n := rand.Intn(16) - 8
			c := img.RGBAAt(x, y)
			img.Set(x, y, color.RGBA{
				uint8(clamp(int(c.R)+n, 0, 255)),
				uint8(clamp(int(c.G)+n, 0, 255)),
				uint8(clamp(int(c.B)+n, 0, 255)),
				255,
			})
		}
	}
}

func drawPremiumBackgroundWithBase(img *image.RGBA, device DeviceProfile, base color.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	top := lighten(base, 1.15)
	bottom := darken(base, 0.3)

	for y := 0; y < h; y++ {
		t := float64(y) / float64(h)
		r := uint8(float64(top.R)*(1-t) + float64(bottom.R)*t)
		g := uint8(float64(top.G)*(1-t) + float64(bottom.G)*t)
		b := uint8(float64(top.B)*(1-t) + float64(bottom.B)*t)

		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	clockBottom := device.ClockBottom()
	cx := float64(w) / 2
	cy := float64(clockBottom) / 2
	radius := float64(w) * 0.6

	for y := 0; y < clockBottom; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			d := math.Sqrt(dx*dx + dy*dy)

			if d < radius {
				k := 1 - d/radius
				c := img.RGBAAt(x, y)
				img.Set(x, y, color.RGBA{
					uint8(min(float64(c.R)+k*18, 255)),
					uint8(min(float64(c.G)+k*18, 255)),
					uint8(min(float64(c.B)+k*18, 255)),
					255,
				})
			}
		}
	}

	addVignette(img, 0.5)
}

/* =========================
   HELPERS
========================= */

func addVignette(img *image.RGBA, power float64) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	cx := float64(w) / 2
	cy := float64(h) / 2
	maxDist := math.Hypot(cx, cy)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			d := math.Hypot(dx, dy) / maxDist

			v := 1 - d*power
			if v < 0 {
				v = 0
			}

			c := img.RGBAAt(x, y)
			img.Set(x, y, color.RGBA{
				uint8(float64(c.R) * v),
				uint8(float64(c.G) * v),
				uint8(float64(c.B) * v),
				255,
			})
		}
	}
}

func lighten(c color.RGBA, k float64) color.RGBA {
	return color.RGBA{
		uint8(min(float64(c.R)*k, 255)),
		uint8(min(float64(c.G)*k, 255)),
		uint8(min(float64(c.B)*k, 255)),
		255,
	}
}

func darken(c color.RGBA, k float64) color.RGBA {
	return color.RGBA{
		uint8(float64(c.R) * k),
		uint8(float64(c.G) * k),
		uint8(float64(c.B) * k),
		255,
	}
}

func clamp(v, minv, maxv int) int {
	if v < minv {
		return minv
	}
	if v > maxv {
		return maxv
	}
	return v
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
