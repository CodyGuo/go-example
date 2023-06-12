package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/data/binding"
)

type Config struct {
	App        fyne.App
	MainWindow fyne.Window
	AddWindow  fyne.Window
	EditWindow fyne.Window
	List       *widget.List
	Data       binding.StringList
	file       string
}

var myApp Config

func main() {
	fyneApp := app.New()

	myApp.App = fyneApp
	myApp.MainWindow = fyneApp.NewWindow("List Data")
	myApp.file = "data.json"

	loadedData := loadJsonData(myApp.file)

	myApp.Data = binding.NewStringList()
	myApp.Data.Set(loadedData)

	defer saveJsonData(myApp.file, myApp.Data)

	myApp.MainWindow.Resize(fyne.NewSize(400, 600))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()
	myApp.MainWindow.CenterOnScreen()

	myApp.makeUI()

	myApp.MainWindow.ShowAndRun()
}

func (app *Config) makeUI() {
	list := app.getDataList()
	add := app.getAddBtn()
	exit := widget.NewButton("Quit", func() {
		app.MainWindow.Close()
	})

	content := container.NewBorder(nil, container.New(layout.NewVBoxLayout(), add, exit), nil, nil, list)
	app.MainWindow.SetContent(content)
}

func (app *Config) getDataList() *widget.List {
	list := widget.NewListWithData(app.Data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	app.List = list
	app.List.OnSelected = app.selectedFunc()

	return list
}

func (app *Config) selectedFunc() func(id widget.ListItemID) {
	editWin := app.App.NewWindow("Edit Data")
	app.EditWindow = editWin
	app.EditWindow.SetCloseIntercept(func() {
		app.EditWindow.Hide()
	})
	return func(id widget.ListItemID) {
		app.List.Unselect(id)
		d, _ := app.Data.GetValue(id)

		itemName := widget.NewEntry()
		itemName.Text = d

		updateData := widget.NewButton("Update", app.updateDataFunc(id, itemName))
		cancel := widget.NewButton("Cancel", app.CancelFunc())
		deleteData := widget.NewButton("Delete", app.DeleteFunc(id))

		app.EditWindow.SetContent(container.New(layout.NewVBoxLayout(), itemName, updateData, deleteData, cancel))
		app.EditWindow.Resize(fyne.NewSize(400, 200))
		app.EditWindow.CenterOnScreen()
		app.EditWindow.Show()
	}
}

func (app *Config) updateDataFunc(id int, item *widget.Entry) func() {
	return func() {
		app.Data.SetValue(id, item.Text)
		app.EditWindow.Hide()
	}
}

func (app *Config) CancelFunc() func() {
	return func() {
		app.EditWindow.Hide()
	}
}

func (app *Config) DeleteFunc(id int) func() {
	return func() {
		var newData []string
		dt, _ := app.Data.Get()

		for index, item := range dt {
			if index != id {
				newData = append(newData, item)
			}
		}

		app.Data.Set(newData)

		app.EditWindow.Hide()
	}
}

func (app *Config) getAddBtn() *widget.Button {
	addWin := app.App.NewWindow("Add Data")
	app.AddWindow = addWin
	app.AddWindow.SetCloseIntercept(func() {
		app.AddWindow.Hide()
	})

	add := widget.NewButton("Add", func() {
		itemName := widget.NewEntry()

		addData := widget.NewButton("Add", func() {
			app.Data.Append(itemName.Text)
			app.AddWindow.Hide()
		})

		cancel := widget.NewButton("Cancel", func() {
			app.AddWindow.Hide()
		})

		app.AddWindow.SetContent(container.New(layout.NewVBoxLayout(), itemName, addData, cancel))
		app.AddWindow.Resize(fyne.NewSize(400, 200))
		app.AddWindow.CenterOnScreen()
		app.AddWindow.Show()
	})

	return add
}

func loadJsonData(file string) []string {
	fmt.Println("Loading data from JSON file", file)

	input, _ := os.ReadFile(file)
	var data []string
	json.Unmarshal(input, &data)

	return data
}

func saveJsonData(file string, data binding.StringList) {
	fmt.Println("Saving data to JSON file", file)
	d, _ := data.Get()
	jsonData, _ := json.Marshal(d)
	os.WriteFile(file, jsonData, 0644)
}
