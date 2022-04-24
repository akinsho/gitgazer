package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"akinsho/gitgazer/domain"
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type IssuesWidget struct {
	component *tview.TextView
	context   *app.Context
}

func (r *IssuesWidget) Open() error {
	panic("not implemented") // TODO: Implement
}

func (i *IssuesWidget) Context() *app.Context {
	return i.context
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

func (r *IssuesWidget) IsEmpty() bool {
	selected := r.context.State.Selected
	if selected == nil || len(selected.Issues.Nodes) == 0 {
		return true
	}
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

func (r *IssuesWidget) Refresh() (err error) {
	r.component.Clear()
	if r.context.State.Selected == nil {
		return
	}
	issues := r.context.State.Selected.Issues.Nodes
	if len(issues) == 0 {
		r.component.SetText("No issues found").SetTextAlign(tview.AlignCenter)
	} else {
		_, _, width, _ := r.component.GetInnerRect()
		header := createHeader(width)
		lines := []string{}
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
			body := convertToMarkdown(issue.Body)
			lines = append(
				lines,
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
			)
		}
		lines = removeBlankLines(lines)
		r.component.SetText(strings.Join(lines, "\n")).SetTextAlign(tview.AlignLeft).ScrollToBeginning()
	}
	return
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

func issuesWidget(ctx *app.Context) *IssuesWidget {
	issues := tview.NewTextView().SetDynamicColors(true)
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
