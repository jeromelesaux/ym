package ui

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/jeromelesaux/ym"
	"github.com/jeromelesaux/ym/cpc"
	w2 "github.com/jeromelesaux/ym/yeti/ui/widget"
)

// nolint: funlen
func (u *ui) LoadUI(app fyne.App) {

	u.ym = ym.NewYm()
	u.ymBackuped = ym.NewYm()
	u.ymToSave = ym.NewYm()
	u.ymCpc = cpc.NewCpcYM()
	u.archiveFilename = "archive.ym"

	format := beep.Format{SampleRate: ym.Frame44Khz}
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while initialising speaker : %v\n", err)
		return
	}
	fmt.Printf("Speaker init ok\n")
	u.fileDescription = widget.NewLabel("File Description")
	u.fileDescription.TextStyle.Monospace = true
	u.fileDescription.SetText("File song's Author :")
	u.fileDescription.Resize(fyne.Size{Height: 10, Width: 50})

	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), u.play)
	stopButton := widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), u.stop)
	u.playerTime = widget.NewLabel("Time:")

	openButton := widget.NewButton("File Open ym file (.ym)", u.OpenFileAction)
	openButton.Resize(fyne.Size{Height: 1, Width: 50})

	saveButton := widget.NewButton("Save file", u.SaveFileAction)
	saveButton.Resize(fyne.Size{Height: 1, Width: 50})

	exportRegisters := widget.NewButton("Export registers", u.ExportRegisters)
	exportRegisters.Resize(fyne.Size{Height: 1, Width: 50})

	cleanButton := widget.NewButton("Clean or reset", u.ResetUI)
	cleanButton.Resize(fyne.Size{Height: 1, Width: 50})

	importExcel := widget.NewButton("Import data (Excel)", u.ImportExcel)
	importExcel.Resize(fyne.Size{Height: 1, Width: 50})

	exportExcel := widget.NewButton("Export data (Excel)", u.ExportExcel)
	exportExcel.Resize(fyne.Size{Height: 1, Width: 50})

	displayChangementsButton := widget.NewButton("Display changements", u.DisplayChange)
	displayChangementsButton.Resize(fyne.Size{Height: 1, Width: 50})

	returnToOriginalButton := widget.NewButton("Cancel changements", u.CancelChange)
	returnToOriginalButton.Resize(fyne.Size{Height: 1, Width: 50})
	goToFrameLabel := widget.NewLabel("Go to frame")
	gotToFrameEntry := widget.NewEntry()
	gotToFrameEntry.OnSubmitted = u.goToFrame

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

	registerCheckLayout := container.New(
		layout.NewGridLayoutWithRows(17),
	)
	registersSelectionCheckedButton[0] = widget.NewCheck("select all registers", registersSelectFuncs[0])
	u.checkAllButton = registersSelectionCheckedButton[0]
	registerCheckLayout.Add(registersSelectionCheckedButton[0])
	for i := 1; i < 17; i++ {
		registersSelectionCheckedButton[i] = widget.NewCheck(fmt.Sprintf("register %d", i-1),
			registersSelectFuncs[i])
		registersSelectionCheckedButton[i].SetChecked(true)
		registerCheckLayout.Add(registersSelectionCheckedButton[i])
	}
	registersSelectionCheckedButton[0].SetChecked(true)

	u.checkCpcYm = widget.NewCheck("CPC", u.applyCpcYmFormat)

	u.rowEndSelected = widget.NewEntry()
	u.rowEndSelected.OnSubmitted = u.endChange
	startFrame := widget.NewLabel("Select the first frame (starts at 0)")
	u.rowStartSelected = widget.NewEntry()
	endFrame := widget.NewLabel("Select the last frame")
	u.rowStartSelected.OnSubmitted = u.startChange

	u.rowSelectionLayout = container.NewVScroll(
		container.New(
			layout.NewGridLayoutWithColumns(6),
			container.NewVScroll(registerCheckLayout),
			u.checkCpcYm,
			startFrame,
			u.rowStartSelected,
			endFrame,
			u.rowEndSelected,
		))
	//u.rowSelectionLayout.Resize(fyne.NewSize(200, 20))
	u.playerProgression = widget.NewProgressBar()
	selectionLayout := container.New(
		layout.NewGridLayoutWithRows(3),
		u.playerProgression,
		container.New(
			layout.NewGridLayoutWithColumns(5),
			goToFrameLabel,
			gotToFrameEntry,
			displayChangementsButton,
			returnToOriginalButton,
			cleanButton,
		),
		container.New(
			layout.NewGridLayoutWithColumns(1),
			u.rowSelectionLayout,
		),
	)
	//selectionLayout.Resize(fyne.NewSize(400, 20))
	u.table = widget.NewTable(
		u.updateTableLength,
		u.updateTableLabel,
		u.updateTableValue,
	)
	u.table.OnSelected = u.selectedTableCell
	u.tableContainer = container.NewVScroll(u.table)

	u.graphic = w2.NewClickableImage(u.Tapped, nil)
	u.generateChart()
	// nolint: staticcheck
	u.graphicContent = container.New(layout.NewStackLayout(), u.graphic)
	//u.graphicContent = container.NewContainerWithLayout(layout.NewMaxLayout())

	u.window = app.NewWindow("YeTi")
	u.window.SetContent(
		container.New(
			layout.NewGridLayoutWithRows(2),
			container.New(
				layout.NewGridLayoutWithRows(3),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					container.NewVScroll(
						container.New(
							layout.NewGridLayoutWithColumns(1),
							container.NewVScroll(u.fileDescription),
						)),
					container.New(
						layout.NewGridLayoutWithColumns(2),
						container.New(layout.NewGridLayoutWithRows(3),
							openButton,
							saveButton,
							exportRegisters,
						),
						container.New(
							layout.NewGridLayoutWithRows(3),
							importExcel,
							exportExcel,
							container.New(
								layout.NewGridLayoutWithColumns(3),
								u.playerTime,
								playButton,
								stopButton,
							),
						),
					),
				),
				u.graphicContent,
				container.New(
					layout.NewGridLayoutWithRows(1),
					selectionLayout,
				)),

			container.New(
				layout.NewGridLayout(1),
				u.tableContainer,
			),
		))
	u.window.Canvas().SetOnTypedRune(u.onTypedRune)
	u.window.Canvas().SetOnTypedKey(u.onTypedKey)
	u.window.Resize(fyne.NewSize(400, 900))
	u.window.SetTitle("YeTi @ " + Appversion)
	u.window.Show()

}
