package ui

import (
	"akinsho/gogazer/database"
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
	favourites  *tview.List
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

func vimMotionInputHandler(
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
	refreshFavouritesList(favourites, err)
}

// refreshFavouritesList fetches all saved repositories from the database and
// adds them to the View.favourites list.
func refreshFavouritesList(favourites []models.FavouriteRepository, err error) {
	if view.favourites.GetItemCount() > 0 {
		view.favourites.Clear()
	}
	if err != nil {
		openErrorModal(err)
		return
	}
	if len(favourites) == 0 {
		view.favourites.AddItem("No favourites found", "", 0, nil)
	}

	repos := favourites
	if len(repos) > 20 {
		repos = favourites[:20]
	}

	for _, repo := range favourites {
		main, secondary, showSecondaryText, onSelect := repositoryEntry(&repo)
		view.favourites.AddItem(main, secondary, 0, onSelect).
			ShowSecondaryText(showSecondaryText)
	}
	app.Draw()
}

// isFavourited checks if the repository is a favourite
// by seeing if the database contains a match by ID
func isFavourited(repo *models.Repository) bool {
	r, err := database.GetFavouriteByRepoID(repo.ID)
	if err != nil {
		return false
	}
	return r != nil
}

func setRepoDescription(repo *models.Repository) {
	title := fmt.Sprintf("%s      🌟%d", repo.GetName(), repo.GetStargazerCount())
	issues := fmt.Sprintf("[red]Issues[white]: %d", repo.GetIssueCount())
	text := fmt.Sprintf("%s\n%s\n%s", title, repo.GetDescription(), issues)
	view.description.SetText(text)
}

func updateFavouriteChange(index int, mainText, secondaryText string, shortcut rune) {
	repo := github.GetFavouriteRepositoryByIndex(index)
	if repo == nil {
		return
	}
	// setRepoDescription(repo)
}

// drawLabels for an issue by pulling out the name and using ascii pill characters on either
// side of the name
// @see: https://github.com/rivo/tview/blob/5508f4b00266dbbac1ebf7bd45438fe6030280f4/doc.go#L65-L129
func drawLabels(labels []*models.Label) string {
	var renderedLabels string
	for _, label := range labels {
		color := "#" + strings.ToUpper(label.Color)
		left := fmt.Sprintf("[%s]%s", color, leftPillIcon)
		right := fmt.Sprintf("[%s:-:]%s", color, rightPillIcon)
		name := fmt.Sprintf(`[black:%s]%s`, color, strings.ToUpper(label.Name))
		renderedLabels += left + name + right
	}
	return renderedLabels
}

func layoutWidget() *Layout {
	pages := tview.NewPages()
	description := tview.NewTextView()
	main := tview.NewFlex()
	favourites := tview.NewList()
	layout := tview.NewFlex()

	repos := reposWidget()
	issues := issuesWidget()
	sidebar := sidebarWidget(repos.component, favourites)

	favourites.SetChangedFunc(updateFavouriteChange)

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
	view = layoutWidget()
	app = tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return appInputHandler(view, event)
	})
	go fetchStarredRepositories(client)
	go fetchFavouriteRepositories()
	if err := app.SetRoot(view.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
