package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SidebarWidget struct {
	currentPanel int
	component    *tview.Flex
}

type panel struct {
	title  string
	widget Widget
	id     string
}

func findCurrenPageByID(entries []panel, id string) int {
	for i, entry := range entries {
		if entry.id == id {
			return i
		}
	}
	return -1
}

func onPanelChange(panels []panel, pages *tview.Pages, sidebar *tview.Flex) func() {
	return func() {
		page, _ := pages.GetFrontPage()
		index := findCurrenPageByID(panels, page)
		if index == -1 {
			return
		}
		e := panels[index]
		title := getPanelTitle(panels, e)
		sidebar.SetTitle(common.Pad(title, 1)).
			SetTitleColor(tcell.ColorBlue).
			SetTitleAlign(tview.AlignLeft)
		go e.widget.Refresh()
		UI.SetFocus(e.widget.Component())
	}
}

func findNext(panels *tview.Pages, entries []panel, reverse bool) func() {
	return func() {
		page, _ := panels.GetFrontPage()
		index := findCurrenPageByID(entries, page)
		if index == -1 {
			return
		}
		tab := index
		if reverse {
			tab = (tab - 1 + len(entries)) % len(entries)
		} else {
			tab = (tab + 1) % len(entries)
		}
		panels.SwitchToPage(entries[tab].id)
	}
}

func panelWidget(ctx *app.Context, focused int, entries []panel) *SidebarWidget {
	sidebar := tview.NewFlex()
	panels := tview.NewPages()
	panels.SetChangedFunc(onPanelChange(entries, panels, sidebar))

	previousTab := findNext(panels, entries, true)
	nextTab := findNext(panels, entries, false)

	for index, panel := range entries {
		panels.AddPage(panel.id, panel.widget.Component(), true, index == focused)
	}

	sidebar.SetBorderPadding(0, 0, 0, 0).
		SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			return sidebarInputHandler(event, nextTab, previousTab)
		})

	sidebar.SetDirection(tview.FlexRow).
		AddItem(panels, 0, 1, false)

	panels.SwitchToPage(entries[focused].id)

	return &SidebarWidget{component: sidebar}
}

func getPanelTitle(entries []panel, e panel) string {
	title := ""
	for i, entry := range entries {
		t := entry.title
		if entry == e {
			title += fmt.Sprintf("(%s)", t)
		} else {
			title += t
		}
		if i < len(entries)-1 {
			title += " - "
		}
	}
	return title
}
