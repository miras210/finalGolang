package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&comics.ID, &comics.CreatedAt, &comics.Version)
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

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
			WHERE id = $4 AND version = $5
			RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		comics.Title,
		comics.Year,
		comics.Pages,
		comics.ID,
		comics.Version,
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query. If no matching row could be found, we know the movie
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&comics.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
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

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m ComicsModel) GetAll(title string, year int, filters Filters) ([]*Comics, Metadata, error) {
	// Define the SQL query for retrieving the comics data.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, title, year, pages, version
		FROM comics
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (year = $2 OR $2 = -1)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, year, filters.limit(), filters.offset()}

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()

	totalRecords := 0
	// Initialize an empty slice to hold the movie data.
	comics := []*Comics{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var comic Comics
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&comic.ID,
			&comic.CreatedAt,
			&comic.Title,
			&comic.Year,
			&comic.Pages,
			&comic.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Movie struct to the slice.
		comics = append(comics, &comic)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of movies.
	return comics, metadata, nil
}
