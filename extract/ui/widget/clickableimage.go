package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	w2 "fyne.io/fyne/v2/widget"
)

type ClickableImage struct {
	w2.BaseWidget
	image         *canvas.Image
	toCallButton0 func(float32, float32)
	toCallButton1 func(float32, float32)
}

func NewClickableImage(f0 func(float32, float32), f1 func(float32, float32)) *ClickableImage {
	c := &ClickableImage{
		BaseWidget:    w2.BaseWidget{},
		toCallButton0: f0,
		toCallButton1: f1,
	}

	c.BaseWidget.ExtendBaseWidget(c)
	return c
}

func (ci *ClickableImage) Tapped(
	pe *fyne.PointEvent) {
	/*fmt.Printf("point X:%f,Y:%f and absolute position X:%f,Y:%f\n",
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

func (ci *ClickableImage) SetImage(i *canvas.Image) {
	ci.image = i
	width := i.Size().Width
	height := i.Size().Height

	ci.image.SetMinSize(
		fyne.NewSize(float32(width), float32(height)))
	ci.BaseWidget.ExtendBaseWidget(ci)
	ci.Refresh()
}

func (ci *ClickableImage) CreateRenderer() fyne.WidgetRenderer {
	ci.BaseWidget.ExtendBaseWidget(ci)
	return &clickableImageRenderer{
		image: ci.image,
		objs:  []fyne.CanvasObject{ci.image},
	}
}

type clickableImageRenderer struct {
	image *canvas.Image
	objs  []fyne.CanvasObject
}

func (ci *clickableImageRenderer) Destroy() {

}

func (ci *clickableImageRenderer) MinSize() fyne.Size {
	return ci.image.MinSize()
}

func (ci *clickableImageRenderer) Objects() []fyne.CanvasObject {
	return ci.objs
}

func (ci *clickableImageRenderer) Refresh() {
	canvas.Refresh(ci.image)
}

func (ci *clickableImageRenderer) Layout(size fyne.Size) {
	ci.image.Resize(size)
}
