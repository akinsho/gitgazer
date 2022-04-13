package ui

import (
	gazerapp "akinsho/gitgazer/app"
	"akinsho/gitgazer/github"

	"github.com/rivo/tview"
)

type FavouritesWidget struct {
	component *tview.List
	context   *gazerapp.Context
}

func updateFavouriteChange(index int, _, _ string, _ rune) {
	repo := github.GetFavouriteRepositoryByIndex(index)
	if repo == nil {
		return
	}
	setRepoDescription(repo)
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
	app.Draw()
}

func (f *FavouritesWidget) Component() *tview.List {
	return f.component
}

func favouritesWidget(ctx *gazerapp.Context) *FavouritesWidget {
	favourites := tview.NewList()
	favourites.SetChangedFunc(updateFavouriteChange)
	return &FavouritesWidget{favourites, ctx}
}
