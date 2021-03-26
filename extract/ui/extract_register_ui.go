package ui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

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

var (
	Appversion = "new layout and force ym6 export"
	dialogSize = fyne.NewSize(1000, 800)
)

type ui struct {
	filename                string
	ym                      *ym.Ym
	fileDescription         *widget.Label
	rowStartSelected        *widget.Entry
	rowEndSelected          *widget.Entry
	frameStartSelectedIndex int
	frameEndSelectedIndex   int
	rowSelectionLayout      *container.Scroll

	table             *widget.Table
	graphicContent    *container.Scroll
	tableContainer    *container.Scroll
	graphic           *canvas.Image
	registersSelected [16]bool
	window            fyne.Window
	lastDirectory     string
	ymToSave          *ym.Ym
	ymBackuped        *ym.Ym
	headerLevel       byte
	compressMethod    int
	archiveFilename   string
	lock              sync.Mutex
	graph             *chart.Chart
}

/*
func (u *ui) onTypedKey(ev *fyne.KeyEvent) {
}

func (u *ui) onTypedRune(r rune) {
}
*/
func (u *ui) generateChart() {
	u.lock.Lock()
	series := []chart.Series{}
	maxX := u.ym.NbFrames
	if maxX > 500 {
		maxX = 500
	}
	xseries := make([]float64, maxX)
	for i := 0; i < int(maxX); i++ {
		xseries[i] = float64(i)
	}
	for i := 0; i < len(u.ym.Data); i++ {
		yseries := make([]float64, maxX)
		for j := 0; j < int(maxX); j++ {
			index := int(u.ym.NbFrames/maxX) * j
			yseries[j] = float64(u.ym.Data[i][index])
		}
		serie := chart.ContinuousSeries{
			XValues: xseries,
			YValues: yseries,
		}
		series = append(series, serie)
	}
	u.lock.Unlock()
	u.graph = &chart.Chart{
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: 0.0,
				Max: 255.0,
			},
		},
		Width:  1800,
		Height: 800,

		Series: series,
	}
	buffer := bytes.NewBuffer([]byte{})
	err := u.graph.Render(chart.PNG, buffer)
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

	frame := id.Row - 1
	register := id.Col - 1
	if frame >= 0 && register >= 0 {
		fmt.Printf("register [%d] , frame [%d]\n", register, frame)
		msg := fmt.Sprintf("Set the value of the register [%d] frame [%d]", register, frame)
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
			label.SetText(fmt.Sprintf("register %d", id.Col-1))
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

func NewUI() *ui {
	u := &ui{}
	return u
}

func (u *ui) LoadUI(app fyne.App) {

	u.ym = ym.NewYm()
	u.ymBackuped = ym.NewYm()
	u.ymToSave = ym.NewYm()
	u.archiveFilename = "archive.ym"

	u.fileDescription = widget.NewLabel("File Description")
	u.fileDescription.TextStyle.Monospace = true
	u.fileDescription.SetText("File song's Author :")
	u.fileDescription.Resize(fyne.Size{Height: 10, Width: 50})

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
		layout.NewGridLayoutWithColumns(4),
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
	u.tableContainer = container.NewVScroll(u.table)
	u.generateChart()
	u.graphicContent = container.NewVScroll(u.graphic)

	u.window = app.NewWindow("YeTi")
	u.window.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewGridLayoutWithRows(2),
			fyne.NewContainerWithLayout(
				layout.NewGridLayoutWithRows(3),
				fyne.NewContainerWithLayout(
					layout.NewGridLayoutWithRows(1),
					container.NewVScroll(
						fyne.NewContainerWithLayout(
							layout.NewGridLayoutWithColumns(2),
							u.fileDescription,
						))),
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(1),
					u.graphicContent,
				),
				fyne.NewContainerWithLayout(
					layout.NewGridLayoutWithRows(2),
					container.New(layout.NewVBoxLayout(),

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
				u.tableContainer,
			),
		))
	//	u.window.Canvas().SetOnTypedRune(u.onTypedRune)
	//	u.window.Canvas().SetOnTypedKey(u.onTypedKey)
	u.window.Resize(fyne.NewSize(700, 600))
	u.window.SetTitle("YeTi @ " + Appversion)
	u.window.Show()

}

func (u *ui) prepareExport() {
	u.ymToSave = ym.CopyYm(u.ym)
	u.ymToSave.LoopFrame = 0
	length := u.frameEndSelectedIndex - u.frameStartSelectedIndex + 1
	if length < 0 {
		return
	}
	for i := 0; i < 16; i++ {
		u.ymToSave.Data[i] = make([]byte, length)
		if u.registersSelected[i] {
			var j, j2 int
			if u.frameStartSelectedIndex != 0 {
				j = u.frameStartSelectedIndex
			}
			for ; j < u.frameEndSelectedIndex; j++ {
				u.ymToSave.Data[i][j2] = u.ym.Data[i][j]
				j2++
			}
		}
	}

	u.ymToSave.NbFrames = uint32(length)
}

func (u *ui) DisplayChange() {
	wait := dialog.NewInformation("Applying changements", "Please wait...", u.window)
	wait.Show()
	var err error
	u.frameEndSelectedIndex, err = strconv.Atoi(u.rowEndSelected.Text)
	if err != nil {
		u.frameEndSelectedIndex = int(u.ym.NbFrames) - 1
	}
	if u.frameEndSelectedIndex < 0 || u.frameEndSelectedIndex > int(u.ym.NbFrames) {
		u.frameEndSelectedIndex = int(u.ym.NbFrames) - 1
	}
	u.frameStartSelectedIndex, err = strconv.Atoi(u.rowStartSelected.Text)
	if err != nil {
		u.frameStartSelectedIndex = 0
	}
	if u.frameStartSelectedIndex < 0 || u.frameStartSelectedIndex > int(u.ym.NbFrames) {
		u.frameStartSelectedIndex = 0
	}
	u.ymBackuped = u.ym
	u.prepareExport()
	u.ym = u.ymToSave
	u.setFileDescription()
	u.generateChart()
	u.graphicContent.Refresh()
	wait.Hide()
	//	u.window.Resize(fyne.NewSize(700, 600))
}

func (u *ui) CancelChange() {
	wait := dialog.NewInformation("Get original file", "Please wait...", u.window)
	wait.Show()
	u.ym = u.ymBackuped
	u.setFileDescription()
	u.generateChart()
	u.graphicContent.Refresh()
	wait.Hide()
	//	u.window.Resize(fyne.NewSize(700, 600))
}

func (u *ui) ResetUI() {
	u.ym = ym.NewYm()
	u.ymToSave = ym.NewYm()
	u.ymBackuped = ym.NewYm()
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
	fd.Resize(dialogSize)
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
		alert := dialog.NewInformation("loading file", "Please Wait", u.window)
		alert.Show()
		u.loadYmFile(reader)
		u.graphicContent.Refresh()
		alert.Hide()

	}, u.window)
	uri, err := storage.ParseURI(u.lastDirectory)
	if err == nil {
		lister, err := storage.ListerForURI(uri)
		if err == nil {
			fd.SetLocation(lister)
		}
	}
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".ym"}))
	fd.Resize(dialogSize)
	fd.Show()
}

func (u *ui) setFileDescription() {
	desc := "File song's Author :" + string(u.ym.AuthorName) + "\n"
	desc += "File song's Name :" + string(u.ym.SongName) + "\n"
	desc += "File song's comment :" + string(u.ym.SongComment) + "\n"
	desc += fmt.Sprintf("Frame rate %d hz", u.ym.FrameHz) + "\n"
	desc += fmt.Sprintf("Number of frame %d ", u.ym.NbFrames) + "\n"
	desc += fmt.Sprintf("Frame loop at %d", u.ym.LoopFrame) + "\n"
	desc += fmt.Sprintf("Number of digidrums: %d", u.ym.DigidrumNb) + "\n"

	var clock string = "ATARI-ST"
	if u.ym.YmMasterClock == encoding.AMSTRAD_CLOCK {
		clock = "AMSTRAD CPC"
	}
	desc += fmt.Sprintf("YM Master clock :%d %s", u.ym.YmMasterClock, clock) + "\n"
	u.fileDescription.SetText(desc)
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
		fmt.Printf("informations: \n\tinternal name:[%s]\n\tCompress Method:[%d]\n\tHeader level:[%d]\n\t",
			string(headers[0].Name),
			archive.CompressMethod,
			headers[0].HeaderLevel)
		u.headerLevel = 0
		u.compressMethod = 5
		u.archiveFilename = string(headers[0].Name)
		content, err = archive.DecompresBytes(headers[0])
		if err != nil && len(content) < headers[0].OriginalSize {
			fmt.Fprintf(os.Stderr, "Error while decompressing file %s, error :%v\n", u.filename, err.Error())
			dialog.ShowError(err, u.window)
			return
		}
		//	f, _ := os.Create("dump.bin")
		//	defer f.Close()
		//	f.Write(content)
		err = encoding.Unmarshall(content, u.ym)
		if err != nil && io.EOF != err {
			fmt.Fprintf(os.Stderr, "Error while decoding ym file %s, error :%v\n", u.filename, err.Error())
			dialog.ShowError(err, u.window)
			return
		}
		fmt.Printf("NB frames:[%d]\n", u.ym.NbFrames)

	}
	// force to last version YM
	if u.ym.FileID != encoding.YM6 {
		u.ym.FileID = encoding.YM6
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
	lha.GenericFormat = true
	return archive.CompressBytes(u.archiveFilename, content, u.compressMethod, int(u.headerLevel))
}
