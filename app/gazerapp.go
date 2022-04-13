package gazerapp

import (
	"akinsho/gitgazer/api"
	"akinsho/gitgazer/storage"
)

type Context struct {
	Client *api.Client
	DB *storage.Database
}
