package main

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidge    *widget.Entry
	PreviewWidge *widget.RichText
	CurrentFile  fyne.URI
	SaveMenuItem *fyne.MenuItem
}

var cfg config

func main() {
	a := app.New()
	// a.Settings().SetTheme(&myTheme{})
	win := a.NewWindow("Markdown")
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(win)

	win.SetContent(container.NewHSplit(edit, preview))
	win.Resize(fyne.Size{Width: 800, Height: 500})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	app.EditWidge = edit
	app.PreviewWidge = preview

	edit.OnChanged = preview.ParseMarkdown
	return edit, preview
}

func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open", app.createOpenFunc(win))

	saveMenuItem := fyne.NewMenuItem("Save", app.createSaveFunc(win))
	saveMenuItem.Disabled = true
	app.SaveMenuItem = saveMenuItem

	saveAsMenuItem := fyne.NewMenuItem("Save As", app.createSaveAsFunc(win))

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)
	win.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) createOpenFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			defer read.Close()
			if err != nil {
				dialog.ShowError(err, win)
			}
			if read == nil {
				return
			}
			data, err := os.ReadFile(read.URI().Path())
			if err != nil {
				dialog.ShowError(err, win)
			}
			app.EditWidge.SetText(string(data))
			app.CurrentFile = read.URI()
			win.SetTitle(win.Title() + " - " + read.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *config) createSaveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			writer, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			defer writer.Close()
			writer.Write([]byte(app.EditWidge.Text))
		}
	}
}

func (app *config) createSaveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			defer write.Close()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if write == nil {
				return
			}
			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please make sure you are opening a .md file.", win)
			}
			// save the file
			write.Write([]byte(app.EditWidge.Text))
			app.CurrentFile = write.URI()

			win.SetTitle(win.Title() + " - " + write.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}
