package widget

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	w2 "fyne.io/fyne/v2/widget"
)

type ClickableImage struct {
	*w2.Icon
	image         *canvas.Image
	toCallButton0 func(float32, float32)
	toCallButton1 func(float32, float32)
}

func NewClickableImage(f0 func(float32, float32), f1 func(float32, float32)) *ClickableImage {
	c := &ClickableImage{
		Icon:          &w2.Icon{},
		image:         &canvas.Image{},
		toCallButton0: f0,
		toCallButton1: f1,
	}
	c.ExtendBaseWidget(c)
	c.BaseWidget.ExtendBaseWidget(c)
	return c
}

func (ci *ClickableImage) Tapped(
	pe *fyne.PointEvent) {
	/*	fmt.Printf("point X:%f,Y:%f and absolute position X:%f,Y:%f\n",
		pe.Position.X,
		pe.Position.Y,
		pe.AbsolutePosition.X,
		pe.AbsolutePosition.Y)*/
	if ci.toCallButton0 != nil {
		(ci.toCallButton0)(pe.Position.X, pe.Position.Y)
	}

}

func (ci *ClickableImage) TappedSecondary(pe *fyne.PointEvent) {
	if ci.toCallButton1 != nil {
		(ci.toCallButton1)(pe.Position.X, pe.Position.Y)
	}
}

func (ci *ClickableImage) SetImage(img image.Image) {
	ci.image.Image = img
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	ci.image.SetMinSize(
		fyne.NewSize(float32(width), float32(height)))
	ci.image.Refresh()
	//ci.ExtendBaseWidget(ci)
	//ci.Refresh()
	//ci.BaseWidget.Refresh()
}

// nolint: ireturn
func (ci *ClickableImage) CreateRenderer() fyne.WidgetRenderer {
	//ci.BaseWidget.ExtendBaseWidget(ci)
	return &clickableImageRenderer{
		image: ci.image,
		objs:  []fyne.CanvasObject{ci.image},
	}
}

func (ci *ClickableImage) Move(position fyne.Position) {
	ci.Icon.Move(position)
}

type clickableImageRenderer struct {
	image *canvas.Image
	objs  []fyne.CanvasObject
}

func (ci *clickableImageRenderer) Destroy() {
	ci.image = nil
}

func (ci *clickableImageRenderer) MinSize() fyne.Size {
	return ci.image.MinSize()
}

func (ci *clickableImageRenderer) Objects() []fyne.CanvasObject {
	return ci.objs
}

func (ci *clickableImageRenderer) Refresh() {
	ci.image.Refresh()
}

func (ci *clickableImageRenderer) Layout(size fyne.Size) {
	//	fmt.Printf("new size layout:%f,%f\n", size.Width, size.Height)
	ci.image.Resize(size)
}

func (ci *clickableImageRenderer) Resize(size fyne.Size) {
	//	fmt.Printf("new size resize :%f,%f\n", size.Width, size.Height)
	ci.image.Resize(size)
}
