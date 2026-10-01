package gen

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/nullable"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	ACTION         MovieGenre = "ACTION"
	ADVENTURE      MovieGenre = "ADVENTURE"
	ANIMATION      MovieGenre = "ANIMATION"
	COMEDY         MovieGenre = "COMEDY"
	CRIME          MovieGenre = "CRIME"
	DOCUMENTARY    MovieGenre = "DOCUMENTARY"
	DRAMA          MovieGenre = "DRAMA"
	FAMILY         MovieGenre = "FAMILY"
	FANTASY        MovieGenre = "FANTASY"
	HISTORY        MovieGenre = "HISTORY"
	HORROR         MovieGenre = "HORROR"
	MUSIC          MovieGenre = "MUSIC"
	MYSTERY        MovieGenre = "MYSTERY"
	ROMANCE        MovieGenre = "ROMANCE"
	SCIENCEFICTION MovieGenre = "SCIENCE_FICTION"
	THRILLER       MovieGenre = "THRILLER"
	WAR            MovieGenre = "WAR"
	WESTERN        MovieGenre = "WESTERN"
)

func (e MovieGenre) Valid() bool {
	switch e {
	case ACTION:
		return true
	case ADVENTURE:
		return true
	case ANIMATION:
		return true
	case COMEDY:
		return true
	case CRIME:
		return true
	case DOCUMENTARY:
		return true
	case DRAMA:
		return true
	case FAMILY:
		return true
	case FANTASY:
		return true
	case HISTORY:
		return true
	case HORROR:
		return true
	case MUSIC:
		return true
	case MYSTERY:
		return true
	case ROMANCE:
		return true
	case SCIENCEFICTION:
		return true
	case THRILLER:
		return true
	case WAR:
		return true
	case WESTERN:
		return true
	default:
		return false
	}
}

const (
	DROPPED  MovieStatus = "DROPPED"
	PLANNED  MovieStatus = "PLANNED"
	WATCHED  MovieStatus = "WATCHED"
	WATCHING MovieStatus = "WATCHING"
)

func (e MovieStatus) Valid() bool {
	switch e {
	case DROPPED:
		return true
	case PLANNED:
		return true
	case WATCHED:
		return true
	case WATCHING:
		return true
	default:
		return false
	}
}

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (e Role) Valid() bool {
	switch e {
	case RoleAdmin:
		return true
	case RoleUser:
		return true
	default:
		return false
	}
}

type AddMovieToCollectionRequest struct {
	MovieId openapi_types.UUID `json:"movie_id"`
	Status  *MovieStatus       `json:"status,omitempty"`
}

type CollectionItem struct {
	CreatedAt      time.Time                 `json:"created_at"`
	Movie          Movie                     `json:"movie"`
	PersonalRating nullable.Nullable[int]    `json:"personal_rating,omitempty"`
	Review         nullable.Nullable[string] `json:"review,omitempty"`
	Status         MovieStatus               `json:"status"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

type CollectionPage struct {
	Items []CollectionItem `json:"items"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`

	Total int `json:"total"`
}

type CreateMovieRequest struct {
	Description *string    `json:"description,omitempty"`
	DurationMin int        `json:"duration_min"`
	Genre       MovieGenre `json:"genre"`
	Name        string     `json:"name"`
	ReleaseYear int        `json:"release_year"`
}

type ErrorResponse struct {
	Code    string             `json:"code"`
	Details *[]ValidationError `json:"details,omitempty"`

	Error string `json:"error"`

	Message string `json:"message"`

	Path string `json:"path"`

	RequestId *openapi_types.UUID `json:"request_id,omitempty"`

	Status int `json:"status"`

	Timestamp time.Time `json:"timestamp"`
}

type LoginRequest struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`

	ExpiresIn int `json:"expires_in"`

	TokenType string `json:"token_type"`
}

type Movie struct {
	CreatedAt   time.Time                 `json:"created_at"`
	Description nullable.Nullable[string] `json:"description,omitempty"`

	DurationMin int                `json:"duration_min"`
	Genre       MovieGenre         `json:"genre"`
	Id          openapi_types.UUID `json:"id"`

	Name string `json:"name"`

	ReleaseYear int       `json:"release_year"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MovieGenre string

type MoviePage struct {
	Items []Movie `json:"items"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`

	Total int `json:"total"`
}

type MovieStatus string

type RegisterRequest struct {
	Email      openapi_types.Email `json:"email"`
	FirstName  string              `json:"first_name"`
	Password   string              `json:"password"`
	SecondName string              `json:"second_name"`
	Username   string              `json:"username"`
}

type Role string

type UpdateCollectionItemRequest struct {
	PersonalRating nullable.Nullable[int]    `json:"personal_rating,omitempty"`
	Review         nullable.Nullable[string] `json:"review,omitempty"`
	Status         *MovieStatus              `json:"status,omitempty"`
}

type UpdateMovieRequest struct {
	Description nullable.Nullable[string] `json:"description,omitempty"`
	DurationMin *int                      `json:"duration_min,omitempty"`
	Genre       *MovieGenre               `json:"genre,omitempty"`
	Name        *string                   `json:"name,omitempty"`
	ReleaseYear *int                      `json:"release_year,omitempty"`
}

type UpdateUserRequest struct {
	IsBlocked bool `json:"is_blocked"`
}

type User struct {
	Email      openapi_types.Email `json:"email"`
	FirstName  string              `json:"first_name"`
	Id         openapi_types.UUID  `json:"id"`
	IsBlocked  bool                `json:"is_blocked"`
	Role       Role                `json:"role"`
	SecondName string              `json:"second_name"`
	Username   string              `json:"username"`
}

type UserPage struct {
	Items []User `json:"items"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`

	Total int `json:"total"`
}

type ValidationError struct {
	Field string `json:"field"`

	Message string `json:"message"`
}

type Limit = int

type Offset = int

type BadRequest = ErrorResponse

type Conflict = ErrorResponse

type Forbidden = ErrorResponse

type InternalServerError = ErrorResponse

type NotFound = ErrorResponse

type TooManyRequests = ErrorResponse

type Unauthorized = ErrorResponse

type GetCatalogParams struct {
	Genre *MovieGenre `form:"genre,omitempty" json:"genre,omitempty"`

	Search *string `form:"search,omitempty" json:"search,omitempty"`

	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
}

type GetCollectionParams struct {
	Status *MovieStatus `form:"status,omitempty" json:"status,omitempty"`

	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
}

type GetUsersParams struct {
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
}

type LoginJSONRequestBody = LoginRequest

type RegisterJSONRequestBody = RegisterRequest

type CreateMovieJSONRequestBody = CreateMovieRequest

type UpdateMovieJSONRequestBody = UpdateMovieRequest

type AddMovieToCollectionJSONRequestBody = AddMovieToCollectionRequest

type UpdateCollectionItemJSONRequestBody = UpdateCollectionItemRequest

type UpdateUserBlockStatusJSONRequestBody = UpdateUserRequest

var swaggerSpec = []string{
	"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",
}

func decodeSpec() ([]byte, error) {
	encoded := strings.Join(swaggerSpec, "")
	compressed, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr := flate.NewReader(bytes.NewReader(compressed))
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(zr); err != nil {
		return nil, fmt.Errorf("read flate: %w", err)
	}
	if err := zr.Close(); err != nil {
		return nil, fmt.Errorf("close flate reader: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

func GetSpec() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

func GetSpecJSON() ([]byte, error) {
	return rawSpec()
}

func GetSwagger() (*openapi3.T, error) {
	return GetSpec()
}
