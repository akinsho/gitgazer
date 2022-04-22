package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/domain"
	"akinsho/gitgazer/github"
	"fmt"

	"github.com/rivo/tview"
)

var heartIcon = "❤"

type StarredWidget struct {
	component *tview.List
	context   *app.Context
}

func (s *StarredWidget) Context() *app.Context {
	return s.context
}

func (r *StarredWidget) Component() tview.Primitive {
	var c interface{} = r.component
	t, ok := c.(tview.Primitive)
	if !ok {
		panic("failed to cast to tview.Primitive")
	}
	return t
}

func (r *StarredWidget) IsEmpty() bool {
	return len(r.context.State.Starred) == 0
}

func (r *StarredWidget) removeFavouriteIndicator(i int, repo *domain.Repository) {
	_, secondary := r.component.GetItemText(i)
	main, _, _, _ := repositoryEntry(repo)
	r.component.SetItemText(i, main, secondary)
}

func (r *StarredWidget) SetSelected(i int) {
	r.component.SetCurrentItem(i)
}

func (r *StarredWidget) Refresh() (err error) {
	r.component.Clear()
	starred := r.context.State.Starred
	if len(starred) == 0 {
		r.component.AddItem("Loading repositories...", "", 0, nil)
		UI.Draw()
		starred, err = github.ListStarredRepositories(r.context.Client)
		if err != nil {
			return err
		}
		r.context.SetStarred(starred)
	}
	r.component.Clear()
	if len(starred) == 0 {
		r.component.AddItem("No repositories found", "", 0, nil)
		return
	}

	repos := starred
	if len(repos) > 20 {
		repos = starred[:20]
	}

	for _, repo := range repos {
		main, secondary, showSecondaryText, onSelect := repositoryEntry(repo)
		r.component.AddItem(main, secondary, 0, onSelect).
			ShowSecondaryText(showSecondaryText)
	}
	view.repos.addFavouriteIndicators()
	return
}

// addFavouriteIndicators loops through all repositories and if they have been previously
// liked, adds a heart icon to the end of the name.
func (r *StarredWidget) addFavouriteIndicators() {
	for i := 0; i < r.component.GetItemCount(); i++ {
		go r.addFavouriteIndicator(i)
	}
}

func (r *StarredWidget) addFavouriteIndicator(i int) {
	repo := r.context.GetStarred(i)
	if isFavourite(r.context, repo) {
		main, secondary := r.component.GetItemText(i)
		r.component.SetItemText(i, fmt.Sprintf("%s [hotpink]%s", main, heartIcon), secondary)
	}
}

func (r *StarredWidget) OnChanged(index int, _, _ string, _ rune) {
	repo := r.context.GetStarred(index)
	if repo == nil {
		return
	}
	updateRepositoryList(r.context, repo)
}

// isFavourite checks if the repository is a favourite
// by seeing if the database contains a match by ID
func isFavourite(ctx *app.Context, repo *domain.Repository) bool {
	r, err := github.GetFavouriteByRepositoryID(ctx, repo.ID)
	if err != nil {
		return false
	}
	return r != nil
}

func onRepoSelect(ctx *app.Context, index int, mainText, secondaryText string, _ rune) {
	if !isFavourite(ctx, ctx.State.Selected) {
		err := github.FavouriteSelectedRepo(ctx)
		if err != nil {
			openErrorModal(err)
			return
		}
		go view.repos.addFavouriteIndicator(index)
	} else {
		err := github.UnfavouriteSelected(ctx, index)
		if err != nil {
			openErrorModal(err)
			return
		}
		go view.repos.removeFavouriteIndicator(index, ctx.State.Selected)
	}
}

func starredWidget(ctx *app.Context) *StarredWidget {
	widget := &StarredWidget{context: ctx}
	repos := listWidget(ListOptions{
		onChanged: widget.OnChanged,
		onSelected: func(i int, s1, s2 string, r rune) {
			onRepoSelect(ctx, i, s1, s2, r)
		},
	})
	widget.component = repos
	return widget
}
