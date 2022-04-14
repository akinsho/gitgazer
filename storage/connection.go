package storage

import (
	"akinsho/gitgazer/domain"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	sqlDB *sql.DB
}

const create string = `
  CREATE TABLE IF NOT EXISTS gazed_repositories (
	id INTEGER NOT NULL PRIMARY KEY,
	repo_id STRING NOT NULL UNIQUE,
	name TEXT NOT NULL,
	owner TEXT NOT NULL,
	description TEXT
  );`

func Setup(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

// Insert a new repository into the database.
func (db *Database) Insert(repo *domain.Repository) (int64, error) {
	if repo == nil {
		return 0, errors.New("could not save repository as it is missing!")
	}
	res, err := db.sqlDB.Exec(
		"INSERT OR IGNORE INTO gazed_repositories (repo_id, name, owner, description) VALUES (?, ?, ?, ?);",
		repo.ID,
		repo.Name,
		repo.Owner.Login,
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

// Delete removes a repository with the matching repo ID from the database.
func (db *Database) DeleteByRepoID(id string) error {
	_, err := db.sqlDB.Exec("DELETE FROM gazed_repositories WHERE repo_id = ?;", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetFavouriteByRepoID(id string) (*domain.FavouriteRepository, error) {
	row := db.sqlDB.QueryRow("SELECT * FROM gazed_repositories WHERE repo_id = ?;", id)
	repo := &domain.FavouriteRepository{}
	if err := row.Scan(
		&repo.ID,
		&repo.RepoID,
		&repo.Name,
		&repo.Owner,
		&repo.Description,
	); err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return repo, nil
}

// ListFavourites pulls the repositories out of the gazers table and returns them as a list
func (db *Database) ListFavourites() ([]*domain.FavouriteRepository, error) {
	rows, err := db.sqlDB.Query("SELECT * FROM gazed_repositories;")
	if err != nil {
		return nil, err
	}
	repositories := []*domain.FavouriteRepository{}
	for rows.Next() {
		repo := &domain.FavouriteRepository{}
		if err := rows.Scan(&repo.ID, &repo.RepoID, &repo.Name, &repo.Owner, &repo.Description); err != nil {
			return nil, err
		}
		repositories = append(repositories, repo)
	}
	return repositories, nil
}
