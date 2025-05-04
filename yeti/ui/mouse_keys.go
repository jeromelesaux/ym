package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (u *ui) onTypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case "S":
		u.stop()
	case "P":
		u.play()
	case "O":
		u.OpenFileAction()
	case "E":
		u.SaveFileAction()
	case "R":
		u.ResetUI()
	case "D":
		u.DisplayChange()
	case "C":
		u.CancelChange()
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
