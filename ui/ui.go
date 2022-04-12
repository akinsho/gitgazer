package ui

import (
	"akinsho/gogazer/github"
	"akinsho/gogazer/models"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"
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
	app  *tview.Application
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
		app.Stop()
	case tcell.KeyTab:
		cycleFocus(app, elements, false)
	case tcell.KeyBacktab:
		cycleFocus(app, elements, true)
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
	view.pages.AddAndSwitchToPage("errors", getErrorModal(err, func(idx int, label string) {
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

func repositoryEntry(repo models.Repo) (string, string, bool, func()) {
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

func fetchStarredRepositories(client *githubv4.Client) {
	repositories, err := github.ListStarredRepositories(client)
	view.repos.refreshRepositoryList(repositories, err)
}

func fetchFavouriteRepositories() {
	favourites, err := github.ListSavedFavourites()
	view.favourites.refreshFavouritesList(favourites, err)
}

func setRepoDescription(repo *models.Repository) {
	view.description.SetTitle(repo.GetName()).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorBlue)
	stars := fmt.Sprintf("[red]Stars[white]: 🌟%d", repo.GetStargazerCount())
	issues := fmt.Sprintf("[red]Issues[white]: %d", repo.GetIssueCount())
	url := fmt.Sprintf("[red]URL[white]: [blue::bu]%s", repo.URL)
	prs := fmt.Sprintf("[red]Open PRs[white]: %d", repo.GetPullRequestCount())
	text := strings.Join([]string{repo.GetDescription(), "", stars, issues, prs, url}, "\n")
	view.description.SetText(text)
}

func layoutWidget() *Layout {
	pages := tview.NewPages()
	description := tview.NewTextView()
	main := tview.NewFlex()
	layout := tview.NewFlex()

	favourites := favouritesWidget()
	repos := reposWidget()
	issues := issuesWidget()
	sidebar := sidebarWidget(repos.component, favourites.component)

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

func Setup(client *githubv4.Client) error {
	app = tview.NewApplication()
	view = layoutWidget()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return appInputHandler(view, event)
	})
	// Only refresh once the application has been mounted
	go fetchStarredRepositories(client)
	go fetchFavouriteRepositories()
	if err := app.SetRoot(view.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
