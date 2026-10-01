package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type MovieGenre string

const (
	MovieGenreACTION         MovieGenre = "ACTION"
	MovieGenreADVENTURE      MovieGenre = "ADVENTURE"
	MovieGenreANIMATION      MovieGenre = "ANIMATION"
	MovieGenreCOMEDY         MovieGenre = "COMEDY"
	MovieGenreCRIME          MovieGenre = "CRIME"
	MovieGenreDOCUMENTARY    MovieGenre = "DOCUMENTARY"
	MovieGenreDRAMA          MovieGenre = "DRAMA"
	MovieGenreFAMILY         MovieGenre = "FAMILY"
	MovieGenreFANTASY        MovieGenre = "FANTASY"
	MovieGenreHISTORY        MovieGenre = "HISTORY"
	MovieGenreHORROR         MovieGenre = "HORROR"
	MovieGenreMUSIC          MovieGenre = "MUSIC"
	MovieGenreMYSTERY        MovieGenre = "MYSTERY"
	MovieGenreROMANCE        MovieGenre = "ROMANCE"
	MovieGenreSCIENCEFICTION MovieGenre = "SCIENCE_FICTION"
	MovieGenreTHRILLER       MovieGenre = "THRILLER"
	MovieGenreWAR            MovieGenre = "WAR"
	MovieGenreWESTERN        MovieGenre = "WESTERN"
)

func (e *MovieGenre) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MovieGenre(s)
	case string:
		*e = MovieGenre(s)
	default:
		return fmt.Errorf("unsupported scan type for MovieGenre: %T", src)
	}
	return nil
}

type NullMovieGenre struct {
	MovieGenre MovieGenre `json:"movie_genre"`
	Valid      bool       `json:"valid"`
}

func (ns *NullMovieGenre) Scan(value interface{}) error {
	if value == nil {
		ns.MovieGenre, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MovieGenre.Scan(value)
}

func (ns NullMovieGenre) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MovieGenre), nil
}

type MovieStatus string

const (
	MovieStatusPLANNED  MovieStatus = "PLANNED"
	MovieStatusWATCHING MovieStatus = "WATCHING"
	MovieStatusWATCHED  MovieStatus = "WATCHED"
	MovieStatusDROPPED  MovieStatus = "DROPPED"
)

func (e *MovieStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MovieStatus(s)
	case string:
		*e = MovieStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for MovieStatus: %T", src)
	}
	return nil
}

type NullMovieStatus struct {
	MovieStatus MovieStatus `json:"movie_status"`
	Valid       bool        `json:"valid"`
}

func (ns *NullMovieStatus) Scan(value interface{}) error {
	if value == nil {
		ns.MovieStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MovieStatus.Scan(value)
}

func (ns NullMovieStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MovieStatus), nil
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role `json:"role"`
	Valid bool `json:"valid"`
}

func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type Movie struct {
	ID          pgtype.UUID        `json:"id"`
	Name        string             `json:"name"`
	DurationMin int32              `json:"duration_min"`
	Genre       MovieGenre         `json:"genre"`
	ReleaseYear int32              `json:"release_year"`
	Description pgtype.Text        `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type User struct {
	ID           pgtype.UUID `json:"id"`
	Role         Role        `json:"role"`
	Username     string      `json:"username"`
	FirstName    string      `json:"first_name"`
	SecondName   string      `json:"second_name"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"password_hash"`
	IsBlocked    bool        `json:"is_blocked"`
}

type UserMovie struct {
	ID             pgtype.UUID        `json:"id"`
	MovieID        pgtype.UUID        `json:"movie_id"`
	UserID         pgtype.UUID        `json:"user_id"`
	PersonalRating pgtype.Int4        `json:"personal_rating"`
	Review         pgtype.Text        `json:"review"`
	Status         MovieStatus        `json:"status"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}
