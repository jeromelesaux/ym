package ui

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

/*
* functions to handle ym file format
* in table widget
 */

func (u *ui) updateTableLength() (int, int) {
	return int(u.ym.NbFrames) + 1, 16 + 1
}

func (u *ui) selectedTableCell(id widget.TableCellID) {

	frame := id.Row - 1
	register := id.Col - 1
	if frame >= 0 && register >= 0 {
		fmt.Printf("register [%d] , frame [%d]\n", register, frame)
		msg := fmt.Sprintf("Set the value of the register [%d] frame [%d]", register, frame)
		// nolint: staticcheck
		de := dialog.NewEntryDialog("Set a new value", msg, func(v string) {
			fmt.Printf("new value [%s] register [%d] , frame [%d]\n", v, register, frame)
			frameValue, err := strconv.ParseInt("0x"+v, 0, 16)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error while set the value :%v\n", err.Error())
				return
			}
			if frameValue > 0xFF {
				fmt.Fprintf(os.Stderr, "Value [%X] exceed 0xff ", frameValue)
				return
			}
			fmt.Printf("new value [%d][%.2X] register [%d] , frame [%d]\n", frameValue, frameValue, register, frame)
			u.ym.Data[register][frame] = byte(frameValue)
			u.table.Refresh()
			//	u.window.Resize(fyne.NewSize(700, 600))
		}, u.window)
		de.Show()
	}

}

// nolint: ireturn
func (u *ui) updateTableLabel() fyne.CanvasObject {
	item := widget.NewLabel("Template")
	item.Resize(fyne.Size{
		Width:  200,
		Height: 20,
	})
	return item
}

func (u *ui) updateTableValue(id widget.TableCellID, cell fyne.CanvasObject) {
	label := cell.(*widget.Label)
	if id.Col >= 17 {
		return
	}
	if id.Row >= int(u.ym.NbFrames)+1 {
		return
	}
	switch id.Col {
	case 0:
		if id.Row != 0 {
			label.SetText(fmt.Sprintf("%d", id.Row-1))
		} else {
			label.SetText("Frame(s)")
		}
	default:
		if id.Row == 0 {
			label.SetText(fmt.Sprintf("r%d", id.Col-1))
		} else {
			label.SetText(fmt.Sprintf("%.2X", u.ym.Data[id.Col-1][id.Row-1]))
		}
	}
	label.Resize(fyne.Size{Height: 20, Width: 20})
	cell.(*widget.Label).Resize(fyne.Size{
		Width:  20,
		Height: 20,
	})
}

func (u *ui) putInFrameCache(start, end int) {

	if u.ym.NbFrames <= uint32(end) {
		return
	}

	u.frameCache = [16][]byte{}
	for i := range 16 {
		u.frameCache[i] = append(u.frameCache[i], u.ym.Data[i][start:end]...)
	}
}

func (u *ui) copyAfterFrame(v int) {
	for i := range 16 {
		u.ym.Data[i] = slices.Insert(u.ym.Data[i], v, u.frameCache[i]...)
	}
	u.ym.NbFrames = uint32(len(u.ym.Data[0]))
	u.table.Refresh()
	canvas.Refresh(u.table)
}

func (u *ui) suppressFrame(start, end int) {
	for i := range 16 {
		u.ym.Data[i] = slices.Delete(u.ym.Data[i], start, end+1)
	}
	u.ym.NbFrames = uint32(len(u.ym.Data[0]))
	u.table.Refresh()
	canvas.Refresh(u.table)
}
