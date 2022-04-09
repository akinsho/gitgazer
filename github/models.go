package github

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
	} `graphql:"issues(first: $issueCount)"`
}
