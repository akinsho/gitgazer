package models

type Author struct {
	Login string
}

type Issue struct {
	State  string
	Closed bool
	Title  string
	Number int
	Author *Author
	Labels struct {
		Nodes []*Label
	} `graphql:"labels(first: $labelCount)"`
}

type Label struct {
	Name  string
	Color string
}

type Repository struct {
	ID             string
	StargazerCount int
	Description    string
	Name           string
	Issues         struct {
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
