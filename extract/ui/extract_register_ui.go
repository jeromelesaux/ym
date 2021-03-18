package ui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

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

var ()

type ui struct {
	filename            string
	ym                  ym.Ym
	fileSongAuthor      *widget.Label
	fileSongName        *widget.Label
	fileSongDescription *widget.Label
	fileFrameHz         *widget.Label
	table               *widget.Table
	graphic             *canvas.Image
	window              fyne.Window
	lastDirectory       string
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
	defer fw.Close()
	png.Encode(fw, img)
	u.graphic = canvas.NewImageFromFile("tmp.png")
	u.graphic.Refresh()
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
	return int(u.ym.NbFrames), 16
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
		}
	default:
		if id.Row == 0 {
			label.SetText(fmt.Sprintf("register %d", id.Col-1))
		} else {
			label.SetText(fmt.Sprintf("%d", u.ym.Data[id.Col][id.Row]))
		}
	}
	label.Resize(fyne.Size{Height: 20, Width: 20})
	cell.(*widget.Label).Resize(fyne.Size{
		Width:  200,
		Height: 20,
	})
}

func NewUI() *ui {
	u := &ui{}
	return u
}

func (u *ui) LoadUI(app fyne.App) {

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

	u.table = widget.NewTable(
		u.updateTableLength,
		u.updateTableLabel,
		u.updateTableValue,
	)

	u.generateChart()

	u.window = app.NewWindow("YeT")
	u.window.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewGridLayoutWithColumns(1),
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
					u.graphic,
				),
				/*	fyne.NewContainerWithLayout(
						layout.NewGridLayout(1),
						u.fileSongName,
					),
					fyne.NewContainerWithLayout(
						layout.NewGridLayout(1),
						u.fileSongDescription,
					),
					fyne.NewContainerWithLayout(
						layout.NewGridLayout(1),
						u.fileFrameHz,
					),*/
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					openButton,
					saveButton,
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

func (u *ui) SaveFileAction() {
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err == nil && writer == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, u.window)
			return
		}

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
	defer f.Close()
	var content []byte
	archive := lha.NewLha(u.filename)
	headers, err := archive.Headers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting lha headers file %s, error :%v\n", u.filename, err.Error())
		dialog.ShowError(err, u.window)
		return
	}
	if len(headers) > 0 {
		content, err = archive.DecompresBytes(headers[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while decompressing file %s, error :%v\n", u.filename, err.Error())
			dialog.ShowError(err, u.window)
			return
		}
		err = encoding.Unmarshall(content, &u.ym)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while decoding ym file %s, error :%v\n", u.filename, err.Error())
			dialog.ShowError(err, u.window)
			return
		}
	}
	u.generateChart()
	u.setFileDescription()

}
