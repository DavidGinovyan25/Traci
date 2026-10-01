package domain

type (
	Role        string
	MovieGenre  string
	MovieStatus string
)

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"

	MovieGenreAction         MovieGenre = "ACTION"
	MovieGenreAdventure      MovieGenre = "ADVENTURE"
	MovieGenreAnimation      MovieGenre = "ANIMATION"
	MovieGenreComedy         MovieGenre = "COMEDY"
	MovieGenreCrime          MovieGenre = "CRIME"
	MovieGenreDocumentary    MovieGenre = "DOCUMENTARY"
	MovieGenreDrama          MovieGenre = "DRAMA"
	MovieGenreFamily         MovieGenre = "FAMILY"
	MovieGenreFantasy        MovieGenre = "FANTASY"
	MovieGenreHistory        MovieGenre = "HISTORY"
	MovieGenreHorror         MovieGenre = "HORROR"
	MovieGenreMusic          MovieGenre = "MUSIC"
	MovieGenreMystery        MovieGenre = "MYSTERY"
	MovieGenreRomance        MovieGenre = "ROMANCE"
	MovieGenreScienceFiction MovieGenre = "SCIENCE_FICTION"
	MovieGenreThriller       MovieGenre = "THRILLER"
	MovieGenreWar            MovieGenre = "WAR"
	MovieGenreWestern        MovieGenre = "WESTERN"

	MovieStatusPlanned  MovieStatus = "PLANNED"
	MovieStatusWatching MovieStatus = "WATCHING"
	MovieStatusWatched  MovieStatus = "WATCHED"
	MovieStatusDropped  MovieStatus = "DROPPED"
)
