package ui

import (
	"image"
	"image/color"
	"image/draw"
	"os"
)

func drawCursor(filePath string, x int) (error, image.Image) {
	fr, err := os.Open(filePath)
	if err != nil {
		return err, nil
	}
	defer fr.Close()

	img, _, err := image.Decode(fr)
	if err != nil {
		return err, nil
	}
	imgDst := image.NewNRGBA(img.Bounds())
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			imgDst.Set(x, y, img.At(x, y))
		}
	}
	red := color.RGBA{R: 255, B: 0, G: 0, A: 255}

	draw.Draw(imgDst, image.Rect(x, 0, x+5, 160), &image.Uniform{red}, image.Pt(100, 0), draw.Src)
	//err := png.Encode(fw, imgDst)
	return err, imgDst
}
