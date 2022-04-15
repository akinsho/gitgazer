package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"akinsho/gitgazer/domain"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/rivo/tview"
)

type IssuesWidget struct {
	component *tview.TextView
	context   *app.Context
	repo      *domain.Repository
}

func issuesWidget(ctx *app.Context) *IssuesWidget {
	issues := tview.NewTextView()
	issues.SetDynamicColors(true)
	return &IssuesWidget{component: issues, context: ctx}
}

// drawLabels for an issue by pulling out the name and using ascii pill characters on either
// side of the name
// @see: https://github.com/rivo/tview/blob/5508f4b00266dbbac1ebf7bd45438fe6030280f4/doc.go#L65-L129
func drawLabels(labels []*domain.Label) string {
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

func (r *IssuesWidget) SetRepo(repo *domain.Repository) {
	r.repo = repo
}

func (r *IssuesWidget) IsEmpty() bool {
	return false
}

func (r *IssuesWidget) Component() tview.Primitive {
	var c interface{} = r.component
	t, ok := c.(tview.Primitive)
	if !ok {
		panic("failed to cast to tview.Primitive")
	}
	return t
}

func (r *IssuesWidget) Refresh() {
	r.component.Clear()
	if r.repo == nil {
		return
	}
	issues := r.repo.Issues.Nodes
	if len(issues) == 0 {
		r.component.SetText("No issues found").SetTextAlign(tview.AlignCenter)
	} else {
		_, _, width, _ := r.component.GetInnerRect()
		header := createHeader(width)
		for _, issue := range issues {
			issueNumber := fmt.Sprintf("#%d", issue.GetNumber())
			title := common.TruncateText(issue.GetTitle(), 80, true)
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
				header,
				fmt.Sprintf(
					"[%s]%s[-::bu] %s %s - %s[-:-:-]",
					issueColor,
					tview.Escape(fmt.Sprintf("[%s]", strings.ToUpper(issue.GetState()))),
					issueNumber,
					title,
					author,
				),
				header,
				fmt.Sprintf("Created at: %s", issue.CreatedAt.Format("02-01-2006 15:04:05")),
				drawLabels(issue.Labels.Nodes),
				body,
			}
			lines := removeBlankLines(list)
			r.component.SetText(strings.Join(lines, "\n")).SetTextAlign(tview.AlignLeft).ScrollToBeginning()
		}
	}
	UI.Draw()
}

func removeBlankLines(lines []string) []string {
	var filtered []string
	for _, line := range lines {
		if line != "" {
			filtered = append(filtered, line)
		}
	}
	return filtered
}

func getIssueBodyMarkdown(issue *domain.Issue) string {
	body, err := glamour.Render(issue.Body, "dark")
	if err != nil {
		body = issue.Body
	} else {
		body = tview.TranslateANSI(body)
	}
	return body
}
