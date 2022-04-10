package main

import (
	"akinsho/gogazer/github"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type View struct {
	pages       *tview.Pages
	layout      *tview.Flex
	main        *tview.Flex
	description *tview.TextView
	repos       *tview.List
	issues      *tview.List
	favourites  *tview.List
	sidebarTabs *tview.TextView
}

type Panel struct {
	title     string
	component *tview.List
}

var (
	leftPillIcon  = "█"
	rightPillIcon = "█"
	repoIcon      = ""
)

//--------------------------------------------------------------------------------------------------
//  Input handlers
//--------------------------------------------------------------------------------------------------

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

func refreshRepositoryList() {
	repositories, err := github.FetchRepositories(client)
	if err != nil {
		openErrorModal(err)
		return
	}
	view.repos.Clear()
	if len(repositories) == 0 {
		view.repos.AddItem("No repositories found", "", 0, nil)
	}

	for _, repo := range repositories[:20] {
		name := repo.Name
		description := repo.Description
		if name != "" {
			showDesc := false
			if len(description) > 0 {
				showDesc = true
			}
			view.repos.AddItem(repoIcon+" "+repo.Name, description, 0, nil).
				ShowSecondaryText(showDesc)
		}
	}
	view.repos.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		_, err := databaseConn.Insert(github.GetRepositoryByIndex(i))
		if err != nil {
			openErrorModal(err)
			return
		}

	})
	app.Draw()
	app.SetFocus(view.repos)
}

func refreshIssuesList(repo *github.Repository) {
	view.issues.Clear()
	issues := repo.Issues.Nodes
	if len(issues) == 0 {
		view.issues.AddItem("No issues found", "", 0, nil)
	} else {
		for _, issue := range issues {
			issueNumber := fmt.Sprintf("#%d", issue.GetNumber())
			title := truncateText(issue.GetTitle(), 80)
			author := ""
			if issue.Author != nil && issue.Author.Login != "" {
				author += "[::bu]@" + issue.Author.Login
			}
			issueColor := "green"
			if issue.Closed {
				issueColor = "red"
			}
			view.issues.AddItem(
				fmt.Sprintf(
					"[%s]%s[-:-:-] %s %s - %s",
					issueColor,
					tview.Escape(fmt.Sprintf("[%s]", strings.ToUpper(issue.GetState()))),
					issueNumber,
					title,
					author,
				),
				drawLabels(issue.Labels.Nodes),
				0,
				nil,
			)
		}
	}
	app.Draw()
}

func updateRepoList() func(index int, mainText, secondaryText string, shortcut rune) {
	var timer *time.Timer
	return func(index int, mainText, secondaryText string, shortcut rune) {
		repo := github.GetRepositoryByIndex(index)
		if repo == nil {
			return
		}
		title := fmt.Sprintf("%s      🌟%d", repo.GetName(), repo.GetStargazerCount())
		issues := fmt.Sprintf("[red]Issues[white]: %d", len(repo.Issues.Nodes))
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

// drawLabels for an issue by pulling out the name and using ascii pill characters on either
// side of the name
// @see: https://github.com/rivo/tview/blob/5508f4b00266dbbac1ebf7bd45438fe6030280f4/doc.go#L65-L129
func drawLabels(labels []*github.Label) string {
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

func getLayout() *tview.Pages {
	view.pages = tview.NewPages()
	view.repos = tview.NewList()
	view.issues = tview.NewList()
	view.description = tview.NewTextView()
	view.main = tview.NewFlex()
	view.favourites = tview.NewList()

	sidebar := getSidebar()

	view.repos.AddItem("Loading repos...", "", 0, nil)
	view.issues.SetBorder(true)

	view.repos.SetChangedFunc(updateRepoList())
	view.repos.SetHighlightFullLine(true)
	view.repos.SetSelectedBackgroundColor(tcell.ColorForestGreen)
	view.repos.SetMainTextColor(tcell.ColorForestGreen)
	view.repos.SetMainTextStyle(tcell.StyleDefault.Bold(true))
	view.repos.SetSecondaryTextColor(tcell.ColorDarkGray)

	title := tview.NewTextView()
	title.
		SetTextAlign(tview.AlignCenter).
		SetText("Go gazer")

	view.description.SetDynamicColors(true).SetBorder(true)
	title.SetBorder(true)

	view.main.SetDirection(tview.FlexRow)
	view.main.
		AddItem(view.description, 0, 1, false).
		AddItem(view.issues, 0, 3, false)

	view.layout = tview.NewFlex().
		AddItem(sidebar, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(title, 3, 0, false).
			AddItem(view.main, 0, 3, false), 0, 3, false)

	view.pages.AddPage("main", view.layout, true, true)

	return view.pages
}

func getSidebar() *tview.Flex {
	entries := []Panel{
		{title: "Repositories", component: view.repos},
		{title: "Favourites", component: view.favourites},
	}
	sidebar := tview.NewFlex()
	panels := tview.NewPages()
	view.sidebarTabs = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			id := added[0]
			panels.SwitchToPage(id)
			num, err := strconv.ParseInt(id, 10, 0)
			if err != nil {
				return
			}
			e := entries[num]
			app.SetFocus(e.component)
		})

	previousTab := func() {
		tab, _ := strconv.Atoi(view.sidebarTabs.GetHighlights()[0])
		tab = (tab - 1 + len(entries)) % len(entries)
		view.sidebarTabs.Highlight(strconv.Itoa(tab)).
			ScrollToHighlight()
	}
	nextTab := func() {
		tab, _ := strconv.Atoi(view.sidebarTabs.GetHighlights()[0])
		tab = (tab + 1) % len(entries)
		view.sidebarTabs.Highlight(strconv.Itoa(tab)).
			ScrollToHighlight()
	}

	for index, panel := range entries {
		panels.AddPage(strconv.Itoa(index), panel.component, true, index == 0)
		fmt.Fprintf(view.sidebarTabs, `["%d"][darkcyan]%s[white][""]  `, index, panel.title)
	}

	sidebar.SetBorder(true).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return sidebarInputHandler(event, nextTab, previousTab)
	})

	// view.sidebarTabs.SetBorder(true)
	// view.sidebarTabs.SetBorderAttributes(tcell.AttrUnderline)

	divider := tview.NewTextView()
	_, _, width, _ := view.sidebarTabs.GetRect()
	divider.SetText(strings.Repeat("—", width*2))

	sidebar.SetDirection(tview.FlexRow).
		AddItem(view.sidebarTabs, 1, 1, false).
		AddItem(divider, 1, 0, false).
		AddItem(panels, 0, 1, false)

	sidebar.SetBorderPadding(0, 1, 1, 1)

	view.sidebarTabs.Highlight("0")

	return sidebar
}
