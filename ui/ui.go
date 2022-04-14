package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/domain"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Layout struct {
	pages       *tview.Pages
	layout      *tview.Flex
	main        *tview.Flex
	description *tview.TextView
	repos       *RepoWidget
	issues      *IssuesWidget
	favourites  *FavouritesWidget
}

var (
	leftPillIcon  = "█"
	rightPillIcon = "█"
	repoIcon      = ""
)

var (
	view *Layout
	UI   *tview.Application
)

//--------------------------------------------------------------------------------------------------
//  Input handlers
//--------------------------------------------------------------------------------------------------

func appInputHandler(layout *Layout, event *tcell.EventKey) *tcell.EventKey {
	elements := []tview.Primitive{
		layout.main,
		layout.issues.component,
		layout.description,
		layout.repos.component,
	}
	switch event.Key() {
	case tcell.KeyCtrlQ:
		UI.Stop()
	case tcell.KeyTab:
		cycleFocus(UI, elements, false)
	case tcell.KeyBacktab:
		cycleFocus(UI, elements, true)
	}
	return event
}

func sidebarInputHandler(
	event *tcell.EventKey,
	nextTab func(),
	previousTab func(),
) *tcell.EventKey {
	if event.Rune() == 'j' {
		return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
	} else if event.Rune() == 'k' {
		return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
	} else if event.Rune() == 'l' {
		return tcell.NewEventKey(tcell.KeyRight, 'l', tcell.ModNone)
	} else if event.Rune() == 'h' {
		return tcell.NewEventKey(tcell.KeyLeft, 'h', tcell.ModNone)
	} else if event.Key() == tcell.KeyCtrlD {
		view.issues.ScrollDown()
	} else if event.Key() == tcell.KeyCtrlU {
		view.issues.ScrollUp()
	} else if event.Key() == tcell.KeyCtrlN {
		nextTab()
		return nil
	} else if event.Key() == tcell.KeyCtrlP {
		previousTab()
		return nil
	}
	return event
}

func cycleFocus(app *tview.Application, elements []tview.Primitive, reverse bool) {
	for i, el := range elements {
		if !el.HasFocus() {
			continue
		}

		if reverse {
			i--
			if i < 0 {
				i = len(elements) - 1
			}
		} else {
			i++
			i = i % len(elements)
		}

		app.SetFocus(elements[i])
		return
	}
}

// openErrorModal opens a modal with the given error message
func openErrorModal(err error) {
	view.pages.AddAndSwitchToPage("errors", getErrorModal(err, func(_ int, _ string) {
		view.pages.SwitchToPage("main")
	}), true)
}

func getErrorModal(err error, onDone func(idx int, label string)) *tview.Modal {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"OK"}).
		SetDoneFunc(onDone)
	return modal
}

func repositoryEntry(repo domain.Repo) (string, string, bool, func()) {
	name := repo.GetName()
	description := repo.GetDescription()
	showSecondaryText := false
	if name != "" {
		if len(description) > 0 {
			showSecondaryText = true
		}
	}
	return repoIcon + " " + name, description, showSecondaryText, nil
}

// throttledRepoList updates the visible issue details in the issue widget when a user
// has paused over a repository in the list for more than interval time
func throttledRepoList(duration time.Duration) func(*domain.Repository) {
	var timer *time.Timer
	return func(repo *domain.Repository) {
		setRepoDescription(repo)
		if timer != nil {
			timer.Stop()
			timer = nil
		}
		timer = time.AfterFunc(duration, func() {
			view.issues.refreshIssuesList(repo)
		})
	}
}

var updateRepoList = throttledRepoList(time.Millisecond * 200)

func setRepoDescription(repo *domain.Repository) {
	view.description.SetTitle(pad(repo.GetName(), 1)).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorBlue)
	stars := fmt.Sprintf("[red]Stars[white]: 🌟%d", repo.GetStargazerCount())
	issues := fmt.Sprintf("[red]Issues[white]: %d", repo.GetIssueCount())
	url := fmt.Sprintf("[red]URL[white]: [blue::bu]%s", repo.URL)
	prs := fmt.Sprintf("[red]Open PRs[white]: %d", repo.GetPullRequestCount())
	text := strings.Join([]string{repo.GetDescription(), "", stars, issues, prs, url}, "\n")
	view.description.SetText(text)
}

func layoutWidget(context *app.Context) *Layout {
	pages := tview.NewPages()
	description := tview.NewTextView()
	main := tview.NewFlex()
	layout := tview.NewFlex()

	favourites := favouritesWidget(context)
	repos := reposWidget(context)
	issues := issuesWidget(context)
	sidebar := sidebarWidget(context, repos, favourites)

	description.SetDynamicColors(true).SetBorder(true)

	main.SetDirection(tview.FlexRow)
	main.
		AddItem(description, 0, 1, false).
		AddItem(issues.component, 0, 3, false)

	layout.
		AddItem(sidebar.component, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(main, 0, 3, false), 0, 3, false)

	navAdvice := "Cycle through sections using TAB/SHIFT-TAB"
	closeAdvice := "Quit using <C-Q> or <C-C>"
	frame := tview.NewFrame(layout).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(fmt.Sprintf("%s | %s", navAdvice, closeAdvice), false, tview.AlignLeft, tcell.ColorWhite)

	pages.AddPage("main", frame, true, true)

	return &Layout{
		pages:       pages,
		main:        main,
		description: description,
		layout:      layout,
		repos:       repos,
		issues:      issues,
		favourites:  favourites,
	}
}

func Setup(context *app.Context) error {
	UI = tview.NewApplication()
	view = layoutWidget(context)

	UI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return appInputHandler(view, event)
	})
	// Only refresh once the application has been mounted
	go view.repos.Refresh()
	go view.favourites.Refresh()
	if err := UI.SetRoot(view.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
