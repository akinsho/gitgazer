package app

import (
	"akinsho/gitgazer/api"
	"akinsho/gitgazer/domain"
	"akinsho/gitgazer/storage"
	"errors"
	"fmt"
)

type State struct {
	Favourites []*domain.Repository
	Starred    []*domain.Repository
	Selected   *domain.Repository
}

type Context struct {
	Client *api.Client
	DB     *storage.Database
	Config *Config
	State  *State
}

func (c *Context) GetStarred(index int) *domain.Repository {
	if len(c.State.Starred) == 0 {
		return nil
	}
	return c.State.Starred[index]
}

func (c *Context) GetFavourite(index int) (*domain.Repository, error) {
	favs := c.State.Favourites
	if len(favs) == 0 {
		return nil, nil
	}
	if index < 0 || index > len(favs)-1 {
		return nil, errors.New(
			fmt.Sprintf(
				"[GetFavourite] Index is out of range: %d, length was %d",
				index,
				len(favs),
			),
		)
	}
	return favs[index], nil
}

func (c *Context) SetFavourites(favourites []*domain.Repository) {
	c.State.Favourites = favourites
}

func (c *Context) SetStarred(starred []*domain.Repository) {
	c.State.Starred = starred
}

func (c *Context) SetSelected(selected *domain.Repository) {
	c.State.Selected = selected
}
