package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"fmt"
	"strconv"

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
}

func onPanelChange(panels []panel, pages *tview.Pages, sidebar *tview.Flex) func() {
	return func() {
		num := getCurrentPage(pages)
		if num == nil {
			return
		}
		e := panels[*num]
		title := getSidebarTitle(panels, e)
		sidebar.SetTitle(common.Pad(title, 1)).
			SetTitleColor(tcell.ColorBlue).
			SetTitleAlign(tview.AlignLeft)
		go e.widget.Refresh()
		UI.SetFocus(e.widget.Component())
	}
}

func getCurrentPage(pages *tview.Pages) *int64 {
	name, _ := pages.GetFrontPage()
	num, err := strconv.ParseInt(name, 10, 0)
	if err != nil {
		return nil
	}
	return &num
}

func sidebarWidget(
	ctx *app.Context,
	repos *RepoWidget,
	favourites *FavouritesWidget,
) *SidebarWidget {
	entries := []panel{
		{title: "Starred", widget: repos},
		{title: "Favourites", widget: favourites},
	}
	sidebar := tview.NewFlex()
	panels := tview.NewPages()
	panels.SetChangedFunc(onPanelChange(entries, panels, sidebar))

	previousTab := func() {
		num := getCurrentPage(panels)
		if num == nil {
			return
		}
		tab := int(*num)
		tab = (tab - 1 + len(entries)) % len(entries)
		panels.SwitchToPage(strconv.Itoa(tab))
	}
	nextTab := func() {
		num := getCurrentPage(panels)
		if num == nil {
			return
		}
		tab := int(*num)
		tab = (tab + 1) % len(entries)
		panels.SwitchToPage(strconv.Itoa(tab))
	}

	// We want to start on the favourites page since we like those the best not the super
	// long list of all repos we've starred.
	focused := 0
	if !favourites.IsEmpty() {
		focused = 1
	}

	for index, panel := range entries {
		panels.AddPage(strconv.Itoa(index), panel.widget.Component(), true, index == focused)
	}

	sidebar.SetBorderPadding(0, 0, 0, 0).
		SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			return sidebarInputHandler(event, nextTab, previousTab)
		})

	sidebar.SetDirection(tview.FlexRow).
		AddItem(panels, 0, 1, false)

	panels.SwitchToPage(strconv.Itoa(focused))

	return &SidebarWidget{component: sidebar}
}

func getSidebarTitle(entries []panel, e panel) string {
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
