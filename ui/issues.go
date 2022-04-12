package ui

import (
	"akinsho/gogazer/models"
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type IssuesWidget struct {
	component *tview.TextView
}

func issuesWidget() *IssuesWidget {
	issues := tview.NewTextView()
	issues.SetDynamicColors(true).SetBorder(true).SetTitle("Issues")
	return &IssuesWidget{component: issues}
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

// formatIssueBody splits the body of an issue at a certain character count into
// a list of lines that are then rejoined by newlines
func formatIssueBody(issue *models.Issue) string {
	if issue.Body == "" {
		return ""
	}
	body := strings.Split(issue.Body, "\n")
	var formattedBody string
	for _, line := range body {
		formattedBody += truncateText(line, 80, false) + "\n"
	}
	return formattedBody
}

func (r *IssuesWidget) refreshIssuesList(repo *models.Repository) {
	r.component.Clear()
	issues := repo.Issues.Nodes
	if len(issues) == 0 {
		r.component.SetText("No issues found").SetTextAlign(tview.AlignCenter)
	} else {
		for _, issue := range issues {
			issueNumber := fmt.Sprintf("#%d", issue.GetNumber())
			title := truncateText(issue.GetTitle(), 80, true)
			author := ""
			if issue.Author != nil && issue.Author.Login != "" {
				author += "[::bu]@" + issue.Author.Login + "[::-]"
			}
			issueColor := "green"
			if issue.Closed {
				issueColor = "red"
			}
			previous := r.component.GetText(false)
			r.component.SetText(
				strings.Join([]string{
					previous,
					fmt.Sprintf(
						"[%s]%s[-:-:-] %s %s - %s",
						issueColor,
						tview.Escape(fmt.Sprintf("[%s]", strings.ToUpper(issue.GetState()))),
						issueNumber,
						title,
						author,
					),
					formatIssueBody(issue),
					drawLabels(issue.Labels.Nodes),
				}, "\n"),
			)
		}
	}
	app.Draw()
}
