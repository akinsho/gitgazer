package main

import (
	"akinsho/gogazer/database"
	"fmt"
	"strings"
	"time"

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

var leftPillSeparator, rightPillSeparator = "█", "█"

// openErrorModal opens a modal with the given error message
func openErrorModal(err error) {
	app.QueueUpdateDraw(func() {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.Stop()
			})
		app.SetRoot(modal, true)
	})
}

func saveSelectedRepository(db database.Gazers, i int, s1, s2 string, r rune) {
	_, err := db.Insert(getRepositoryByIndex(i))
	if err != nil {
		openErrorModal(err)
		return
	}
}

func refreshRepositoryList(user string, db *database.Gazers) {
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
	view.repoList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		saveSelectedRepository(*db, i, s1, s2, r)
	})
	app.Draw()
}

// renderLabels for an issue by pulling out the name and using ascii pill characters on either
// side of the name
func renderLabels(labels []*github.Label) string {
	var renderedLabels string
	for _, label := range labels {
		color := "#" + strings.ToUpper(label.GetColor())
		name := strings.ToUpper(label.GetName())
		renderedLabels += fmt.Sprintf(
			"[%s]%s%s[%s]%s",
			color,
			leftPillSeparator,
			name,
			color,
			rightPillSeparator,
		)
	}
	return renderedLabels
}

func refreshIssuesList(repo *github.Repository) {
	issues, err := getRepositoryIssues(repo)
	if err != nil {
		openErrorModal(err)
		return
	}
	view.issuesList.Clear()
	if len(issues) == 0 {
		view.issuesList.AddItem("No issues found", "", 0, nil)
	} else {
		for _, issue := range issues {
			issueNumber := fmt.Sprintf("#%d", issue.GetNumber())
			title := truncateText(issue.GetTitle(), 50)
			view.issuesList.AddItem(
				fmt.Sprintf(
					"%s %s [red](%s)",
					issueNumber,
					title,
					strings.ToUpper(issue.GetState()),
				),
				issue.GetUser().GetLogin()+"  "+renderLabels(issue.Labels),
				0,
				nil,
			)
		}
	}
	app.Draw()
}

func vimInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'j' {
		return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
	} else if event.Rune() == 'k' {
		return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
	} else if event.Rune() == 'l' {
		return tcell.NewEventKey(tcell.KeyRight, 'l', tcell.ModNone)
	} else if event.Rune() == 'h' {
		return tcell.NewEventKey(tcell.KeyLeft, 'h', tcell.ModNone)
	}
	return event
}

func updateRepoList() func(index int, mainText, secondaryText string, shortcut rune) {
	var timer *time.Timer
	return func(index int, mainText, secondaryText string, shortcut rune) {
		repo := getRepositoryByIndex(index)
		if repo == nil {
			return
		}
		title := fmt.Sprintf("%s      🌟%d", repo.GetName(), repo.GetStargazersCount())
		issues := fmt.Sprintf("[red]issue count[white]: %d", repo.GetOpenIssuesCount())
		text := fmt.Sprintf("%s\n%s\n%s", title, repo.GetDescription(), issues)
		view.description.SetText(text)
		if timer != nil {
			timer.Stop()
			timer = nil
		}
		timer = time.AfterFunc(time.Second, func() {
			refreshIssuesList(repo)
		})
	}
}

func getLayout() *tview.Flex {
	view.repoList = tview.NewList()
	view.issuesList = tview.NewList()
	view.description = tview.NewTextView()
	view.main = tview.NewFlex()
	sidebar := tview.NewFlex()

	view.repoList.AddItem("Loading repos...", "", 0, nil)
	view.issuesList.SetBorder(true)

	view.repoList.SetChangedFunc(updateRepoList())

	sidebar.AddItem(view.repoList, 0, 1, true).SetBorder(true).SetTitle("Repositories")
	sidebar.SetInputCapture(vimInputHandler)

	title := textWidget("Go Gazer")
	footer := tview.NewBox().SetBorder(true)

	view.description.SetDynamicColors(true).SetBorder(true)
	title.SetBorder(true)

	view.main.
		AddItem(view.description, 0, 1, false).
		AddItem(view.issuesList, 0, 2, false)

	flex := tview.NewFlex().
		AddItem(sidebar, 0, 1, true).
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
