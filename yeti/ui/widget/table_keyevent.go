package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TableKeyEvent struct {
	*widget.Table
	onTypeKey func(ev *fyne.KeyEvent)
}

func NewTable(length func() (rows int, cols int), create func() fyne.CanvasObject, update func(widget.TableCellID, fyne.CanvasObject), onTypedKey func(ev *fyne.KeyEvent)) *TableKeyEvent {
	t := &TableKeyEvent{
		Table:     widget.NewTable(length, create, update),
		onTypeKey: onTypedKey,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *TableKeyEvent) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeyDown,
		fyne.KeyLeft, fyne.KeyRight,
		fyne.KeyUp:
		t.Table.TypedKey(event)
	default:
		if t.onTypeKey != nil {
			t.onTypeKey(event)
		}
	}
	t.Table.Refresh()
}
