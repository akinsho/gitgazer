package ui

import (
	"akinsho/gitgazer/github"
	"akinsho/gitgazer/models"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RepoWidget struct {
	component *tview.List
}

var heartIcon = "❤"

// isFavourite checks if the repository is a favourite
// by seeing if the database contains a match by ID
func isFavourite(repo *models.Repository) bool {
	r, err := github.GetFavouriteByRepositoryID(repo.ID)
	if err != nil {
		return false
	}
	return r != nil
}

func onRepoSelect(index int, mainText, secondaryText string, shortcut rune) {
	repo := github.GetRepositoryByIndex(index)
	if !isFavourite(repo) {
		err := github.FavouriteRepo(index, mainText, secondaryText)
		if err != nil {
			openErrorModal(err)
			return
		}
		go view.repos.addFavouriteIndicator(index)
	} else {
		err := github.UnfavouriteRepo(index)
		if err != nil {
			openErrorModal(err)
			return
		}
		go view.repos.removeFavouriteIndicator(index, repo)
	}
	go fetchFavouriteRepositories()
}

// throttledRepoList updates the visible issue details in the issue widget when a user
// has paused over a repository in the list for more than interval time
func throttledRepoList(duration time.Duration) func(int, string, string, rune) {
	var timer *time.Timer
	return func(index int, mainText, secondaryText string, shortcut rune) {
		repo := github.GetRepositoryByIndex(index)
		if repo == nil {
			return
		}
		setRepoDescription(repo)
		if timer != nil {
			timer.Stop()
			timer = nil
		}
		timer = time.AfterFunc(duration, func() {
			view.issues.refreshIssuesList(repo)
		})
	}
}

var updateRepoList = throttledRepoList(time.Millisecond * 200)

func (r *RepoWidget) removeFavouriteIndicator(i int, repo *models.Repository) {
	main, secondary := r.component.GetItemText(i)
	main, _, _, _ = repositoryEntry(repo)
	r.component.SetItemText(i, main, secondary)
}

func (r *RepoWidget) refreshRepositoryList(repositories []*models.Repository, err error) {
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
	app.Draw()
	app.SetFocus(r.component)
}

// addFavouriteIndicators loops through all repositories and if they have been previously
// favourited, adds a heart icon to the end of the name.
func (r *RepoWidget) addFavouriteIndicators() {
	for i := 0; i < r.component.GetItemCount(); i++ {
		go r.addFavouriteIndicator(i)
	}
}

func (r *RepoWidget) addFavouriteIndicator(i int) {
	if isFavourite(github.GetRepositoryByIndex(i)) {
		main, secondary := r.component.GetItemText(i)
		r.component.SetItemText(i, fmt.Sprintf("%s [hotpink]%s", main, heartIcon), secondary)
	}
}

func reposWidget() *RepoWidget {
	repos := tview.NewList()
	repos.AddItem("Loading repos...", "", 0, nil)

	repos.SetChangedFunc(updateRepoList).
		SetSelectedFunc(onRepoSelect).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorForestGreen).
		SetMainTextColor(tcell.ColorForestGreen).
		SetMainTextStyle(tcell.StyleDefault.Bold(true)).
		SetSecondaryTextColor(tcell.ColorDarkGray).SetBorderPadding(0, 0, 1, 1)
	return &RepoWidget{component: repos}
}
