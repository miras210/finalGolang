package data

import (
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
