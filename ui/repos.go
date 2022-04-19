package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/domain"
	"akinsho/gitgazer/github"
	"fmt"

	"github.com/rivo/tview"
)

type RepoWidget struct {
	component *tview.List
	context   *app.Context
}

var heartIcon = "❤"

// isFavourite checks if the repository is a favourite
// by seeing if the database contains a match by ID
func isFavourite(ctx *app.Context) bool {
	r, err := github.GetFavouriteByRepositoryID(ctx, ctx.State.Selected.ID)
	if err != nil {
		return false
	}
	return r != nil
}

func onRepoSelect(ctx *app.Context, index int, mainText, secondaryText string, _ rune) {
	if !isFavourite(ctx) {
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

func (r *RepoWidget) Component() tview.Primitive {
	var c interface{} = r.component
	t, ok := c.(tview.Primitive)
	if !ok {
		panic("failed to cast to tview.Primitive")
	}
	return t
}

func (r *RepoWidget) IsEmpty() bool {
	return len(r.context.State.Starred) == 0
}

func (r *RepoWidget) removeFavouriteIndicator(i int, repo *domain.Repository) {
	main, secondary := r.component.GetItemText(i)
	main, _, _, _ = repositoryEntry(repo)
	r.component.SetItemText(i, main, secondary)
}

func (r *RepoWidget) SetSelected(i int) {
	r.component.SetCurrentItem(i)
}

func (r *RepoWidget) Refresh() {
	repositories, err := github.ListStarredRepositories(r.context.Client)
	if err != nil {
		openErrorModal(err)
		return
	}
	r.context.SetStarred(repositories)
	r.component.Clear()
	if len(repositories) == 0 {
		r.component.AddItem("No repositories found", "", 0, nil)
	}

	repos := repositories
	if len(repos) > 20 {
		repos = repositories[:20]
	}

	for _, repo := range repos {
		main, secondary, showSecondaryText, onSelect := repositoryEntry(repo)
		r.component.AddItem(main, secondary, 0, onSelect).
			ShowSecondaryText(showSecondaryText)
	}
	view.repos.addFavouriteIndicators()
	UI.Draw()
}

// addFavouriteIndicators loops through all repositories and if they have been previously
// favourited, adds a heart icon to the end of the name.
func (r *RepoWidget) addFavouriteIndicators() {
	for i := 0; i < r.component.GetItemCount(); i++ {
		go r.addFavouriteIndicator(i)
	}
}

func (r *RepoWidget) addFavouriteIndicator(i int) {
	if isFavourite(r.context) {
		main, secondary := r.component.GetItemText(i)
		r.component.SetItemText(i, fmt.Sprintf("%s [hotpink]%s", main, heartIcon), secondary)
	}
}

func (r *RepoWidget) OnChanged(index int, _, _ string, _ rune) {
	repo := r.context.GetStarred(index)
	if repo == nil {
		return
	}
	updateRepositoryList(r.context, repo)
}

func reposWidget(ctx *app.Context) *RepoWidget {
	widget := &RepoWidget{context: ctx}
	repos := listWidget(ListOptions{
		onChanged: widget.OnChanged,
		onSelected: func(i int, s1, s2 string, r rune) {
			onRepoSelect(ctx, i, s1, s2, r)
		},
	})
	repos.AddItem("Loading repos...", "", 0, nil)
	widget.component = repos
	return widget
}
