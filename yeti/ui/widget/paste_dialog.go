package widget

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type PasteDialog struct {
	*dialog.CustomDialog
	start    *widget.Entry
	onClosed func(start int)
	win      fyne.Window
}

func NewPasteDialog(onClosed func(start int), parentWindow fyne.Window) *PasteDialog {
	p := &PasteDialog{
		onClosed: onClosed,
		win:      parentWindow,
	}
	p.start = widget.NewEntry()
	p.start.SetText("")
	p.CustomDialog = dialog.NewCustom("Paste frames", "Cancel",
		container.New(
			layout.NewGridLayoutWithRows(3),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel("Copy frame"),
			),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel("Copy after frame"),
				p.start,
			),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewButton("Apply", p.apply),
			),
		),
		p.win)
	return p
}

func (c *PasteDialog) apply() {

	startFrame, err := strconv.ParseUint(c.start.Text, 10, 32)
	if err != nil {
		return
	}

	if c.onClosed != nil {
		c.onClosed(int(startFrame))
	}
}
