package ui

import "github.com/rivo/tview"

type Widget interface {
	Refresh() error
	Component() tview.Primitive
	IsEmpty() bool
}

type ListWidget interface {
	Widget
	SetSelected(int)
}

type TextWidget interface {
	Widget
	ScrollUp()
	ScrollDown()
}
