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

func (f *FavouritesWidget) IsEmpty() bool {
	favs, err := github.ListSavedFavourites(f.context)
	if err != nil {
		openErrorModal(err)
		return true
	}
	if len(favs) > 0 {
		return false
	}
	return github.FavouriteRepositoryCount() == 0
}

func (f *FavouritesWidget) Component() tview.Primitive {
	var c interface{} = f.component
	t, ok := c.(tview.Primitive)
	if !ok {
		panic("failed to cast to tview.Primitive")
	}
	return t
}

func favouritesWidget(ctx *app.Context) *FavouritesWidget {
	favourites := listWidget(ListOptions{
		onSelected: func(int, string, string, rune) {},
		onChanged:  updateFavouriteChange,
	})
	return &FavouritesWidget{favourites, ctx}
}
