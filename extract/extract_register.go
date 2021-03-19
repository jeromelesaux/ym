package main

import (
	"errors"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/jeromelesaux/ym/extract/ui"
)

var ErrorIsNotDirectory = errors.New("Is not a directory, Quiting.")

var ()

func main() {
	os.Setenv("FYNE_SCALE", "0.75")
	/* main application */
	app := app.NewWithID("Ym extract tool (YET^^)")
	/* set icon application */
	app.SetIcon(theme.FyneLogo())
	/* set new window */

	u := ui.NewUI()
	u.LoadUI(app)
	app.Run()
}
