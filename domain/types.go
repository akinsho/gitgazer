package domain

type Repo interface {
	GetDescription() string
	GetName() string
}

type FavouriteRepository struct {
	ID          int64
	RepoID      string
	Owner       string
	Description string
	Name        string
}

func (r *FavouriteRepository) GetDescription() string {
	if r == nil {
		return ""
	}
	return r.Description
}

func (r *FavouriteRepository) GetName() string {
	if r == nil {
		return ""
	}
	return r.Name
}

type Author struct {
	Login string
}

type Issue struct {
	State  string
	Closed bool
	Title  string
	Number int
	Author *Author
	Body   string
	Labels struct {
		Nodes []*Label
	} `graphql:"labels(first: $labelCount)"`
}

type Label struct {
	Name  string
	Color string
}

type PullRequest struct {
	Title string
	ID    string
}

type RepositoryOwner struct {
	ID    string
	Login string
}

type Repository struct {
	ID             string
	Owner          *RepositoryOwner
	StargazerCount int
	Description    string
	Name           string
	URL            string
	PullRequests   struct {
		TotalCount int
		Nodes      []*PullRequest
	} `graphql:"pullRequests(first: $prCount, states: $prState, orderBy: $pullRequestOrderBy)"`
	Issues struct {
		Nodes []*Issue
	} `graphql:"issues(first: $issueCount, orderBy: $issuesOrderBy)"`
}

//--------------------------------------------------------------------------------------------------
//  Repository Getters
//--------------------------------------------------------------------------------------------------

func (r *Repository) GetID() string {
	if r == nil {
		return ""
	}
	return r.ID
}

func (r *Repository) GetDescription() string {
	if r == nil {
		return ""
	}
	return r.Description
}

func (r *Repository) GetName() string {
	if r == nil {
		return ""
	}
	return r.Name
}

func (r *Repository) GetStargazerCount() int {
	if r == nil {
		return 0
	}
	return r.StargazerCount
}

func (r *Repository) GetPullRequestCount() int {
	if r == nil {
		return 0
	}
	return r.PullRequests.TotalCount
}

func (r *Repository) GetIssueCount() int {
	if r == nil {
		return 0
	}
	return len(r.Issues.Nodes)
}

func (r *Repository) GetIssues() []*Issue {
	if r == nil {
		return []*Issue{}
	}
	return r.Issues.Nodes
}

// Getters for the Issue struct
func (i *Issue) GetState() string {
	if i == nil {
		return ""
	}
	return i.State
}

func (i *Issue) GetClosed() bool {
	if i == nil {
		return false
	}
	return i.Closed
}

func (i *Issue) GetTitle() string {
	if i == nil {
		return ""
	}
	return i.Title
}

func (i *Issue) GetNumber() int {
	if i == nil {
		return 0
	}
	return i.Number
}
