package ui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SidebarWidget struct {
	component *tview.Flex
}

type panel struct {
	title     string
	component *tview.List
}

func sidebarWidget(repos *tview.List, favourites *tview.List) *SidebarWidget {
	entries := []panel{
		{title: "Repositories", component: repos},
		{title: "Favourites", component: favourites},
	}
	sidebar := tview.NewFlex()
	panels := tview.NewPages()
	sidebarTabs := tview.NewTextView().
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
			sidebar.SetTitle(e.title).
				SetTitleColor(tcell.ColorBlue).
				SetTitleAlign(tview.AlignLeft)
			app.SetFocus(e.component)
		})
	sidebarTabs.SetBorder(true)

	previousTab := func() {
		tab, _ := strconv.Atoi(sidebarTabs.GetHighlights()[0])
		tab = (tab - 1 + len(entries)) % len(entries)
		sidebarTabs.Highlight(strconv.Itoa(tab)).
			ScrollToHighlight()
	}
	nextTab := func() {
		tab, _ := strconv.Atoi(sidebarTabs.GetHighlights()[0])
		tab = (tab + 1) % len(entries)
		sidebarTabs.Highlight(strconv.Itoa(tab)).
			ScrollToHighlight()
	}

	for index, panel := range entries {
		panels.AddPage(strconv.Itoa(index), panel.component, true, index == 0)
		fmt.Fprintf(sidebarTabs, `["%d"][darkcyan]%s[white][""]`, index, pad(panel.title, 1))
		if index == 0 {
			fmt.Fprintf(sidebarTabs, "|")
		}
	}

	sidebar.SetBorderPadding(0, 0, 0, 0).
		SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			return sidebarInputHandler(event, nextTab, previousTab)
		})

	sidebar.SetDirection(tview.FlexRow).
		AddItem(panels, 0, 1, false).
		AddItem(sidebarTabs, 3, 0, false)

	sidebarTabs.Highlight("0")

	return &SidebarWidget{component: sidebar}
}