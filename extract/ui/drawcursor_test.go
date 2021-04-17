package ui

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"
	"testing"
)

func TestDrawCursor(t *testing.T) {
	fr, err := os.Open("../yeti-gfx-cache.png")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer fr.Close()

	img, _, err := image.Decode(fr)
	if err != nil {
		t.Fatalf(err.Error())
	}
	imgDst := image.NewNRGBA(img.Bounds())
	fw, err := os.Create("yeti-gfx-cache-cursored.png")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer fw.Close()
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			imgDst.Set(x, y, img.At(x, y))
		}
	}
	red := color.RGBA{R: 255, B: 0, G: 0, A: 255}
	x := 300
	draw.Draw(imgDst, image.Rect(x, 0, x+5, 160), &image.Uniform{red}, image.Pt(100, 0), draw.Src)
	png.Encode(fw, imgDst)
}
