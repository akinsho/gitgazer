package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"akinsho/gitgazer/domain"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
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
	details     *TabbedPanelWidget
	repos       *StarredWidget
	issues      *IssuesWidget
	prs         *PullRequestsWidget
	sidebar     *TabbedPanelWidget
	favourites  *FavouritesWidget
}

func (l *Layout) ActiveList() ListWidget {
	if view.favourites.component.HasFocus() {
		return l.favourites
	} else if view.repos.component.HasFocus() {
		return l.repos
	} else {
		return view.sidebar.CurrentItem().(ListWidget)
	}
}

func (l *Layout) ActiveDetails() TextWidget {
	if view.issues.component.HasFocus() {
		return view.issues
	} else if view.prs.component.HasFocus() {
		return view.prs
	} else {
		return view.details.CurrentTextView()
	}
}

//--------------------------------------------------------------------------------------------------
//  Input handlers
//--------------------------------------------------------------------------------------------------

func appInputHandler(layout *Layout, event *tcell.EventKey) *tcell.EventKey {
	elements := []tview.Primitive{
		layout.ActiveList().Component(),
		layout.description,
		layout.ActiveDetails().Component(),
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

func openErrorModal(err error) {
	current := UI.GetFocus()
	modal := getErrorModal(err, "Sorry! looks like something went wrong", func(_ int, _ string) {
		view.pages.SwitchToPage("main")
		UI.SetFocus(current)
	})
	view.pages.AddPage("errors", modal, true, true)
}

func getErrorModal(err error, title string, onDone func(int, string)) *tview.Modal {
	lines := strings.Join([]string{title, "message: " + err.Error()}, "\n")
	modal := tview.NewModal().
		SetText(lines).
		AddButtons([]string{"OK"}).
		SetDoneFunc(onDone)
	return modal
}

func createHeader(width int) string {
	return strings.Repeat(headerChar, width)
}

func convertToMarkdown(body string) string {
	body, err := glamour.Render(body, "dark")
	if err != nil {
		return body
	}
	return tview.TranslateANSI(body)
}

func repositoryEntry(repo domain.Repo) (string, string, bool, func()) {
	name := repo.GetName()
	description := repo.GetDescription()
	showSecondaryText := false
	if len(strings.TrimSpace(description)) > 0 {
		showSecondaryText = true
	}
	return repoIcon + " " + name, description, showSecondaryText, nil
}

func updateRepositoryList(ctx *app.Context, repo *domain.Repository) {
	setRepoDescription(repo)
	ctx.SetSelected(repo)
	err := view.ActiveDetails().Refresh()
	if err != nil {
		openErrorModal(err)
	}
}

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
		SetSecondaryTextColor(tcell.ColorDarkGrey).
		SetSelectedBackgroundColor(tcell.ColorRebeccaPurple).
		SetMainTextColor(tcell.ColorForestGreen).
		SetMainTextStyle(tcell.StyleDefault.Bold(true)).
		SetBorderPadding(0, 0, 1, 1)
	return list
}

func repositoryPanelWidget(
	context *app.Context,
	favourites *FavouritesWidget,
	starred *StarredWidget,
) *TabbedPanelWidget {
	focused := 0
	if !favourites.IsEmpty() {
		focused = 1
	}
	sidebar := panelWidget(context, focused, []panel{
		{id: domain.StarredRepositoriesPanel.String(), title: "Starred", widget: starred},
		{id: domain.FavouriteRepositoriesPanel.String(), title: "Favourites", widget: favourites},
	})
	return sidebar
}

func repositoryDetailsPanelWidget(
	ctx *app.Context,
	issues *IssuesWidget,
	prs *PullRequestsWidget,
) *TabbedPanelWidget {
	focused := 0
	preferred := ctx.Config.UserConfig.Panels.Details.Preferred
	if preferred == domain.PullRequestPanel {
		focused = 1
	}
	return panelWidget(nil, focused, []panel{
		{id: domain.IssuesPanel.String(), title: "Issues", widget: issues},
		{id: domain.PullRequestPanel.String(), title: "PRs", widget: prs},
	})
}

// TODO: pull colour values from config
// setupTheme sets up the theme for the application which can be derived from the app's config
func setupTheme(_ *app.Config) {
	theme := tview.Theme{
		PrimitiveBackgroundColor:    tview.Styles.PrimitiveBackgroundColor,
		ContrastBackgroundColor:     tcell.ColorDimGray,
		MoreContrastBackgroundColor: tcell.ColorRebeccaPurple,
		BorderColor:                 tview.Styles.BorderColor,
		TitleColor:                  tcell.ColorBlue,
		GraphicsColor:               tview.Styles.GraphicsColor,
		PrimaryTextColor:            tview.Styles.PrimaryTextColor,
		SecondaryTextColor:          tview.Styles.SecondaryTextColor,
		TertiaryTextColor:           tview.Styles.TertiaryTextColor,
		InverseTextColor:            tview.Styles.InverseTextColor,
		ContrastSecondaryTextColor:  tview.Styles.ContrastSecondaryTextColor,
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
	prs := pullRequestsWidget(ctx)

	sidebar := repositoryPanelWidget(ctx, favourites, repos)
	details := repositoryDetailsPanelWidget(ctx, issues, prs)

	description.SetDynamicColors(true).SetBorder(true)

	main.SetDirection(tview.FlexRow)
	main.
		AddItem(description, 0, 1, false).
		AddItem(details.component, 0, 3, false)

	layout.
		AddItem(sidebar.component, 0, 1, false).
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
		sidebar:     sidebar,
		details:     details,
		prs:         prs,
		favourites:  favourites,
	}
}

func Setup(context *app.Context) error {
	setupTheme(context.Config)
	UI = tview.NewApplication()
	view = layoutWidget(context)

	UI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return appInputHandler(view, event)
	})
	if err := UI.SetRoot(view.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
