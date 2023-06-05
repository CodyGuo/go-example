package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"testing"
)

func Test_RunApp(t *testing.T) {
	var testCfg config
	testApp := test.NewApp()
	testCfg.Title = "Test Markdown"
	testWin := testApp.NewWindow(testCfg.Title)

	edit, preview := testCfg.makeUI(testWin)
	testCfg.createMenuItems(testWin)
	testWin.SetContent(container.NewHSplit(preview, edit))
	testApp.Run()

	test.Type(edit, "Some text")

	if preview.String() != "Some text" {
		t.Error("Run app failed")
	}
}
