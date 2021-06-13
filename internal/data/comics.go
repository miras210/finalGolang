package data

import (
	"database/sql"
	"github.com/miras210/finalGolang/internal/validator"
	"time"
)

type Comics struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Pages     Pages     `json:"pages,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateComics(v *validator.Validator, comics *Comics) {
	v.Check(comics.Title != "", "title", "must be provided")
	v.Check(len(comics.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(comics.Year != 0, "year", "must be provided")
	v.Check(comics.Year >= 1888, "year", "must be greater than 1888")
	v.Check(comics.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(comics.Pages != 0, "pages", "must be provided")
	v.Check(comics.Pages > 0, "pages", "must be a positive integer")

}

type ComicsModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m ComicsModel) Insert(comics *Comics) error {
	query := `INSERT INTO comics (title, year, pages)
			VALUES ($1, $2, $3)
			RETURNING id, created_at, version`

	args := []interface{}{comics.Title, comics.Year, comics.Pages}
	return m.DB.QueryRow(query, args...).Scan(&comics.ID, &comics.CreatedAt, &comics.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m ComicsModel) Get(id int64) (*Comics, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m ComicsModel) Update(comics *Comics) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m ComicsModel) Delete(id int64) error {
	return nil
}
