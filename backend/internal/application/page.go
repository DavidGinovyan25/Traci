package application

import "traci/backend/internal/domain"

type Page[T any] struct {
	Items  []T
	Limit  int
	Offset int
	Total  int
}

type CollectionEntry struct {
	Item  domain.CollectionItem
	Movie domain.Movie
}

type MoviePage = Page[domain.Movie]
type UserPage = Page[domain.User]
type CollectionPage = Page[CollectionEntry]
