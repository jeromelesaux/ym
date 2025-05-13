package widget

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FrameRange struct {
	From      uint64
	To        uint64
	fromEntry *widget.Entry
	toEntry   *widget.Entry
}

func newFrameRange() FrameRange {
	f := FrameRange{
		fromEntry: widget.NewEntry(),
		toEntry:   widget.NewEntry(),
	}
	f.fromEntry.SetText("")
	f.toEntry.SetText("")
	return f
}

type ReplaceDialog struct {
	*dialog.CustomDialog
	suppress FrameRange
	win      fyne.Window
	onClosed func(start, end int)
}

func NewReplaceDialog(onClosed func(start, end int), parentWindow fyne.Window) *ReplaceDialog {
	r := &ReplaceDialog{
		suppress: newFrameRange(),
		onClosed: onClosed,
		win:      parentWindow,
	}
	r.CustomDialog = dialog.NewCustom("Suppress Frames", "Cancel",
		container.New(
			layout.NewGridLayoutWithRows(4),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel("Suppress frame"),
			),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewLabel("Suppress frame range : "),
				widget.NewLabel("From start frame"),
				r.suppress.fromEntry,
				widget.NewLabel("To end frame"),
				r.suppress.toEntry,
			),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewButton("Apply", r.apply),
			),
		), r.win)
	return r
}

func (r *ReplaceDialog) apply() {
	var err error
	r.suppress.From, err = strconv.ParseUint(r.suppress.fromEntry.Text, 10, 32)
	if err != nil {
		return
	}
	r.suppress.To, err = strconv.ParseUint(r.suppress.toEntry.Text, 10, 32)
	if err != nil {
		return
	}

	if r.suppress.To < r.suppress.From {
		return
	}
	if r.onClosed != nil {
		r.onClosed(
			int(r.suppress.From),
			int(r.suppress.To))
	}
	r.Hide()
}
