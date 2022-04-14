package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/github"

	"github.com/rivo/tview"
)

type FavouritesWidget struct {
	component *tview.List
	context   *app.Context
}

func updateFavouriteChange(index int, _, _ string, _ rune) {
	repo := github.GetFavouriteRepositoryByIndex(index)
	if repo == nil {
		return
	}
	updateRepositoryList(repo)
}

// refreshFavouritesList fetches all saved repositories from the database and
// adds them to the View.favourites list.
func (f *FavouritesWidget) Refresh() {
	favourites, err := github.RetrieveFavouriteRepositories(f.context)
	if err != nil {
		openErrorModal(err)
		return
	}
	if f.component.GetItemCount() > 0 {
		f.component.Clear()
	}
	if len(favourites) == 0 {
		f.component.AddItem("No favourites found", "", 0, nil)
	}

	repos := favourites
	if len(repos) > 20 {
		repos = favourites[:20]
	}

	for _, repo := range repos {
		main, secondary, showSecondaryText, onSelect := repositoryEntry(repo)
		f.component.AddItem(main, secondary, 0, onSelect).
			ShowSecondaryText(showSecondaryText)
	}
	UI.Draw()
}

func (f *FavouritesWidget) Component() *tview.List {
	return f.component
}

func favouritesWidget(ctx *app.Context) *FavouritesWidget {
	favourites := tview.NewList()
	favourites.SetBorderPadding(0, 0, 1, 1)
	favourites.SetChangedFunc(updateFavouriteChange)
	return &FavouritesWidget{favourites, ctx}
}
