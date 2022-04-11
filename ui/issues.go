package ui

import (
	"akinsho/gogazer/models"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type IssuesWidget struct {
	component *tview.List
}

func issuesWidget() *IssuesWidget {
	issues := tview.NewList()
	issues.SetSelectedStyle(tcell.StyleDefault.Underline(true)).SetBorder(true)
	return &IssuesWidget{component: issues}
}

func (r *IssuesWidget) refreshIssuesList(repo *models.Repository) {
	r.component.Clear()
	issues := repo.Issues.Nodes
	if len(issues) == 0 {
		r.component.AddItem("No issues found", "", 0, nil)
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
			r.component.AddItem(
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
