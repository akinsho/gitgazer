package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Create a struct representing a repository.
type repo struct {
	name        string
	description string
}

var view = View{}

type View struct {
	repoList *tview.List
}

func textElement(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

func boxWidget() *tview.Box {
	return tview.NewBox().SetBackgroundColor(tcell.Color100)
}

func refreshRepositoryList() {
	data := []repo{
		{"Repo 1", "This is a description of the first repo"},
		{"Repo 2", "This is a description of the second repo"},
		{"Repo 3", "This is a description of the third repo"},
	}
	view.repoList.Clear()
	for _, repo := range data {
		view.repoList.AddItem(repo.name, repo.description, 0, nil)
	}
}

func layoutGrid() *tview.Flex {
	view.repoList = tview.NewList()
	sidebar := tview.NewFlex()

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

func main() {
	grid := layoutGrid()
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.Stop()
		}
		return event
	})
	refreshRepositoryList()
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		log.Panicln(err)
	}
}
