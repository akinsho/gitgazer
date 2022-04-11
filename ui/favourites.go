package ui

import (
	"akinsho/gogazer/github"
	"akinsho/gogazer/models"

	"github.com/rivo/tview"
)

type FavouritesWidget struct {
	component *tview.List
}

func updateFavouriteChange(index int, mainText, secondaryText string, shortcut rune) {
	repo := github.GetFavouriteRepositoryByIndex(index)
	if repo == nil {
		return
	}
	// setRepoDescription(repo)
}

// refreshFavouritesList fetches all saved repositories from the database and
// adds them to the View.favourites list.
func (f *FavouritesWidget) refreshFavouritesList(
	favourites []models.FavouriteRepository,
	err error,
) {
	if f.component.GetItemCount() > 0 {
		f.component.Clear()
	}
	if err != nil {
		openErrorModal(err)
		return
	}
	if len(favourites) == 0 {
		f.component.AddItem("No favourites found", "", 0, nil)
	}

	repos := favourites
	if len(repos) > 20 {
		repos = favourites[:20]
	}

	for _, repo := range favourites {
		main, secondary, showSecondaryText, onSelect := repositoryEntry(&repo)
		f.component.AddItem(main, secondary, 0, onSelect).
			ShowSecondaryText(showSecondaryText)
	}
	app.Draw()
}

func favouritesWidget() *FavouritesWidget {
	favourites := tview.NewList()
	favourites.SetChangedFunc(updateFavouriteChange)
	return &FavouritesWidget{favourites}
}
