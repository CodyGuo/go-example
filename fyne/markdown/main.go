package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/codyguo/go-example/fyne/pkg/fynex"
	"io"
	"strings"
)

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

type config struct {
	Title          string
	EditWidget     *fynex.ShortcutEntry
	PreviewWidget  *widget.RichText
	CurrentFile    fyne.URI
	SaveMenuItem   *fyne.MenuItem
	SaveAsMenuItem *fyne.MenuItem
}

var cfg config

func main() {
	a := app.New()
	a.Settings().SetTheme(&customTheme{})

	cfg.Title = "Markdown"
	win := a.NewWindow(cfg.Title)

	edit, preview := cfg.makeUI(win)
	cfg.createMenuItems(win)

	win.SetContent(container.NewHSplit(preview, edit))
	win.Resize(fyne.NewSize(1024, 800))
	win.CenterOnScreen()

	win.ShowAndRun()
}

func (app *config) makeUI(win fyne.Window) (*fynex.ShortcutEntry, *widget.RichText) {
	edit := fynex.NewMultiLineShortcutEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.FocusGained()
	edit.ShortcutFunc = func(shortcut fyne.Shortcut) {
		v, ok := shortcut.(*desktop.CustomShortcut)
		if !ok {
			return
		}
		switch v.ShortcutName() {
		case fynex.CtrlS.ShortcutName():
			if strings.TrimSpace(app.EditWidget.Text) == "" {
				return
			}
			if app.CurrentFile == nil {
				app.SaveAsMenuItem.Action()
			} else {
				app.SaveMenuItem.Action()
			}
		}
	}
	edit.OnChanged = func(s string) {
		app.PreviewWidget.ParseMarkdown(s)
		app.SaveMenuItem.Disabled = false
		app.setWinTitle(win, true)
	}

	preview.Wrapping = fyne.TextWrapWord

	return edit, preview
}

func (app *config) setWinTitle(win fyne.Window, changed bool) {
	if app.CurrentFile == nil {
		return
	}
	str := " - "
	if changed {
		str = " - *"
	}
	win.SetTitle(app.Title + str + app.CurrentFile.Name())
}

func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open ...", app.openFunc(win))
	saveMenuItem := fyne.NewMenuItem("Save", app.SaveFunc(win))
	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.SaveAsFunc(win))

	app.SaveMenuItem = saveMenuItem
	app.SaveAsMenuItem = saveAsMenuItem

	app.SaveMenuItem.Disabled = true
	app.SaveMenuItem.Shortcut = fynex.CtrlS
	win.Canvas().AddShortcut(app.SaveMenuItem.Shortcut, func(shortcut fyne.Shortcut) {
		app.SaveMenuItem.Action()
	})

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if read == nil {
				return
			}

			defer read.Close()

			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			app.EditWidget.SetText(string(data))
			app.SaveMenuItem.Disabled = true
			app.CurrentFile = read.URI()
			app.setWinTitle(win, false)
		}, win)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *config) SaveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile == nil {
			app.SaveAsMenuItem.Action()
			return
		}

		write, err := storage.Writer(app.CurrentFile)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		if write == nil {
			return
		}

		defer write.Close()
		write.Write([]byte(app.EditWidget.Text))

		app.SaveMenuItem.Disabled = true
		app.setWinTitle(win, false)
	}
}

func (app *config) SaveAsFunc(win fyne.Window) func() {
	return func() {
		saveAsDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if write == nil {
				return
			}

			defer write.Close()

			if !filter.Matches(write.URI()) {
				dialog.ShowInformation("Error", "Please name your file with a .md extension!", win)
				return
			}

			write.Write([]byte(app.EditWidget.Text))

			app.CurrentFile = write.URI()
			app.SaveMenuItem.Disabled = true
			app.setWinTitle(win, false)
		}, win)

		saveAsDialog.SetFileName("untitled.md")
		saveAsDialog.SetFilter(filter)
		saveAsDialog.Show()
	}
}
