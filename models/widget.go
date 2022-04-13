package models

import "github.com/rivo/tview"

type Widget interface {
	Refresh()
	Component() *tview.List
}
