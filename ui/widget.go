package ui

import (
	"akinsho/gitgazer/app"

	"github.com/rivo/tview"
)

type Widget interface {
	Refresh() error
	Open() error
	Context() *app.Context
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
