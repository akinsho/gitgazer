package ui

import (
	"fmt"

	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type panel struct {
	title  string
	widget Widget
	id     string
}

type TabbedPanelWidget struct {
	currentPanel int
	component    *tview.Flex
	entries      []panel
}

func (s *TabbedPanelWidget) SetCurrentIndex(index int) {
	s.currentPanel = index
}

func (s *TabbedPanelWidget) CurrentItem() Widget {
	return s.entries[s.currentPanel].widget
}

func (s *TabbedPanelWidget) CurrentTextView() TextWidget {
	widget, ok := s.entries[s.currentPanel].widget.(TextWidget)
	if !ok {
		return nil
	}
	return widget
}

func (s *TabbedPanelWidget) OnChange(
	panels []panel,
	pages *tview.Pages,
	sidebar *tview.Flex,
) func() {
	return func() {
		page, _ := pages.GetFrontPage()
		index := findCurrentPageByID(panels, page)
		if index == -1 {
			return
		}
		s.SetCurrentIndex(index)
		e := panels[index]
		sidebar.SetTitle(common.Pad(getPanelTitle(panels, e), 1))
		go handleRefresh(e, sidebar, panels)
	}
}

func handleRefresh(selected panel, tabbedPanel *tview.Flex, panels []panel) {
	err := selected.widget.Refresh()
	UI.QueueUpdateDraw(func() {
		if err != nil {
			openErrorModal(err)
		} else {
			tabbedPanel.SetTitle(common.Pad(getPanelTitle(panels, selected), 1))
			UI.SetFocus(selected.widget.Component())
		}
	})
}

func getListItemCount(c tview.Primitive) int {
	list, ok := c.(*tview.List)
	itemCount := 0
	if ok {
		itemCount = list.GetItemCount()
	}
	return itemCount
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
		view.ActiveDetails().ScrollDown()
	} else if event.Key() == tcell.KeyCtrlU {
		view.ActiveDetails().ScrollUp()
	} else if event.Key() == tcell.KeyCtrlN {
		nextTab()
		return nil
	} else if event.Key() == tcell.KeyCtrlP {
		previousTab()
		return nil
	} else if event.Key() == tcell.KeyCtrlO {
		err := view.ActiveList().Open()
		if err != nil {
			openErrorModal(err)
			return nil
		}
		return nil
	}
	return event
}

func findCurrentPageByID(entries []panel, id string) int {
	for i, entry := range entries {
		if entry.id == id {
			return i
		}
	}
	return -1
}

func findNext(panels *tview.Pages, entries []panel, reverse bool) func() {
	return func() {
		page, _ := panels.GetFrontPage()
		index := findCurrentPageByID(entries, page)
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

func panelWidget(_ *app.Context, focused int, entries []panel) *TabbedPanelWidget {
	tabbedPanel := tview.NewFlex()
	panels := tview.NewPages()
	widget := &TabbedPanelWidget{component: tabbedPanel, entries: entries}
	panels.SetChangedFunc(widget.OnChange(entries, panels, tabbedPanel))

	tabbedPanel.
		SetDirection(tview.FlexRow).
		SetTitleColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft)

	previousTab := findNext(panels, entries, true)
	nextTab := findNext(panels, entries, false)

	for index, panel := range entries {
		panels.AddPage(panel.id, panel.widget.Component(), true, index == focused)
	}

	tabbedPanel.SetBorderPadding(0, 0, 0, 0).
		SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			return sidebarInputHandler(event, nextTab, previousTab)
		})

	tabbedPanel.AddItem(panels, 0, 1, false)

	panels.SwitchToPage(entries[focused].id)

	return widget
}

func getPanelTitle(panels []panel, p panel) string {
	title := ""
	for i, entry := range panels {
		itemCount := getListItemCount(entry.widget.Component())
		count := ""
		if itemCount > 0 {
			count += fmt.Sprintf("(%d)", itemCount)
		}
		t := entry.title
		if entry == p {
			title += tview.Escape(fmt.Sprintf("[%s%s]", t, count))
		} else {
			title += t
		}
		if i < len(panels)-1 {
			title += " - "
		}
	}
	return title
}
