package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/domain"
	"akinsho/gitgazer/github"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RepoWidget struct {
	component *tview.List
	context   *app.Context
}

var heartIcon = "❤"

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
	repo := github.GetRepositoryByIndex(index)
	if !isFavourite(ctx, repo) {
		err := github.FavouriteRepo(ctx, index, mainText, secondaryText)
		if err != nil {
			openErrorModal(err)
			return
		}
		go view.repos.addFavouriteIndicator(index)
	} else {
		err := github.UnfavouriteRepo(ctx, index)
		if err != nil {
			openErrorModal(err)
			return
		}
		go view.repos.removeFavouriteIndicator(index, repo)
	}
}

func (r *RepoWidget) Component() *tview.List {
	return r.component
}

func (r *RepoWidget) removeFavouriteIndicator(i int, repo *domain.Repository) {
	main, secondary := r.component.GetItemText(i)
	main, _, _, _ = repositoryEntry(repo)
	r.component.SetItemText(i, main, secondary)
}

func (r *RepoWidget) Refresh() {
	repositories, err := github.ListStarredRepositories(r.context.Client)
	if err != nil {
		openErrorModal(err)
		return
	}
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
	UI.SetFocus(r.component)
}

// addFavouriteIndicators loops through all repositories and if they have been previously
// favourited, adds a heart icon to the end of the name.
func (r *RepoWidget) addFavouriteIndicators() {
	for i := 0; i < r.component.GetItemCount(); i++ {
		go r.addFavouriteIndicator(i)
	}
}

func (r *RepoWidget) addFavouriteIndicator(i int) {
	if isFavourite(r.context, github.GetRepositoryByIndex(i)) {
		main, secondary := r.component.GetItemText(i)
		r.component.SetItemText(i, fmt.Sprintf("%s [hotpink]%s", main, heartIcon), secondary)
	}
}

func updateStarredList(index int, _, _ string, _ rune) {
	repo := github.GetRepositoryByIndex(index)
	if repo == nil {
		return
	}
	updateRepoList(repo)
}

func reposWidget(ctx *app.Context) *RepoWidget {
	repos := tview.NewList()
	repos.AddItem("Loading repos...", "", 0, nil)

	repos.SetChangedFunc(updateStarredList).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			onRepoSelect(ctx, i, s1, s2, r)
		}).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorForestGreen).
		SetMainTextColor(tcell.ColorForestGreen).
		SetMainTextStyle(tcell.StyleDefault.Bold(true)).
		SetSecondaryTextColor(tcell.ColorDarkGray).SetBorderPadding(0, 0, 1, 1)
	return &RepoWidget{component: repos, context: ctx}
}
