package main

import (
	"github.com/rivo/tview"
)

type View struct {
	repoList *tview.List
}

func Layout() *tview.Flex {
	view.repoList = tview.NewList()
	sidebar := tview.NewFlex()

	view.repoList.AddItem("Loading repos...", "", 0, nil)

	sidebar.AddItem(view.repoList, 0, 1, true).SetBorder(true).SetTitle("Repositories")
	title := tview.NewBox().SetBorder(true).SetTitle("Go gazer")
	content := tview.NewBox().SetBorder(true)
	footer := tview.NewBox().SetBorder(true)

	flex := tview.NewFlex().
		AddItem(sidebar, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(title, 3, 0, false).
			AddItem(content, 0, 3, false).
			AddItem(footer, 5, 1, false), 0, 3, false)

	return flex
}

func TextWidget(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}
