package ui

import (
	"akinsho/gitgazer/app"

	"github.com/rivo/tview"
)

type PullRequestsWidget struct {
	component *tview.TextView
}

func pullRequestsWidget(ctx *app.Context) *PullRequestsWidget {
	prs := tview.NewTextView()
	return &PullRequestsWidget{prs}
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
	return
}

func (p *PullRequestsWidget) IsEmpty() bool {
	panic("not implemented")
}
