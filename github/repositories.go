package github

import (
	"akinsho/gitgazer/api"
	"akinsho/gitgazer/database"
	"akinsho/gitgazer/models"
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/sync/errgroup"
)

var (
	repositories   []*models.Repository
	favourites     []*models.Repository
	issuesByRepoID = make(map[int64][]*models.Issue)
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
func ListStarredRepositories() ([]*models.Repository, error) {
	//  TODO: We need a way to invalidate previous fetched repositories
	// and refetch but this is necessary for now to prevent DDOSing the API.
	if len(repositories) > 0 {
		return repositories, nil
	}

	var starredRepositoriesQuery struct {
		Viewer struct {
			StarredRepositories struct {
				Nodes []*models.Repository `graphql:"nodes"`
			} `graphql:"starredRepositories(first: $repoCount, orderBy: {field: STARRED_AT, direction: DESC})"`
		}
	}

	err := api.Client.Query(
		context.Background(),
		&starredRepositoriesQuery,
		map[string]interface{}{
			"labelCount": githubv4.Int(20),
			"issueCount": githubv4.Int(20),
			"repoCount":  githubv4.Int(20),
			"issuesOrderBy": githubv4.IssueOrder{
				Direction: githubv4.OrderDirectionDesc,
				Field:     githubv4.IssueOrderFieldUpdatedAt,
			},
			"prCount": githubv4.Int(5),
			"prState": []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
			"pullRequestOrderBy": githubv4.IssueOrder{
				Direction: githubv4.OrderDirectionDesc,
				Field:     githubv4.IssueOrderFieldUpdatedAt,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	// FIXME: can we do better than relying on these globals
	repositories = starredRepositoriesQuery.Viewer.StarredRepositories.Nodes
	return repositories, nil
}

func fetchFavouriteRepo(repo *models.FavouriteRepository, results chan *models.Repository) error {
	var repositoryQuery struct {
		Repository models.Repository `graphql:"repository(name: $name, owner: $owner)"`
	}
	variables := map[string]interface{}{
		"name":       githubv4.String(repo.Name),
		"owner":      githubv4.String(repo.Owner),
		"labelCount": githubv4.Int(20),
		"issueCount": githubv4.Int(20),
		"issuesOrderBy": githubv4.IssueOrder{
			Direction: githubv4.OrderDirectionDesc,
			Field:     githubv4.IssueOrderFieldUpdatedAt,
		},
		"prCount": githubv4.Int(5),
		"prState": []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
		"pullRequestOrderBy": githubv4.IssueOrder{
			Direction: githubv4.OrderDirectionDesc,
			Field:     githubv4.IssueOrderFieldUpdatedAt,
		},
	}
	err := api.Client.Query(context.Background(), &repositoryQuery, variables)
	if err != nil {
		return err
	}
	results <- &repositoryQuery.Repository
	return nil
}

func RetrieveFavouriteRepositories() ([]*models.Repository, error) {
	saved, err := ListSavedFavourites()
	if err != nil {
		return nil, err
	}
	g := new(errgroup.Group)
	results := make(chan *models.Repository, len(saved))
	for _, repo := range saved {
		repo := repo
		g.Go(func() error {
			return fetchFavouriteRepo(repo, results)
		})
	}
	if g.Wait(); err != nil {
		return nil, err
	}
	close(results)
	repos := []*models.Repository{}
	for result := range results {
		repos = append(repos, result)
	}

	// FIXME: can we do better than relying on these globals
	favourites = repos
	return repos, nil
}

func GetFavouriteByRepositoryID(id string) (favourite *models.FavouriteRepository, err error) {
	return database.GetFavouriteByRepoID(id)
}

func GetRepositoryByIndex(index int) *models.Repository {
	if len(repositories) > 0 {
		return repositories[index]
	}
	return nil
}

func GetFavouriteRepositoryByIndex(index int) *models.Repository {
	if favourites != nil && len(favourites) == 0 {
		return nil
	}
	return favourites[index]
}

func ListSavedFavourites() (repos []*models.FavouriteRepository, err error) {
	repos, err = database.ListFavourites()
	if err != nil {
		return
	}
	return
}

func FavouriteRepo(index int, main, secondary string) (err error) {
	repo := GetRepositoryByIndex(index)
	if repo == nil {
		return
	}
	_, err = database.Insert(repo)
	if err != nil {
		return err
	}
	return nil
}

func UnfavouriteRepo(index int) (err error) {
	repo := GetRepositoryByIndex(index)
	if repo == nil {
		return
	}
	err = database.DeleteByRepoID(repo.ID)
	if err != nil {
		return err
	}
	return nil
}
