package data

import (
	"database/sql"
	"errors"
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
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the comics data.
	query := `SELECT id, created_at, title, year, pages, version
			FROM comics
			WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var comics Comics
	err := m.DB.QueryRow(query, id).Scan(
		&comics.ID,
		&comics.CreatedAt,
		&comics.Title,
		&comics.Year,
		&comics.Pages,
		&comics.Version,
	)
	// Handle any errors. If there was no matching comics found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &comics, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m ComicsModel) Update(comics *Comics) error {
	query := `UPDATE comics
			SET title = $1, year = $2, pages = $3, version = version + 1
			WHERE id = $4
			RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		comics.Title,
		comics.Year,
		comics.Pages,
		comics.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&comics.Version)

}

// Add a placeholder method for deleting a specific record from the movies table.
func (m ComicsModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `DELETE FROM comics
			WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
