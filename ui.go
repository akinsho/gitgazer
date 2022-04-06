package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type View struct {
	repoList *tview.List
	mainView *tview.TextView
}

func getAppUI() *tview.Flex {
	view.repoList = tview.NewList()
	view.mainView = tview.NewTextView()
	sidebar := tview.NewFlex()

	view.repoList.AddItem("Loading repos...", "", 0, nil)

	view.repoList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		repo := getRepositoryByIndex(index)
		text := fmt.Sprintf(
			"%s\n%s\n[red]issue count[white]: %d",
			repo.GetName(),
			repo.GetDescription(),
			repo.GetOpenIssuesCount(),
		)
		view.mainView.SetText(text)
	})

	sidebar.AddItem(view.repoList, 0, 1, true).SetBorder(true).SetTitle("Repositories")

	title := textWidget("Go Gazer")
	footer := tview.NewBox().SetBorder(true)

	view.mainView.SetDynamicColors(true).SetBorder(true)
	title.SetBorder(true)

	flex := tview.NewFlex().
		AddItem(sidebar, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(title, 3, 0, false).
			AddItem(view.mainView, 0, 3, false).
			AddItem(footer, 5, 1, false), 0, 3, false)

	return flex
}

func textWidget(text string) *tview.TextView {
	widget := tview.NewTextView()
	widget.
		SetTextAlign(tview.AlignCenter).
		SetText(text)
	return widget
}
