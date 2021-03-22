package ui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/lha"
	"github.com/jeromelesaux/ym"
	"github.com/jeromelesaux/ym/encoding"
	chart "github.com/wcharczuk/go-chart"
)

type ui struct {
	filename                string
	ym                      *ym.Ym
	fileSongAuthor          *widget.Label
	fileSongName            *widget.Label
	fileSongDescription     *widget.Label
	fileFrameHz             *widget.Label
	rowStartSelected        *widget.Entry
	rowEndSelected          *widget.Entry
	frameStartSelectedIndex int
	frameEndSelectedIndex   int
	rowSelectionLayout      *container.Scroll

	table             *widget.Table
	graphicContent    *container.Scroll
	graphic           *canvas.Image
	registersSelected [16]bool
	window            fyne.Window
	lastDirectory     string
	ymToSave          *ym.Ym
	ymBackuped        *ym.Ym
	headerLevel       byte
	compressMethod    int
}

func (u *ui) onTypedKey(ev *fyne.KeyEvent) {
}

func (u *ui) onTypedRune(r rune) {
}

func (u *ui) generateChart() {
	series := []chart.Series{}
	xseries := make([]float64, u.ym.NbFrames)
	for i := 0; i < int(u.ym.NbFrames); i++ {
		xseries[i] = float64(i)
	}
	for i := 0; i < len(u.ym.Data); i++ {
		yseries := make([]float64, u.ym.NbFrames)
		for j := 0; j < int(u.ym.NbFrames); j++ {
			yseries[j] = float64(u.ym.Data[i][j])
		}
		serie := chart.ContinuousSeries{
			XValues: xseries,
			YValues: yseries,
		}
		series = append(series, serie)
	}
	graph := chart.Chart{
		Series: series,
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating chart : %v \n", err)
	}
	img, err := png.Decode(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while decoding png image : %v \n", err)
	}
	fw, _ := os.Create("tmp.png")

	png.Encode(fw, img)
	fw.Close()
	u.graphic = canvas.NewImageFromFile("tmp.png")
}

func (u *ui) updateTableLabel() fyne.CanvasObject {
	item := widget.NewLabel("Template")
	item.Resize(fyne.Size{
		Width:  200,
		Height: 20,
	})
	return item
}

func (u *ui) updateTableLength() (int, int) {
	return int(u.ym.NbFrames) + 1, 16 + 1
}

func (u *ui) selectedTableCell(id widget.TableCellID) {
}

func (u *ui) checkAllChanger(v bool) {
	u.check0Changer(v)
	u.check1Changer(v)
	u.check2Changer(v)
	u.check3Changer(v)
	u.check4Changer(v)
	u.check5Changer(v)
	u.check6Changer(v)
	u.check7Changer(v)
	u.check8Changer(v)
	u.check9Changer(v)
	u.check10Changer(v)
	u.check11Changer(v)
	u.check12Changer(v)
	u.check13Changer(v)
	u.check14Changer(v)
	u.check15Changer(v)
}

func (u *ui) check0Changer(v bool) {
	u.registersSelected[0] = v
}
func (u *ui) check1Changer(v bool) {
	u.registersSelected[1] = v
}
func (u *ui) check2Changer(v bool) {
	u.registersSelected[2] = v
}
func (u *ui) check3Changer(v bool) {
	u.registersSelected[3] = v
}
func (u *ui) check4Changer(v bool) {
	u.registersSelected[4] = v
}
func (u *ui) check5Changer(v bool) {
	u.registersSelected[5] = v
}
func (u *ui) check6Changer(v bool) {
	u.registersSelected[6] = v
}
func (u *ui) check7Changer(v bool) {
	u.registersSelected[7] = v
}
func (u *ui) check8Changer(v bool) {
	u.registersSelected[8] = v
}
func (u *ui) check9Changer(v bool) {
	u.registersSelected[9] = v
}
func (u *ui) check10Changer(v bool) {
	u.registersSelected[10] = v
}
func (u *ui) check11Changer(v bool) {
	u.registersSelected[11] = v
}
func (u *ui) check12Changer(v bool) {
	u.registersSelected[12] = v
}
func (u *ui) check13Changer(v bool) {
	u.registersSelected[13] = v
}
func (u *ui) check14Changer(v bool) {
	u.registersSelected[14] = v
}
func (u *ui) check15Changer(v bool) {
	u.registersSelected[15] = v
}

func (u *ui) endChange(v string) {
	end, err := strconv.Atoi(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error input %v, for value [%s]\n", err, v)
		return
	}
	u.frameEndSelectedIndex = end
	if end < 0 || end > int(u.ym.NbFrames) {
		u.frameEndSelectedIndex = int(u.ym.NbFrames)
	}
}

func (u *ui) startChange(v string) {
	start, err := strconv.Atoi(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error input %v, for value [%s]\n", err, v)
		return
	}
	u.frameStartSelectedIndex = start
	if start < 0 || start > int(u.ym.NbFrames) {
		u.frameStartSelectedIndex = 0
	}
}

func (u *ui) updateTableValue(id widget.TableCellID, cell fyne.CanvasObject) {
	label := cell.(*widget.Label)
	if id.Col > 16 {
		return
	}
	if id.Row >= int(u.ym.NbFrames)+1 {
		return
	}
	switch id.Col {
	case 0:
		if id.Row != 0 {
			label.SetText(fmt.Sprintf("frame %d", id.Row))
		} else {
			label.SetText("")
		}
	default:
		if id.Row == 0 {
			label.SetText(fmt.Sprintf("register %d", id.Col-1))
		} else {
			label.SetText(fmt.Sprintf("%d", u.ym.Data[id.Col-1][id.Row-1]))
		}
	}
	label.Resize(fyne.Size{Height: 20, Width: 20})
	cell.(*widget.Label).Resize(fyne.Size{
		Width:  20,
		Height: 20,
	})
}

func NewUI() *ui {
	u := &ui{}
	return u
}

func (u *ui) LoadUI(app fyne.App) {

	u.ym = ym.NewYm()

	u.fileSongAuthor = &widget.Label{Alignment: fyne.TextAlignTrailing}
	u.fileSongAuthor.TextStyle.Monospace = true
	u.fileSongAuthor.SetText("File song's Author :")
	u.fileSongAuthor.Resize(fyne.Size{Height: 10, Width: 50})

	u.fileSongName = &widget.Label{Alignment: fyne.TextAlignTrailing}
	u.fileSongName.TextStyle.Monospace = true
	u.fileSongName.SetText("File song's Name :")
	u.fileSongName.Resize(fyne.Size{Height: 10, Width: 50})

	u.fileSongDescription = &widget.Label{Alignment: fyne.TextAlignTrailing}
	u.fileSongDescription.TextStyle.Monospace = true
	u.fileSongDescription.SetText("File comment :")
	u.fileSongDescription.Resize(fyne.Size{Height: 10, Width: 50})

	u.fileFrameHz = &widget.Label{Alignment: fyne.TextAlignTrailing}
	u.fileFrameHz.TextStyle.Monospace = true
	u.fileFrameHz.SetText("Song frame in Hz:")
	u.fileFrameHz.Resize(fyne.Size{Height: 10, Width: 50})

	openButton := widget.NewButton("File Open ym file (.ym)", u.OpenFileAction)
	openButton.Resize(fyne.Size{Height: 1, Width: 50})

	saveButton := widget.NewButton("Save file", u.SaveFileAction)
	saveButton.Resize(fyne.Size{Height: 1, Width: 50})

	cleanButton := widget.NewButton("Clean or reset", u.ResetUI)
	cleanButton.Resize(fyne.Size{Height: 1, Width: 50})

	displayChangementsButton := widget.NewButton("Display changements", u.DisplayChange)
	displayChangementsButton.Resize(fyne.Size{Height: 1, Width: 50})

	returnToOriginalButton := widget.NewButton("Cancel changements", u.CancelChange)
	returnToOriginalButton.Resize(fyne.Size{Height: 1, Width: 50})

	/* registers check boxes selection */
	var registersSelectionCheckedButton = make([]*widget.Check, 17)
	type registerCheckFunc func(bool)
	var registersSelectFuncs = [17]registerCheckFunc{
		u.checkAllChanger,
		u.check0Changer,
		u.check1Changer,
		u.check2Changer,
		u.check3Changer,
		u.check4Changer,
		u.check5Changer,
		u.check6Changer,
		u.check7Changer,
		u.check8Changer,
		u.check9Changer,
		u.check10Changer,
		u.check11Changer,
		u.check12Changer,
		u.check13Changer,
		u.check14Changer,
		u.check15Changer}

	registerCheckLayout := fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithRows(17),
	)
	registersSelectionCheckedButton[0] = widget.NewCheck("select all registers", registersSelectFuncs[0])
	registerCheckLayout.Add(registersSelectionCheckedButton[0])
	for i := 1; i < 17; i++ {
		registersSelectionCheckedButton[i] = widget.NewCheck(fmt.Sprintf("register %d", i),
			registersSelectFuncs[i])
		registerCheckLayout.Add(registersSelectionCheckedButton[i])
	}

	/* end of creation  */

	u.rowEndSelected = widget.NewEntry()
	u.rowEndSelected.OnSubmitted = u.endChange
	startFrame := widget.NewLabel("Select the first frame (starts at 0)")
	u.rowStartSelected = widget.NewEntry()
	endFrame := widget.NewLabel("Select the last frame")
	u.rowStartSelected.OnSubmitted = u.startChange

	u.rowSelectionLayout = container.NewVScroll(fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithRows(4),
		startFrame,
		u.rowStartSelected,
		endFrame,
		u.rowEndSelected,
	))

	selectionLayout := fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(2),
		container.NewVScroll(registerCheckLayout),
		u.rowSelectionLayout,
	)

	u.table = widget.NewTable(
		u.updateTableLength,
		u.updateTableLabel,
		u.updateTableValue,
	)
	u.table.OnSelected = u.selectedTableCell

	u.generateChart()
	u.graphicContent = container.NewVScroll(u.graphic)

	u.window = app.NewWindow("YeTi")
	u.window.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewGridLayoutWithRows(2),
			fyne.NewContainerWithLayout(
				layout.NewGridLayoutWithRows(3),
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(4),
					u.fileSongAuthor,
					u.fileSongName,
					u.fileSongDescription,
					u.fileFrameHz,
				),
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(1),
					u.graphicContent,
				),
				fyne.NewContainerWithLayout(
					layout.NewGridLayoutWithRows(2),
					fyne.NewContainerWithLayout(
						layout.NewGridLayout(2),
						fyne.NewContainerWithLayout(
							layout.NewGridLayout(2),
							openButton,
							saveButton,
						),
						fyne.NewContainerWithLayout(
							layout.NewGridLayout(1),
							selectionLayout,
						)),
					fyne.NewContainerWithLayout(
						layout.NewGridLayout(3),
						fyne.NewContainerWithLayout(
							layout.NewGridLayout(1),
							displayChangementsButton,
						),
						fyne.NewContainerWithLayout(
							layout.NewGridLayout(1),
							returnToOriginalButton,
						),
						fyne.NewContainerWithLayout(
							layout.NewGridLayout(1),
							cleanButton,
						)),
				)),

			fyne.NewContainerWithLayout(
				layout.NewGridLayout(1),
				container.NewVScroll(u.table),
			),
		))
	u.window.Canvas().SetOnTypedRune(u.onTypedRune)
	u.window.Canvas().SetOnTypedKey(u.onTypedKey)
	u.window.Resize(fyne.NewSize(700, 600))

	u.window.Show()

}

func (u *ui) prepareExport() {
	u.ymToSave = ym.NewYm()
	copy(u.ymToSave.AuthorName, u.ym.AuthorName)
	copy(u.ymToSave.SongName, u.ym.SongName)
	copy(u.ymToSave.SongComment, u.ym.SongComment)
	copy(u.ymToSave.CheckString[:], u.ym.CheckString[:])
	u.ymToSave.DigidrumNb = u.ym.DigidrumNb
	u.ymToSave.EndID = u.ym.EndID
	u.ymToSave.FileID = u.ym.FileID
	u.ymToSave.FrameHz = u.ym.FrameHz
	u.ymToSave.LoopFrame = u.ym.LoopFrame
	u.ymToSave.Size = u.ym.Size
	u.ymToSave.SongAttributes = u.ym.SongAttributes
	u.ymToSave.YmMasterClock = u.ym.YmMasterClock
	length := u.frameEndSelectedIndex - u.frameStartSelectedIndex
	if length < 0 {
		return
	}
	for i := 0; i < 16; i++ {
		if u.registersSelected[i] {
			var j int
			if u.frameStartSelectedIndex != 0 {
				j = u.frameStartSelectedIndex - 1
			}
			for ; j < u.frameEndSelectedIndex; j++ {
				u.ymToSave.Data[i] = append(u.ymToSave.Data[i], u.ym.Data[i][j])
			}
		} else {
			u.ymToSave.Data[i] = make([]byte, length)
		}
	}

	u.ymToSave.NbFrames = uint32(length)
}

func (u *ui) DisplayChange() {
	var err error
	u.frameEndSelectedIndex, err = strconv.Atoi(u.rowEndSelected.Text)
	if err != nil {
		return
	}
	if u.frameEndSelectedIndex < 0 || u.frameEndSelectedIndex > int(u.ym.NbFrames) {
		u.frameEndSelectedIndex = 0
	}
	u.frameStartSelectedIndex, err = strconv.Atoi(u.rowStartSelected.Text)
	if err != nil {
		return
	}
	if u.frameStartSelectedIndex < 0 || u.frameStartSelectedIndex > int(u.ym.NbFrames) {
		u.frameStartSelectedIndex = int(u.ym.NbFrames)
	}
	u.ymBackuped = u.ym
	u.prepareExport()
	u.ym = u.ymToSave
	u.generateChart()
	u.graphicContent.Refresh()
	//	u.window.Resize(fyne.NewSize(700, 600))
}

func (u *ui) CancelChange() {
	u.ym = u.ymBackuped
	u.generateChart()
	u.graphicContent.Refresh()

	//	u.window.Resize(fyne.NewSize(700, 600))
}

func (u *ui) ResetUI() {
	u.ym = ym.NewYm()
	u.ymToSave = ym.NewYm()
	u.generateChart()
	u.graphicContent.Refresh()
}

func (u *ui) SaveFileAction() {
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err == nil && writer == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, u.window)
			return
		}
		filePath := strings.Replace(writer.URI().String(), writer.URI().Scheme()+"://", "", -1)
		if err = u.saveNewYm(filePath, writer); err != nil {
			dialog.ShowError(err, u.window)
			return
		}
		dialog.ShowInformation("file saving", "You file "+filePath+" is saved.", u.window)
	}, u.window)
	fd.Show()
}

func (u *ui) OpenFileAction() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, u.window)
			return
		}
		u.filename = reader.URI().Path()
		u.lastDirectory = reader.URI().Scheme() + "://" + filepath.Dir(reader.URI().Path())
		u.loadYmFile(reader)
		u.graphicContent.Refresh()

	}, u.window)
	uri, err := storage.ParseURI(u.lastDirectory)
	if err == nil {
		lister, err := storage.ListerForURI(uri)
		if err == nil {
			fd.SetLocation(lister)
		}
	}
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".ym"}))
	fd.Show()
}

func (u *ui) setFileDescription() {
	u.fileSongAuthor.SetText("File song's Author :" + string(u.ym.AuthorName))
	u.fileSongName.SetText("File song's Name :" + string(u.ym.SongName))
	u.fileSongDescription.SetText("File song's comment :" + string(u.ym.SongComment))
	u.fileFrameHz.SetText(fmt.Sprintf("Frame rate %d hz", u.ym.FrameHz))
}

func (u *ui) loadYmFile(f fyne.URIReadCloser) {
	f.Close()
	u.ym = ym.NewYm()
	var content []byte
	archive := lha.NewLha(u.filename)
	headers, err := archive.Headers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting lha headers file %s, error :%v\n", u.filename, err.Error())
		dialog.ShowError(err, u.window)
		return
	}
	if len(headers) > 0 {
		u.headerLevel = headers[0].HeaderLevel
		u.compressMethod = archive.CompressMethod
		content, err = archive.DecompresBytes(headers[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while decompressing file %s, error :%v\n", u.filename, err.Error())
			dialog.ShowError(err, u.window)
			return
		}
		err = encoding.Unmarshall(content, u.ym)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while decoding ym file %s, error :%v\n", u.filename, err.Error())
			dialog.ShowError(err, u.window)
			return
		}
	}
	u.generateChart()
	u.setFileDescription()

}

func (u *ui) saveNewYm(filePath string, writer fyne.URIWriteCloser) error {

	writer.Close()
	os.Remove(filePath)
	ym := u.ymToSave
	content, err := encoding.Marshall(ym)
	if err != nil {
		return err
	}
	archive := lha.NewLha(filePath)
	return archive.CompressBytes("archive.ym", content, u.compressMethod, int(u.headerLevel))
}
