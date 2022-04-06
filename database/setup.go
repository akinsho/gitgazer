package database

import (
	"database/sql"
	"errors"

	"github.com/google/go-github/v43/github"
	_ "github.com/mattn/go-sqlite3"
)

const file string = "gazers.db"

type Gazers struct {
	db *sql.DB
}

const create string = `
  CREATE TABLE IF NOT EXISTS gazed_repositories (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT
  );`

func Setup() (*Gazers, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &Gazers{
		db: db,
	}, nil
}

// Insert a new repository into the database.
func (g *Gazers) Insert(repo *github.Repository) (int64, error) {
	if repo == nil {
		return 0, errors.New("could not save repository as it is missing!")
	}
	res, err := g.db.Exec(
		"INSERT INTO gazed_repositories (id, name, description) VALUES (?, ?, ?);",
		repo.GetID(),
		repo.GetName(),
		repo.GetDescription(),
	)
	if err != nil {
		return 0, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, nil
	}
	return id, nil
}
