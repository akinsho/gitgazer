package github

import (
	"akinsho/gitgazer/api"
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"akinsho/gitgazer/domain"

	"golang.org/x/sync/errgroup"
)

// Fetch the users starred repositories from github
// ```
// query {
//   viewer
//     login
//	   starredRepositories(first: 20, orderBy: {field: STARRED_AT, direction: DESC}) {
// 	        nodes {
// 	          stargazerCount
// 	          description
// 	          name
// 	          issues(first: 20, orderBy: {field: UPDATED_AT, direction: DESC }) {
// 	            nodes {
// 	              state
// 	              closed
// 	              title
//  	     }
//			pullRequests(first: 5, states: OPEN, orderBy: {field: UPDATED_AT, direction:DESC}) {
// 			    edges {
// 			       node {
//					id
//					title
// 			    }
// 			}
// 	      }
// 	    }
//   }
// }
// ```
func ListStarredRepositories(client *api.Client) ([]*domain.Repository, error) {
	repos, err := client.ListStarredRepositories()
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func fetchFavouriteRepo(
	client *api.Client,
	repo *domain.FavouriteRepository,
	results chan *domain.Repository,
) error {
	r, err := client.FetchRepositoryByName(repo.Name, repo.Owner)
	if err != nil {
		return err
	}
	results <- r
	return nil
}

func RetrieveFavouriteRepositories(ctx *app.Context) ([]*domain.Repository, error) {
	saved, err := ListSavedFavourites(ctx)
	if err != nil {
		return nil, err
	}
	g := new(errgroup.Group)
	results := make(chan *domain.Repository, len(saved))
	for _, repo := range saved {
		repo := repo
		g.Go(func() error {
			return fetchFavouriteRepo(ctx.Client, repo, results)
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(results)
	repos := []*domain.Repository{}
	for result := range results {
		repos = append(repos, result)
	}

	return repos, nil
}

func GetFavouriteByRepositoryID(ctx *app.Context,
	id string,
) (favourite *domain.FavouriteRepository, err error) {
	return ctx.DB.GetFavouriteByRepoID(id)
}

func GetFavouriteRepositoryByIndex(ctx *app.Context, index int) *domain.Repository {
	favourites := ctx.State.Favourites
	if favourites != nil && len(favourites) == 0 {
		return nil
	}
	return favourites[index]
}

func ListSavedFavourites(ctx *app.Context) (repos []*domain.FavouriteRepository, err error) {
	repos, err = ctx.DB.ListFavourites()
	if err != nil {
		return
	}
	return
}

func FavouriteSelectedRepo(ctx *app.Context) (err error) {
	repo := ctx.State.Selected
	if repo == nil {
		return
	}
	_, err = ctx.DB.Insert(repo)
	if err != nil {
		return err
	}
	return nil
}

func UnfavouriteSelected(ctx *app.Context, index int) (err error) {
	repo := ctx.State.Selected
	if repo == nil {
		return
	}
	err = ctx.DB.DeleteByRepoID(repo.ID)
	ctx.SetFavourites(common.RemoveIndex(ctx.State.Favourites, index))
	if err != nil {
		return err
	}
	return nil
}
