package main

import (
	"errors"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/jeromelesaux/ym/extract/ui"
)

// nolint: deadcode
var ErrorIsNotDirectory = errors.New("is not a directory, Quiting")

var ()

func main() {
	os.Setenv("FYNE_SCALE", "0.6")
	/* main application */
	app := app.NewWithID("Ym extract tool By ImPact(YETI^)")
	/* set icon application */
	icon, err := fyne.LoadResourceFromPath("icon/YeTi.png")
	if err != nil {
		app.SetIcon(icon)
	} else {
		// nolint: staticcheck
		app.SetIcon(theme.FyneLogo())
	}
	/* set new window */

	u := ui.NewUI()
	u.LoadUI(app)
	app.Run()
}
