package github

import (
	"akinsho/gogazer/database"
	"akinsho/gogazer/models"
	"context"

	"github.com/shurcooL/githubv4"
)

var (
	repositories   []*models.Repository
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
// 	      }
// 	    }
//   }
// }
// ```
func FetchRepositories(client *githubv4.Client) ([]*models.Repository, error) {
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

	err := client.Query(
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
		},
	)
	if err != nil {
		return nil, err
	}
	repositories = starredRepositoriesQuery.Viewer.StarredRepositories.Nodes
	return repositories, nil
}

func GetRepositoryByIndex(index int) *models.Repository {
	if len(repositories) > 0 {
		return repositories[index]
	}
	return nil
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
