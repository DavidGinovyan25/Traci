package httptransport

import (
	"math"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"traci/backend/internal/domain"
	"traci/backend/internal/gen"
)

func (h *Handler) GetCatalog(w http.ResponseWriter, r *http.Request, params gen.GetCatalogParams) {
	if h.actor(w, r, false) == nil {
		return
	}
	limit, offset := pagination(params.Limit, params.Offset)
	filter := domain.MovieFilter{Limit: limit, Offset: offset}
	if params.Genre != nil {
		genre := domain.MovieGenre(*params.Genre)
		filter.Genre = &genre
	}
	if params.Search != nil {
		filter.Search = *params.Search
	}
	page, err := h.service.ListMovies(r.Context(), filter)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	items := make([]gen.Movie, len(page.Items))
	for i, movie := range page.Items {
		items[i] = movieResponse(movie)
	}
	writeJSON(w, http.StatusOK, gen.MoviePage{Items: items, Limit: page.Limit, Offset: page.Offset, Total: page.Total})
}

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	actor := h.actor(w, r, true)
	if actor == nil {
		return
	}
	body, ok := decodeBody[gen.CreateMovieRequest](w, r)
	if !ok {
		return
	}
	if strings.TrimSpace(body.Name) == "" || body.DurationMin > math.MaxInt32 || body.ReleaseYear > math.MaxInt32 {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid movie fields")
		return
	}
	movie := domain.NewMovie(body.Name, body.DurationMin, domain.MovieGenre(body.Genre), body.ReleaseYear, body.Description)
	if err := h.service.CreateMovie(r.Context(), actor, movie); err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Location", "/catalog/"+movie.ID.String())
	writeJSON(w, http.StatusCreated, movieResponse(*movie))
}

func (h *Handler) GetMovieById(w http.ResponseWriter, r *http.Request, movieID uuid.UUID) {
	if h.actor(w, r, false) == nil {
		return
	}
	movie, err := h.service.GetMovie(r.Context(), movieID)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, movieResponse(*movie))
}

func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request, movieID uuid.UUID) {
	actor := h.actor(w, r, true)
	if actor == nil {
		return
	}
	body, ok := decodeBody[gen.UpdateMovieRequest](w, r)
	if !ok {
		return
	}
	if (body.Name != nil && strings.TrimSpace(*body.Name) == "") || (body.DurationMin != nil && *body.DurationMin > math.MaxInt32) || (body.ReleaseYear != nil && *body.ReleaseYear > math.MaxInt32) {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid movie fields")
		return
	}
	movie, err := h.service.GetMovie(r.Context(), movieID)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if body.Name != nil {
		movie.Name = *body.Name
	}
	if body.DurationMin != nil {
		movie.DurationMin = *body.DurationMin
	}
	if body.Genre != nil {
		movie.Genre = domain.MovieGenre(*body.Genre)
	}
	if body.ReleaseYear != nil {
		movie.ReleaseYear = *body.ReleaseYear
	}
	if body.Description.IsSpecified() {
		movie.Description = nullablePointer(body.Description)
	}
	movie, err = h.service.UpdateMovie(r.Context(), actor, movie)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, movieResponse(*movie))
}

func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request, movieID uuid.UUID) {
	actor := h.actor(w, r, true)
	if actor == nil {
		return
	}
	if err := h.service.DeleteMovie(r.Context(), actor, movieID); err != nil {
		h.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
