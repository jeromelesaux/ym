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
	"time"

	"image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	wav2 "github.com/faiface/beep/wav"
	"github.com/jeromelesaux/lha"
	"github.com/jeromelesaux/ym"
	"github.com/jeromelesaux/ym/cpc"
	"github.com/jeromelesaux/ym/encoding"
	"github.com/jeromelesaux/ym/wav"
	w2 "github.com/jeromelesaux/ym/yeti/ui/widget"
	"github.com/jeromelesaux/ym/yeti/ui/xls"
	chart "github.com/wcharczuk/go-chart"
)

var (
	Appversion = "compilation_issue"
	dialogSize = fyne.NewSize(1000, 800)
)

type ui struct {
	filename                string
	ym                      *ym.Ym
	ymOrg                   *ym.Ym
	fileDescription         *widget.Label
	rowStartSelected        *widget.Entry
	rowEndSelected          *widget.Entry
	checkAllButton          *widget.Check
	checkCpcYm              *widget.Check
	frameStartSelectedIndex int
	frameEndSelectedIndex   int
	rowSelectionLayout      *container.Scroll

	table             *widget.Table
	graphicContent    *fyne.Container
	tableContainer    *container.Scroll
	graphic           *w2.ClickableImage
	registersSelected [16]bool
	window            fyne.Window
	lastDirectory     string
	ymToSave          *ym.Ym
	ymBackuped        *ym.Ym
	ymCpc             *cpc.CpcYM
	headerLevel       byte
	compressMethod    int
	archiveFilename   string
	lock              sync.Mutex
	graph             *chart.Chart
	speakerDone       chan bool
	playerTime        *widget.Label
	playerTimeTicker  *time.Ticker
	playerTimeChan    chan bool
	playerTimeValue   float64
	playerIsPlaying   bool
	playerProgression *widget.ProgressBar
	currentFrame      int
	frameCache        [16][]byte
}

func (u *ui) getCurrentYM() *ym.Ym {
	if u.checkCpcYm.Checked {
		return u.ymCpc.Ym
	}
	return u.ym
}

func (u *ui) generateChart() {
	u.lock.Lock()
	u.graph = nil
	series := []chart.Series{}
	maxX := u.ym.NbFrames
	if maxX > 800 {
		maxX = 800
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
		TitleStyle: chart.Style{Hidden: true},
		Elements:   []chart.Renderable{},
		YAxis: chart.YAxis{
			Style: chart.Style{Hidden: true},
			Range: &chart.ContinuousRange{
				Min: 0.0,
				Max: 255.0,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{Hidden: true},
		},
		Width:  1200,
		Height: 180,
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
	/*fw, _ := os.Create(graphicFileTemporaryFile)

	err = png.Encode(fw, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while encoding png image : %v \n", err)
	}
	fw.Close()*/
	u.graphic.SetImage(img)
	u.table.Select(widget.TableCellID{Row: 0, Col: 0})
}

func (u *ui) checkAllChanger(v bool) {
	if v {
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

func (u *ui) applyCpcYmFormat(v bool) {
	if v {
		u.ymOrg = ym.CopyYm(u.ym)
		u.ymCpc.Ym = CopyCPCYm(u.ymOrg)
		u.ymToSave = u.ymCpc.Ym
		// change the table functions to use cpc functions
		u.table.Length = u.updateCpcTableLength
		u.table.CreateCell = u.updateCpcTableLabel
		u.table.UpdateCell = u.updateCpcTableValue
		u.table.OnSelected = u.selectedCpcTableCell
		u.table.Refresh()
	} else {
		// change the table functions to use ym functions
		u.ymOrg = ym.CopyYm(u.ym)
		u.ymToSave = u.ym
		u.table.Length = u.updateTableLength
		u.table.CreateCell = u.updateTableLabel
		u.table.UpdateCell = u.updateTableValue
		u.table.OnSelected = u.selectedTableCell
		u.table.Refresh()

	}
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

func (u *ui) goToFrame(v string) {
	frame, err := strconv.Atoi(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error input %v, for value [%s]\n", err, v)
		return
	}
	frame++
	if frame > int(u.ym.NbFrames) {
		frame = int(u.ym.NbFrames)
	}
	u.currentFrame = frame
}
func NewUI() *ui {
	u := &ui{}
	return u
}

// nolint: funlen
func (u *ui) play() {
	if u.playerIsPlaying {
		return
	}
	u.playerTime.SetText("Decoding file.")

	var y *ym.Ym
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
	if u.frameEndSelectedIndex-u.frameStartSelectedIndex <= 0 {
		// no ym loaded.
		u.playerIsPlaying = false
		return
	}

	u.playerTimeTicker = time.NewTicker(time.Millisecond * 10)
	u.playerTimeValue = 0
	u.playerIsPlaying = true
	currentYm := u.getCurrentYM()
	y = currentYm.Extract(u.frameStartSelectedIndex, u.frameEndSelectedIndex)
	totalTime := float64(y.NbFrames) / float64(y.FrameHz)
	nbFrames := y.NbFrames

	go func() {
		for {
			select {
			case <-u.playerTimeChan:
				u.playerIsPlaying = false
				fyne.DoAndWait(func() {
					u.playerProgression.SetValue(0)
					u.playerTime.SetText(u.playerTime.Text + "\nPlayer stopped.")
					u.table.Select(widget.TableCellID{Row: u.currentFrame, Col: 0})
				})
				return
			case <-u.playerTimeTicker.C:
				u.playerTimeValue += .01
				u.currentFrame = int(u.playerTimeValue / totalTime * float64(nbFrames))
				label := fmt.Sprintf("Time: %.2f seconds  Frame: %d", u.playerTimeValue, u.currentFrame)
				fyne.DoAndWait(
					func() {
						u.playerTime.SetText(label)
						u.playerProgression.SetValue(u.playerTimeValue / totalTime)
					},
				)

			}
		}
	}()
	go func() {

		v := wav.NewYMMusic()
		err := v.LoadMemory(y)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while loading memory error:%v\n", err)
		}
		content, err := v.Wave()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while converting ym to wave with error :%v\n", err.Error())
			return
		}
		r := bytes.NewReader(content)
		u.playerTimeValue = 0
		streamer, format, err := wav2.Decode(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while streaming wave with error :%v\n", err.Error())
			return
		}
		defer streamer.Close()

		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			u.speakerDone <- true
			//	u.playerTimeChan <- true
			fmt.Printf("Googbye go routine\n")
		})))
		fmt.Printf("Speaker play the new file %v\n", format)

		for {
			select {
			case <-u.speakerDone:
				streamer.Close()
				speaker.Clear()
				u.playerTimeChan <- true
				fmt.Printf("Now the speaker is cleared\n")
				return
			}
		}

	}()

}

func (u *ui) stop() {
	if u.playerIsPlaying {
		u.speakerDone <- true
	}

}

func (u *ui) prepareExport() {
	u.ymToSave = ym.CopyYm(u.getCurrentYM())
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
	//	u.graphicContent.Refresh()
	u.rowStartSelected.SetText("")
	u.rowEndSelected.SetText("")
	wait.Hide()
	//	u.window.Resize(fyne.NewSize(700, 600))
}

func (u *ui) CancelChange() {
	wait := dialog.NewInformation("Get original file", "Please wait...", u.window)
	wait.Show()
	u.ym = u.ymBackuped
	u.setFileDescription()
	u.generateChart()
	//	u.graphicContent.Refresh()
	wait.Hide()
	//	u.window.Resize(fyne.NewSize(700, 600))
}

func (u *ui) ResetUI() {
	u.ym = ym.NewYm()
	u.ymToSave = ym.NewYm()
	u.ymBackuped = ym.NewYm()
	u.ymCpc = cpc.NewCpcYM()
	u.setFileDescription()
	u.generateChart()
	u.graphicContent.Refresh()
}

func (u *ui) ExportExcel() {
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err == nil && writer == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, u.window)
			return
		}

		filePath := writer.URI().Path()
		if filepath.Ext(filePath) != ".xlsx" {
			filePath += ".xlsx"
		}
		xl := xls.XlsFile{}
		if err := xl.New(filePath, u.ym.Data); err != nil {
			dialog.ShowError(err, u.window)
			return
		}
		dialog.ShowInformation("Xlsx file saved", "Your file is saved in ["+filePath+"].", u.window)
	}, u.window)
	fd.Resize(dialogSize)
	fd.Show()
}

func (u *ui) ImportExcel() {
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
		xl := xls.XlsFile{}
		reader.Close()
		data, err := xl.Get(reader.URI().Path())
		if err != nil {
			dialog.ShowError(err, u.window)
			return
		}
		u.newYM(data)
		u.generateChart()
		u.setFileDescription()
		//	u.graphicContent.Refresh()
		alert.Hide()

	}, u.window)
	uri, err := storage.ParseURI(u.lastDirectory)
	if err == nil {
		lister, err := storage.ListerForURI(uri)
		if err == nil {
			fd.SetLocation(lister)
		}
	}
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".xls", ".xlsx"}))
	fd.Resize(dialogSize)
	fd.Show()
}

func (u *ui) ExportRegisters() {
	fd := dialog.NewFolderOpen(func(lister fyne.ListableURI, err error) {
		if err == nil && lister == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, u.window)
			return
		}
		folderPath := lister.Path()
		for i := range 16 {
			filePath := folderPath + string(filepath.Separator) + filepath.Base(u.filename) + ".r" + fmt.Sprintf("%.2d", i)
			fw, err := os.Create(filePath)
			if err != nil {
				dialog.ShowError(err, u.window)
				return
			}
			reg := u.ym.Data[i]
			_, err = fw.Write(reg)
			if err != nil {
				fw.Close()
				dialog.ShowError(err, u.window)
				return
			}
			fw.Close()
		}

		dialog.ShowInformation("registers saved", "Your files are saved in ["+folderPath+"].", u.window)
	}, u.window)
	fd.Resize(dialogSize)
	fd.Show()
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
		//	u.graphicContent.Refresh()
		alert.Hide()

	}, u.window)
	uri, err := storage.ParseURI(u.lastDirectory)
	if err == nil {
		lister, err := storage.ListerForURI(uri)
		if err == nil {
			fd.SetLocation(lister)
		}
	}
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".ym", ".ym5", ".bin"}))
	fd.Resize(dialogSize)
	fd.Show()
}

func (u *ui) setFileDescription() {
	desc := fmt.Sprintf("File song's Author :%+q\n", string(u.ym.AuthorName))
	desc += fmt.Sprintf("File song's Name :%+q\n", string(u.ym.SongName))
	desc += fmt.Sprintf("File song's comment :%+q\n", string(u.ym.SongComment))
	desc += fmt.Sprintf("Duration %.2f seconds\n", float64(u.ym.NbFrames)/float64(u.ym.FrameHz))
	desc += fmt.Sprintf("Frame rate %d hz", u.ym.FrameHz) + "\n"
	desc += fmt.Sprintf("Number of frame %d ", u.ym.NbFrames) + "\n"
	desc += fmt.Sprintf("Frame loop at %d", u.ym.LoopFrame) + "\n"
	desc += fmt.Sprintf("Number of digidrums: %d", u.ym.DigidrumNb) + "\n"
	desc += fmt.Sprintf("YM format :%s\n", u.ym.FormatType())

	var clock string = "ATARI-ST"
	if u.ym.YmMasterClock == ym.AMSTRAD_CLOCK {
		clock = "AMSTRAD CPC"
	}
	desc += fmt.Sprintf("YM Master clock :%d %s", u.ym.YmMasterClock, clock) + "\n"
	u.fileDescription.SetText(desc)
}

func (u *ui) loadYmFile(f fyne.URIReadCloser) {
	f.Close()
	u.ym = nil
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
	} else {
		content, err = os.ReadFile(u.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while decoding ym file %s, error :%v\n", u.filename, err.Error())
			dialog.ShowError(err, u.window)
			return
		}
	}
	err = encoding.Unmarshall(content, u.ym)
	if err != nil && io.EOF != err {
		fmt.Fprintf(os.Stderr, "Error while decoding ym file %s, error :%v\n", u.filename, err.Error())
		dialog.ShowError(err, u.window)
		return
	}
	fmt.Printf("NB frames:[%d]\n", u.ym.NbFrames)

	u.generateChart()
	u.setFileDescription()

}

func (u *ui) saveNewYm(filePath string, writer fyne.URIWriteCloser) error {
	writer.Close()
	os.Remove(filePath)
	// force to last version YM
	if u.ymToSave.FileID != ym.YM1 && u.ymToSave.FileID != ym.YM2 {
		if u.ymToSave.FileID != ym.YM6 {
			u.ymToSave.FileID = ym.YM6
		}
	}

	// check if the ymTosave is not empty
	if len(u.ymToSave.Data[0]) == 0 {
		u.frameEndSelectedIndex = int(u.ym.NbFrames)
		u.prepareExport()
	}
	content, err := encoding.Marshall(u.ymToSave)
	if err != nil {
		return err
	}
	archive := lha.NewLha(filePath)
	lha.GenericFormat = true
	return archive.CompressBytes(u.archiveFilename, content, u.compressMethod, int(u.headerLevel))
}

func (u *ui) newYM(data [16][]byte) {
	u.ym = ym.NewYm()
	u.ym.Data = data
	u.ym.NbFrames = uint32(len(u.ym.Data[0]))
	u.compressMethod = 5
	u.ymToSave = u.ym

}
