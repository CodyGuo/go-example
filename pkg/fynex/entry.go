package fynex

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ShortcutEntry struct {
	widget.Entry
	ShortcutFunc func(shortcut fyne.Shortcut)
}

func NewMultiLineShortcutEntry() *ShortcutEntry {
	se := &ShortcutEntry{Entry: widget.Entry{MultiLine: true, Wrapping: fyne.TextTruncate}}
	se.ExtendBaseWidget(se)
	return se
}

func (se *ShortcutEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if _, ok := shortcut.(*desktop.CustomShortcut); !ok {
		se.Entry.TypedShortcut(shortcut)
		return
	}

	if se.ShortcutFunc == nil {
		return
	}

	se.ShortcutFunc(shortcut)
}
