package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/google/go-github/v43/github"
	"github.com/rivo/tview"
)

type View struct {
	main        *tview.Flex
	description *tview.TextView
	repoList    *tview.List
	issuesList  *tview.List
}

// openErrorModal opens a modal with the given error message
func openErrorModal(err error) {
	app.QueueUpdateDraw(func() {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				// app.Stop()
			})
		app.SetRoot(modal, true)
	})
}

func refreshRepositoryList(user string) {
	repositories, err := fetchRepositories(user)
	if err != nil {
		openErrorModal(err)
		return
	}
	view.repoList.Clear()
	if len(repositories) == 0 {
		view.repoList.AddItem("No repositories found", "", 0, nil)
	}

	for _, repo := range repositories[:20] {
		name := repo.GetName()
		description := repo.GetDescription()
		if name != "" {
			showDesc := false
			if len(description) > 0 {
				showDesc = true
			}
			view.repoList.AddItem(repo.GetName(), description, 0, nil).ShowSecondaryText(showDesc)
		}
		view.repoList.SetSelectedBackgroundColor(tcell.Color101)
	}
	app.Draw()
}

func refreshIssuesList(repo *github.Repository) {
	issues, err := getRepositoryIssues(repo)
	if err != nil {
		openErrorModal(err)
		return
	}
	view.issuesList.Clear()
	for _, issue := range issues {
		view.issuesList.AddItem(
			fmt.Sprintf("%s [red](%s)", issue.GetTitle(), strings.ToUpper(issue.GetState())),
			fmt.Sprintf("#%d", issue.GetNumber()),
			0,
			nil,
		)
	}
	app.Draw()
}

func getLayout() *tview.Flex {
	view.repoList = tview.NewList()
	view.issuesList = tview.NewList()
	view.description = tview.NewTextView()
	view.main = tview.NewFlex()
	sidebar := tview.NewFlex()

	view.repoList.AddItem("Loading repos...", "", 0, nil)
	view.issuesList.SetBorder(true)

	view.repoList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		repo := getRepositoryByIndex(index)
		if repo == nil {
			return
		}
		text := fmt.Sprintf(
			"%s\n%s\n[red]issue count[white]: %d",
			repo.GetName(),
			repo.GetDescription(),
			repo.GetOpenIssuesCount(),
		)
		view.description.SetText(text)
		go refreshIssuesList(repo)
	})

	sidebar.AddItem(view.repoList, 0, 1, true).SetBorder(true).SetTitle("Repositories")

	title := textWidget("Go Gazer")
	footer := tview.NewBox().SetBorder(true)

	view.description.SetDynamicColors(true).SetBorder(true)
	title.SetBorder(true)

	view.main.
		AddItem(view.description, 0, 1, false).
		AddItem(view.issuesList, 0, 2, false)

	flex := tview.NewFlex().
		AddItem(sidebar, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(title, 3, 0, false).
			AddItem(view.main, 0, 3, false).
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
