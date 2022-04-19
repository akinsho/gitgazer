package ui

import "github.com/rivo/tview"

type Widget interface {
	Refresh()
	Component() tview.Primitive
	IsEmpty() bool
}

type TextWidget interface {
	Widget
	ScrollUp()
	ScrollDown()
}
