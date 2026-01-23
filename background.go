package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

func drawBackground(img *image.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// Верх и низ градиента
	top := color.RGBA{10, 10, 14, 255} // чуть светлее
	bottom := color.RGBA{0, 0, 0, 255} // чёрный

	// --- Градиент ---
	for y := 0; y < h; y++ {
		t := float64(y) / float64(h)

		r := uint8(float64(top.R)*(1-t) + float64(bottom.R)*t)
		g := uint8(float64(top.G)*(1-t) + float64(bottom.G)*t)
		b := uint8(float64(top.B)*(1-t) + float64(bottom.B)*t)

		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// --- Виньетка ---
	cx := float64(w) / 2
	cy := float64(h) / 2
	maxDist := math.Hypot(cx, cy)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			d := math.Hypot(dx, dy) / maxDist

			// коэффициент затемнения
			v := 1.0 - d*0.45
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

func drawPremiumBackground(img *image.RGBA, device DeviceProfile) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// === ГРАДИЕНТ ===
	top := color.RGBA{12, 12, 16, 255}
	bottom := color.RGBA{0, 0, 0, 255}

	for y := 0; y < h; y++ {
		t := float64(y) / float64(h)
		r := uint8(float64(top.R)*(1-t) + float64(bottom.R)*t)
		g := uint8(float64(top.G)*(1-t) + float64(bottom.G)*t)
		b := uint8(float64(top.B)*(1-t) + float64(bottom.B)*t)

		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// === ПОДСВЕТКА ЗОНЫ ЧАСОВ ===
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

	// === ВИНЬЕТКА ===
	cvx := float64(w) / 2
	cvy := float64(h) / 2
	maxDist := math.Hypot(cvx, cvy)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cvx
			dy := float64(y) - cvy
			d := math.Hypot(dx, dy) / maxDist

			v := 1 - d*0.5
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

	// === ШУМ (GRAIN) ===
	rand.Seed(time.Now().UnixNano())
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			n := rand.Intn(12) - 6
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

func drawCalendarPanel(img *image.RGBA, top, bottom int) {
	w := img.Bounds().Dx()

	for y := top; y < bottom; y++ {
		for x := 0; x < w; x++ {
			c := img.RGBAAt(x, y)
			img.Set(x, y, color.RGBA{
				uint8(float64(c.R) * 0.85),
				uint8(float64(c.G) * 0.85),
				uint8(float64(c.B) * 0.85),
				255,
			})
		}
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
