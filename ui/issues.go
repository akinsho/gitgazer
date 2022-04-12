package ui

import (
	"akinsho/gitgazer/models"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/rivo/tview"
)

type IssuesWidget struct {
	component *tview.TextView
}

func issuesWidget() *IssuesWidget {
	issues := tview.NewTextView()
	issues.SetDynamicColors(true).SetBorder(true).SetTitle("Issues").SetBorderPadding(0, 0, 1, 1)
	return &IssuesWidget{component: issues}
}

// drawLabels for an issue by pulling out the name and using ascii pill characters on either
// side of the name
// @see: https://github.com/rivo/tview/blob/5508f4b00266dbbac1ebf7bd45438fe6030280f4/doc.go#L65-L129
func drawLabels(labels []*models.Label) string {
	renderedLabels := []string{}
	for _, label := range labels {
		color := "#" + strings.ToUpper(label.Color)
		left := fmt.Sprintf("[%s]%s", color, leftPillIcon)
		right := fmt.Sprintf("[%s:-:]%s", color, rightPillIcon)
		name := fmt.Sprintf(`[black:%s]%s`, color, strings.ToUpper(label.Name))
		renderedLabels = append(renderedLabels, left+name+right)
	}
	return strings.Join(renderedLabels, " ")
}

// scrollUp scroll the issues widget's text view up from the current position by 1 line
func (r *IssuesWidget) ScrollUp() {
	row, col := r.component.GetScrollOffset()
	r.component.ScrollTo(row-1, col)
}

func (r *IssuesWidget) ScrollDown() {
	row, col := r.component.GetScrollOffset()
	r.component.ScrollTo(row+1, col)
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
			body := getIssueBodyMarkdown(issue)
			previous := r.component.GetText(false)
			list := []string{
				previous,
				fmt.Sprintf(
					"[%s]%s[-:-:bu] %s %s - %s",
					issueColor,
					tview.Escape(fmt.Sprintf("[%s]", strings.ToUpper(issue.GetState()))),
					issueNumber,
					title,
					author,
				),
				body,
				drawLabels(issue.Labels.Nodes),
			}, "\n")
			r.component.SetText(lines).SetTextAlign(tview.AlignLeft).ScrollToBeginning()
		}
	}
	app.Draw()
}

func getIssueBodyMarkdown(issue *models.Issue) string {
	body, err := glamour.Render(issue.Body, "dark")
	if err != nil {
		body = issue.Body
	} else {
		body = tview.TranslateANSI(body)
	}
	return body
}
