package ui

import (
	"akinsho/gitgazer/app"
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type PullRequestsWidget struct {
	component *tview.TextView
	context   *app.Context
}

func pullRequestsWidget(ctx *app.Context) *PullRequestsWidget {
	prs := tview.NewTextView().SetDynamicColors(true).SetWrap(true)
	return &PullRequestsWidget{prs, ctx}
}

func (p *PullRequestsWidget) Component() tview.Primitive {
	var c interface{} = p.component
	t, ok := c.(tview.Primitive)
	if !ok {
		panic("cannot convert to tview.TextView")
	}
	return t
}

func (p *PullRequestsWidget) Refresh() (err error) {
	p.component.Clear()
	if p.context.State.Selected == nil {
		return
	}
	prs := []string{}
	_, _, w, _ := p.component.GetInnerRect()
	hr := createHeader(w)
	pullRequests := p.context.State.Selected.PullRequests.Nodes
	if len(pullRequests) == 0 {
		p.component.SetText("No pull requests").SetTextAlign(tview.AlignCenter)
	} else {
		for _, pr := range pullRequests {
			text := convertToMarkdown(pr.Body)
			stateColor := "green"
			if pr.Closed {
				stateColor = "red"
			}
			status := fmt.Sprintf(" [%s]", stateColor) + tview.Escape(fmt.Sprintf("[%s]", pr.State)) + "[-:-:-]"
			author := ""
			if pr.Author != nil && pr.Author.Login != "" {
				author += "[::bu]@" + pr.Author.Login + "[::-]"
			}
			list := []string{hr, pr.Title + status, hr, author, text}
			prs = append(prs, list...)
		}
		p.component.SetText(strings.Join(prs, "\n")).SetTextAlign(tview.AlignLeft).ScrollToBeginning()
	}
	return
}

// scrollUp scroll the issues widget's text view up from the current position by 1 line
func (r *PullRequestsWidget) ScrollUp() {
	row, col := r.component.GetScrollOffset()
	r.component.ScrollTo(row-1, col)
}

func (r *PullRequestsWidget) ScrollDown() {
	row, col := r.component.GetScrollOffset()
	r.component.ScrollTo(row+1, col)
}

func (p *PullRequestsWidget) IsEmpty() bool {
	panic("not implemented")
}
