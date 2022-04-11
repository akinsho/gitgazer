package ui

import (
	"fmt"
	"strconv"
	"strings"

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
			app.SetFocus(e.component)
		})

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
		fmt.Fprintf(sidebarTabs, `["%d"][darkcyan] %s [white][""]  `, index, panel.title)
	}

	sidebar.SetBorder(true).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return vimMotionInputHandler(event, nextTab, previousTab)
	})

	divider := tview.NewTextView()
	_, _, width, _ := sidebarTabs.GetRect()
	divider.SetText(strings.Repeat("—", width*2))

	sidebar.SetDirection(tview.FlexRow).
		AddItem(sidebarTabs, 1, 1, false).
		AddItem(divider, 1, 0, false).
		AddItem(panels, 0, 1, false)

	sidebar.SetBorderPadding(0, 1, 1, 1)

	sidebarTabs.Highlight("0")

	return &SidebarWidget{component: sidebar}
}
