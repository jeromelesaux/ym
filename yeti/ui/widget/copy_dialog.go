package widget

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CopyDialog struct {
	*dialog.CustomDialog
	start    *widget.Entry
	end      *widget.Entry
	onClosed func(start, end int)
	win      fyne.Window
}

func NewCopyDialog(onClosed func(start, end int), parentWindow fyne.Window) *CopyDialog {
	c := &CopyDialog{
		onClosed: onClosed,
		win:      parentWindow,
	}
	c.start = widget.NewEntry()
	c.start.SetText("")
	c.end = widget.NewEntry()
	c.end.SetText("")

	c.CustomDialog = dialog.NewCustom("Copy frames", "Cancel",
		container.New(
			layout.NewGridLayoutWithRows(3),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel("Copy frame"),
			),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel("Copy from frame"),
				c.start,
				widget.NewLabel("To frame"),
				c.end,
			),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewButton("Apply", c.apply),
			),
		), c.win)
	return c
}

func (c *CopyDialog) apply() {

	startFrame, err := strconv.ParseUint(c.start.Text, 10, 32)
	if err != nil {
		return
	}

	endFrame, err := strconv.ParseUint(c.end.Text, 10, 32)
	if err != nil {
		return
	}

	if endFrame < startFrame {
		return
	}

	if c.onClosed != nil {
		c.onClosed(int(startFrame), int(endFrame))
	}
	c.Hide()
}
