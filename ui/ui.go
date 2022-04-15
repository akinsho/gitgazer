package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"akinsho/gitgazer/domain"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	leftPillIcon  = "█"
	rightPillIcon = "█"
	repoIcon      = ""
	headerChar    = "─"
)

var (
	view *Layout
	UI   *tview.Application
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

func (l *Layout) ActiveList() *tview.List {
	if view.favourites.component.HasFocus() {
		return l.favourites.component
	}
	return l.repos.component
}

//--------------------------------------------------------------------------------------------------
//  Input handlers
//--------------------------------------------------------------------------------------------------

func appInputHandler(layout *Layout, event *tcell.EventKey) *tcell.EventKey {
	elements := []tview.Primitive{
		layout.ActiveList(),
		layout.description,
		layout.issues.component,
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

func createHeader(width int) string {
	return strings.Repeat(headerChar, width)
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

// throttledListUpdate updates the visible issue details in the issue widget when a user
// has paused over a repository in the list for more than interval time
func throttledListUpdate(duration time.Duration) func(*domain.Repository) {
	var timer *time.Timer
	return func(repo *domain.Repository) {
		setRepoDescription(repo)
		if timer != nil {
			timer.Stop()
			timer = nil
		}
		// TODO: this setter only exists because repo cannot be passed to
		// issues directly as the Widget type does not expect Refresh to have an argument
		// and the argument cannot be sufficiently generic anyway.
		timer = time.AfterFunc(duration, func() {
			view.issues.SetRepo(repo)
			view.issues.Refresh()
		})
	}
}

var updateRepositoryList = throttledListUpdate(time.Millisecond * 200)

func setRepoDescription(repo *domain.Repository) {
	view.description.SetTitle(common.Pad(repo.GetName(), 1)).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorBlue)
	stars := fmt.Sprintf("[red]Stars[white]: 🌟%d", repo.GetStargazerCount())
	issues := fmt.Sprintf("[red]Issues[white]: %d", repo.GetIssueCount())
	url := fmt.Sprintf("[red]URL[white]: [blue::bu]%s", repo.URL)
	prs := fmt.Sprintf("[red]Open PRs[white]: %d", repo.GetPullRequestCount())
	text := strings.Join([]string{repo.GetDescription(), "", stars, issues, prs, url}, "\n")
	view.description.SetText(text)
}

func helpWidget() *tview.TextView {
	navAdvice := "Cycle through sections using [::b]TAB/SHIFT-TAB[::-]"
	closeAdvice := "Quit using [::b]<C-Q>[::-] or [::b]<C-C>[::-]"
	listNavAdvice := "Navigate through the list using [::b]j/k[::-]"
	listNavScrollAdvice := "Scroll through the issues list using [::b]C-D/C-U[::-]"
	helpText := strings.Join([]string{
		navAdvice,
		closeAdvice,
		listNavAdvice,
		listNavScrollAdvice,
	}, " | ")
	help := tview.NewTextView().SetText(helpText).SetDynamicColors(true)
	help.SetBorder(true)
	return help
}

type ListOptions struct {
	onSelected func(int, string, string, rune)
	onChanged  func(int, string, string, rune)
}

func listWidget(opts ListOptions) *tview.List {
	list := tview.NewList()
	list.SetChangedFunc(opts.onChanged).
		SetSelectedFunc(opts.onSelected).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorRebeccaPurple).
		SetMainTextColor(tcell.ColorForestGreen).
		SetMainTextStyle(tcell.StyleDefault.Bold(true)).
		SetBorderPadding(0, 0, 1, 1)
	return list
}

func repositoryPanelWidget(
	favourites *FavouritesWidget,
	context *app.Context,
	repos *RepoWidget,
) *SidebarWidget {
	leftSidebarFocused := 0
	if !favourites.IsEmpty() {
		leftSidebarFocused = 1
	}
	sidebar := panelWidget(context, leftSidebarFocused, []panel{
		{id: "starred", title: "Starred", widget: repos},
		{id: "favourites", title: "Favourites", widget: favourites},
	})
	return sidebar
}

// setupTheme sets up the theme for the application
// which can be derived from the app's config
// TODO: pull colour values from config
func setupTheme(_ *app.Config) {
	theme := tview.Theme{
		TitleColor:                  tcell.ColorBlue,
		MoreContrastBackgroundColor: tcell.ColorGray,
	}
	tview.Styles = theme
}

func layoutWidget(ctx *app.Context) *Layout {
	pages := tview.NewPages()
	description := tview.NewTextView()
	main := tview.NewFlex()
	frame := tview.NewFlex().SetDirection(tview.FlexRow)
	layout := tview.NewFlex()

	favourites := favouritesWidget(ctx)
	repos := reposWidget(ctx)
	issues := issuesWidget(ctx)

	sidebar := repositoryPanelWidget(favourites, ctx, repos)
	issuesPanel := panelWidget(ctx, 0, []panel{{title: "Issues", widget: issues}})

	description.SetDynamicColors(true).SetBorder(true)

	main.SetDirection(tview.FlexRow)
	main.
		AddItem(description, 0, 1, false).
		AddItem(issuesPanel.component, 0, 3, false)

	layout.
		AddItem(sidebar.component, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(main, 0, 3, false), 0, 3, false)

	frame.AddItem(layout, 0, 1, false).AddItem(helpWidget(), 3, 0, false)

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

func focusActiveList() {
	if !view.favourites.IsEmpty() {
		view.favourites.Refresh()
		UI.SetFocus(view.favourites.component)
	} else {
		view.repos.Refresh()
		UI.SetFocus(view.repos.component)
	}
}

func Setup(context *app.Context) error {
	setupTheme(context.Config)
	UI = tview.NewApplication()
	view = layoutWidget(context)

	UI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return appInputHandler(view, event)
	})
	// Only focus once the application has been mounted
	go focusActiveList()
	if err := UI.SetRoot(view.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
