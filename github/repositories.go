package github

import (
	"akinsho/gitgazer/api"
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/common"
	"akinsho/gitgazer/domain"

	"golang.org/x/sync/errgroup"
)

var (
	repositories   []*domain.Repository
	favourites     []*domain.Repository
	issuesByRepoID = make(map[int64][]*domain.Issue)
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
	//  TODO: We need a way to invalidate previous fetched repositories
	// and refetch but this is necessary for now to prevent DDOSing the API.
	if len(repositories) > 0 {
		return repositories, nil
	}
	repos, err := client.ListStarredRepositories()
	if err != nil {
		return nil, err
	}
	// FIXME: can we do better than relying on these globals
	repositories = repos
	return repositories, nil
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
	if g.Wait(); err != nil {
		return nil, err
	}
	close(results)
	repos := []*domain.Repository{}
	for result := range results {
		repos = append(repos, result)
	}

	// FIXME: can we do better than relying on these globals
	favourites = repos
	return repos, nil
}

func GetFavouriteByRepositoryID(
	ctx *app.Context,
	id string,
) (favourite *domain.FavouriteRepository, err error) {
	return ctx.DB.GetFavouriteByRepoID(id)
}

func GetRepositoryByIndex(index int) *domain.Repository {
	if len(repositories) > 0 {
		return repositories[index]
	}
	return nil
}

func GetFavouriteRepositoryByIndex(index int) *domain.Repository {
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

func FavouriteRepositoryCount() int {
	return len(favourites)
}

func StarredRepositoryCount() int {
	return len(repositories)
}

func FavouriteRepo(ctx *app.Context, index int, main, secondary string) (err error) {
	repo := GetRepositoryByIndex(index)
	if repo == nil {
		return
	}
	_, err = ctx.DB.Insert(repo)
	if err != nil {
		return err
	}
	return nil
}

func UnfavouriteRepo(ctx *app.Context, index int) (err error) {
	repo := GetRepositoryByIndex(index)
	if repo == nil {
		return
	}
	err = ctx.DB.DeleteByRepoID(repo.ID)
	favourites = common.RemoveIndex(favourites, index)
	if err != nil {
		return err
	}
	return nil
}
