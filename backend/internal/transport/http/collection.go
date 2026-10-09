package httptransport

import (
	"net/http"

	"github.com/google/uuid"

	"traci/backend/internal/domain"
	"traci/backend/internal/gen"
)

func (h *Handler) GetCollection(w http.ResponseWriter, r *http.Request, params gen.GetCollectionParams) {
	actor := h.actor(w, r, false)
	if actor == nil {
		return
	}
	limit, offset := pagination(params.Limit, params.Offset)
	filter := domain.CollectionFilter{Limit: limit, Offset: offset}
	if params.Status != nil {
		status := domain.MovieStatus(*params.Status)
		filter.Status = &status
	}
	page, err := h.service.ListCollection(r.Context(), actor.ID, filter)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	items := make([]gen.CollectionItem, len(page.Items))
	for i, entry := range page.Items {
		items[i] = collectionResponse(entry)
	}
	writeJSON(w, http.StatusOK, gen.CollectionPage{Items: items, Limit: page.Limit, Offset: page.Offset, Total: page.Total})
}

func (h *Handler) AddMovieToCollection(w http.ResponseWriter, r *http.Request) {
	actor := h.actor(w, r, false)
	if actor == nil {
		return
	}
	body, ok := decodeBody[gen.AddMovieToCollectionRequest](w, r)
	if !ok {
		return
	}
	status := domain.MovieStatusPlanned
	if body.Status != nil {
		status = domain.MovieStatus(*body.Status)
	}
	entry, err := h.service.AddCollectionItem(r.Context(), actor.ID, body.MovieId, status)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Location", "/collection/"+body.MovieId.String())
	writeJSON(w, http.StatusCreated, collectionResponse(entry))
}

func (h *Handler) GetMovieFromCollection(w http.ResponseWriter, r *http.Request, movieID uuid.UUID) {
	actor := h.actor(w, r, false)
	if actor == nil {
		return
	}
	entry, err := h.service.GetCollectionItem(r.Context(), actor.ID, movieID)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, collectionResponse(entry))
}

func (h *Handler) UpdateCollectionItem(w http.ResponseWriter, r *http.Request, movieID uuid.UUID) {
	actor := h.actor(w, r, false)
	if actor == nil {
		return
	}
	body, ok := decodeBody[gen.UpdateCollectionItemRequest](w, r)
	if !ok {
		return
	}
	entry, err := h.service.GetCollectionItem(r.Context(), actor.ID, movieID)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if body.Status != nil {
		entry.Item.Status = domain.MovieStatus(*body.Status)
	}
	if body.PersonalRating.IsSpecified() {
		entry.Item.PersonalRating = nullablePointer(body.PersonalRating)
	}
	if body.Review.IsSpecified() {
		entry.Item.Review = nullablePointer(body.Review)
	}
	entry, err = h.service.UpdateCollectionItem(r.Context(), &entry.Item)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, collectionResponse(entry))
}

func (h *Handler) DeleteMovieFromCollection(w http.ResponseWriter, r *http.Request, movieID uuid.UUID) {
	actor := h.actor(w, r, false)
	if actor == nil {
		return
	}
	if err := h.service.DeleteCollectionItem(r.Context(), actor.ID, movieID); err != nil {
		h.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
