package ui

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/github"
	"fmt"

	"github.com/rivo/tview"
)

type FavouritesWidget struct {
	component *tview.List
	context   *app.Context
}

func (f *FavouritesWidget) Context() *app.Context {
	return f.context
}

func (f *FavouritesWidget) OnChanged(index int, main, _ string, _ rune) {
	repo, err := f.context.GetFavourite(index)
	if err != nil {
		openErrorModal(err)
		return
	}
	if repo == nil {
		return
	}
	updateRepositoryList(f.context, repo)
}

func (f *FavouritesWidget) SetSelected(i int) {
	f.component.SetCurrentItem(i)
}

// refreshFavouritesList fetches all saved repositories from the database and
// adds them to the View.favourites list.
func (f *FavouritesWidget) Refresh() (err error) {
	// FIXME: this happens twice on startup rather than once which causes weird intermittent
	// race conditions.
	f.context.Logger.Write("refreshing favourites list")
	f.component.Clear()
	favourites := f.context.State.Favourites
	if len(favourites) == 0 {
		f.component.AddItem("Loading favourites...", "", 0, nil)
		UI.Draw()
		favourites, err = github.RetrieveFavouriteRepositories(f.context)
		if err != nil {
			return err
		}
		f.context.SetFavourites(favourites)
	}
	f.component.Clear()
	if len(favourites) == 0 {
		f.component.AddItem("No favourites found", "", 0, nil)
		return
	}

	favs := favourites
	if len(favs) > 20 {
		favs = favourites[:20]
	}

	for _, repo := range favs {
		main, secondary, showSecondaryText, onSelect := repositoryEntry(repo)
		f.component.AddItem(main, secondary, 0, onSelect).ShowSecondaryText(showSecondaryText)
	}
	f.context.Logger.Write(fmt.Sprintf("Favourites item count: %d", f.component.GetItemCount()))
	return
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
	return len(favs) == 0
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
	widget := &FavouritesWidget{context: ctx}
	favourites := listWidget(ListOptions{
		onSelected: func(int, string, string, rune) {},
		onChanged:  widget.OnChanged,
	})
	widget.component = favourites
	return widget
}
