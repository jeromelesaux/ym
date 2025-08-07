package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (u *ui) onTypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyS:
		u.stop()
	case fyne.KeyP:
		u.play()
	case fyne.KeyO:
		u.OpenFileAction()
	case fyne.KeyE:
		u.SaveFileAction()
	case fyne.KeyR:
		u.ResetUI()
	case fyne.KeyD:
		u.DisplayChange()
	case fyne.KeyC:
		u.CancelChange()
	case fyne.KeySpace:
		u.playFrame()
	case fyne.KeyDown,
		fyne.KeyLeft, fyne.KeyRight,
		fyne.KeyUp:
		u.table.TypedKey(ev)
	default:
		fmt.Printf("name:%s\n", ev.Name)

	}
}

func (u *ui) onTypedRune(r rune) {
	switch r {
	default:
		fmt.Printf("name:%v\n", r)

	}
}

func (u *ui) Tapped(
	x float32, y float32) {
	//size := u.graphicContent.Size()
	size := u.graphic.Size()
	percentage := x / size.Width * 100.
	frame := 0
	//fmt.Printf("percentage :%f\n", percentage)
	if percentage < 0.5 {
		frame = 0
	} else {
		frame = int((float32(u.ym.NbFrames) * percentage) / 100.)
	}
	fmt.Printf("gotoframe %d\n", frame)
	u.table.Select(widget.TableCellID{Row: frame, Col: 0})
	u.table.Refresh()
	// min 1,4 %
	// max 95 %
}
