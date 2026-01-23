package main

import (
	"image"
	"image/color"
	"math"
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
