package models

type FavouriteRepository struct {
	ID          int64
	RepoID      string
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

type Repo interface {
	GetDescription() string
	GetName() string
}
