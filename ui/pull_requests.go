package ui

import (
	"akinsho/gitgazer/app"
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

func (p *PullRequestsWidget) Refresh() {
	p.component.Clear()
	prs := []string{}
	_, _, w, _ := p.component.GetInnerRect()
	hr := createHeader(w)
	for _, pr := range p.context.State.Selected.PullRequests.Nodes {
		text := convertToMarkdown(pr.Body)
		list := []string{hr, pr.Title, hr, text}
		prs = append(prs, list...)
	}
	p.component.SetText(strings.Join(prs, "\n"))
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