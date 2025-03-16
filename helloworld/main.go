package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	output *widget.Label
}

var myApp App

func main() {
	a := app.New()
	w := a.NewWindow("Hello world")

	output, entry, btn := myApp.makeUI()

	w.SetContent(container.NewVBox(
		output, entry, btn,
	))
	w.Resize(fyne.Size{Width: 500, Height: 500})
	w.ShowAndRun()
}

func (app *App) makeUI() (*widget.Label, *widget.Entry, *widget.Button) {
	output := widget.NewLabel("Hello world")
	app.output = output
	entry := widget.NewEntry()
	button := widget.NewButton("Enter", func() {
		app.output.SetText(entry.Text)
	})

	return output, entry, button
}
