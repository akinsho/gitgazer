package github

import (
	"context"

	"github.com/shurcooL/githubv4"
)

var (
	repositories   []*Repository
	issuesByRepoID = make(map[int64][]*Issue)
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
// 	          issues(first: 20) {
// 	            nodes {
// 	              state
// 	              closed
// 	              title
// 	            }
// 	          }
// 	        }
// 	      }
// 	    }
//   }
// }
// ```
func FetchRepositories(client *githubv4.Client) ([]*Repository, error) {
	//  TODO: We need a way to invalidate previous fetched repositories
	// and refetch but this is necessary for now to prevent DDOSing the API.
	if len(repositories) > 0 {
		return repositories, nil
	}

	var starredRepositoriesQuery struct {
		Viewer struct {
			StarredRepositories struct {
				Nodes []*Repository `graphql:"nodes"`
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
		},
	)
	if err != nil {
		return nil, err
	}
	repositories = starredRepositoriesQuery.Viewer.StarredRepositories.Nodes
	return repositories, nil
}

func GetRepositoryByIndex(index int) *Repository {
	if len(repositories) > 0 {
		return repositories[index]
	}
	return nil
}
