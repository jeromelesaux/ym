package ui

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/ym/cpc"
)

/*
* functions to handle ymcpc file format
* in table widget
 */

var (
	nbRegisterMax = 16 + 1
)

func (u *ui) updateCpcTableLength() (int, int) {
	return int(u.ymCpc.NbFrames) + 1, nbRegisterMax
}

func (u *ui) selectedCpcTableCell(id widget.TableCellID) {

	frame := id.Row - 1
	register := id.Col - 1
	if frame >= 0 && register >= 0 {
		fmt.Printf("register [%d] , frame [%d]\n", register, frame)
		msg := fmt.Sprintf("Set the value of the register [%d] frame [%d]", register, frame)
		// nolint: staticcheck
		de := dialog.NewEntryDialog("Set a new value", msg, func(v string) {
			fmt.Printf("new value [%s] register [%d] , frame [%d]\n", v, register, frame)
			frameValue, err := strconv.ParseInt("0x"+v, 0, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error while set the value :%v\n", err.Error())
				return
			}

			if register >= cpc.Register16bitsMaxIndice {
				if frameValue > 0xFF {
					fmt.Fprintf(os.Stderr, "Value [%X] exceed 0xff ", frameValue)
					return
				}
				fmt.Printf("new value [%d][%.2X] register [%d] , frame [%d]\n", frameValue, frameValue, register, frame)
				err := u.ymCpc.SetRegister8bits(register, frame, byte(frameValue))
				if err != nil {
					fmt.Printf("error with new value [%d][%.4X] register [%d] , frame [%d], error :%v\n", frameValue, frameValue, register, frame, err)

				}
			} else {
				fmt.Printf("new value [%d][%.4X] register [%d] , frame [%d]\n", frameValue, frameValue, register, frame)
				err := u.ymCpc.SetRegister16bits(register, frame, uint16(frameValue))
				if err != nil {
					fmt.Printf("error with new value [%d][%.4X] register [%d] , frame [%d], error :%v\n", frameValue, frameValue, register, frame, err)

				}
			}
			u.table.Refresh()
		}, u.window)
		de.Show()
	}

}

// nolint: ireturn
func (u *ui) updateCpcTableLabel() fyne.CanvasObject {
	item := widget.NewLabel("Template")
	item.Resize(fyne.Size{
		Width:  200,
		Height: 20,
	})
	return item
}

func (u *ui) updateCpcTableValue(id widget.TableCellID, cell fyne.CanvasObject) {
	label := cell.(*widget.Label)
	if id.Col >= 17 {
		return
	}
	if id.Row >= int(u.ymCpc.NbFrames)+1 {
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
			label.SetText(fmt.Sprintf("register %d", id.Col))
		} else {
			label.SetText(fmt.Sprintf("%.2X", u.ymCpc.Data[id.Col-1][id.Row-1]))
		}
	}
	label.Resize(fyne.Size{Height: 20, Width: 20})
	cell.(*widget.Label).Resize(fyne.Size{
		Width:  20,
		Height: 20,
	})
}
