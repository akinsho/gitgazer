package database

import (
	"akinsho/gogazer/models"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

const file string = "gazers.db"

type Gazers struct {
	db *sql.DB
}

const create string = `
  CREATE TABLE IF NOT EXISTS gazed_repositories (
	id INTEGER NOT NULL PRIMARY KEY,
	repo_id STRING NOT NULL UNIQUE,
	name TEXT NOT NULL,
	description TEXT
  );`

var connection *Gazers

func Setup() error {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}
	if _, err := db.Exec(create); err != nil {
		return err
	}
	connection = &Gazers{
		db: db,
	}
	return nil
}

// Insert a new repository into the database.
func Insert(repo *models.Repository) (int64, error) {
	if repo == nil {
		return 0, errors.New("could not save repository as it is missing!")
	}
	res, err := connection.db.Exec(
		"INSERT INTO gazed_repositories (repo_id, name, description) VALUES (?, ?, ?);",
		repo.ID,
		repo.Name,
		repo.Description,
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

// ListFavourites pulls the repositories out of the gazers table and returns them as a list
func ListFavourites() ([]models.FavouriteRepository, error) {
	rows, err := connection.db.Query("SELECT * FROM gazed_repositories;")
	if err != nil {
		return nil, err
	}
	repositories := []models.FavouriteRepository{}
	for rows.Next() {
		var id int64
		var repoID string
		var name string
		var description string
		if err := rows.Scan(&id, &repoID, &name, &description); err != nil {
			return nil, err
		}
		repositories = append(repositories, models.FavouriteRepository{
			ID:          id,
			RepoID:      repoID,
			Name:        name,
			Description: description,
		})
	}
	return repositories, nil
}
