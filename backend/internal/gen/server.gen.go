package gen

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ServerInterface interface {
	Login(w http.ResponseWriter, r *http.Request)

	Register(w http.ResponseWriter, r *http.Request)

	GetCatalog(w http.ResponseWriter, r *http.Request, params GetCatalogParams)

	CreateMovie(w http.ResponseWriter, r *http.Request)

	DeleteMovie(w http.ResponseWriter, r *http.Request, movieId openapi_types.UUID)

	GetMovieById(w http.ResponseWriter, r *http.Request, movieId openapi_types.UUID)

	UpdateMovie(w http.ResponseWriter, r *http.Request, movieId openapi_types.UUID)

	GetCollection(w http.ResponseWriter, r *http.Request, params GetCollectionParams)

	AddMovieToCollection(w http.ResponseWriter, r *http.Request)

	DeleteMovieFromCollection(w http.ResponseWriter, r *http.Request, movieId openapi_types.UUID)

	GetMovieFromCollection(w http.ResponseWriter, r *http.Request, movieId openapi_types.UUID)

	UpdateCollectionItem(w http.ResponseWriter, r *http.Request, movieId openapi_types.UUID)

	GetUsers(w http.ResponseWriter, r *http.Request, params GetUsersParams)

	DeleteUser(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID)

	GetUserById(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID)

	UpdateUserBlockStatus(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) Register(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Register(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetCatalog(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var params GetCatalogParams

	err = runtime.BindQueryParameterWithOptions("form", true, false, "genre", r.URL.Query(), &params.Genre, runtime.BindQueryParameterOptions{Type: "string", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "genre"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "genre", Err: err})
		}
		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "search", r.URL.Query(), &params.Search, runtime.BindQueryParameterOptions{Type: "string", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "search"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "search", Err: err})
		}
		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "limit", r.URL.Query(), &params.Limit, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "limit"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		}
		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "offset", r.URL.Query(), &params.Offset, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "offset"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		}
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetCatalog(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) CreateMovie(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateMovie(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var movieId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "movie_id", r.PathValue("movie_id"), &movieId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "movie_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteMovie(w, r, movieId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetMovieById(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var movieId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "movie_id", r.PathValue("movie_id"), &movieId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "movie_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetMovieById(w, r, movieId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) UpdateMovie(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var movieId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "movie_id", r.PathValue("movie_id"), &movieId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "movie_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateMovie(w, r, movieId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetCollection(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var params GetCollectionParams

	err = runtime.BindQueryParameterWithOptions("form", true, false, "status", r.URL.Query(), &params.Status, runtime.BindQueryParameterOptions{Type: "string", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "status"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "status", Err: err})
		}
		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "limit", r.URL.Query(), &params.Limit, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "limit"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		}
		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "offset", r.URL.Query(), &params.Offset, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "offset"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		}
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetCollection(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) AddMovieToCollection(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddMovieToCollection(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) DeleteMovieFromCollection(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var movieId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "movie_id", r.PathValue("movie_id"), &movieId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "movie_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteMovieFromCollection(w, r, movieId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetMovieFromCollection(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var movieId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "movie_id", r.PathValue("movie_id"), &movieId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "movie_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetMovieFromCollection(w, r, movieId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) UpdateCollectionItem(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var movieId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "movie_id", r.PathValue("movie_id"), &movieId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "movie_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateCollectionItem(w, r, movieId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetUsers(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var params GetUsersParams

	err = runtime.BindQueryParameterWithOptions("form", true, false, "limit", r.URL.Query(), &params.Limit, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "limit"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		}
		return
	}

	err = runtime.BindQueryParameterWithOptions("form", true, false, "offset", r.URL.Query(), &params.Offset, runtime.BindQueryParameterOptions{Type: "integer", Format: ""})
	if err != nil {
		var requiredError *runtime.RequiredParameterError
		if errors.As(err, &requiredError) {
			siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "offset"})
		} else {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		}
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUsers(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) DeleteUser(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var userId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "user_id", r.PathValue("user_id"), &userId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteUser(w, r, userId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) GetUserById(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var userId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "user_id", r.PathValue("user_id"), &userId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUserById(w, r, userId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) UpdateUserBlockStatus(w http.ResponseWriter, r *http.Request) {

	var err error
	_ = err

	var userId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "user_id", r.PathValue("user_id"), &userId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true, Type: "string", Format: "uuid", ValueIsUnescaped: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateUserBlockStatus(w, r, userId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	http.Handler
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc(http.MethodPost+" "+options.BaseURL+"/auth/register", wrapper.Register)
	m.HandleFunc(http.MethodPost+" "+options.BaseURL+"/auth/login", wrapper.Login)
	m.HandleFunc(http.MethodGet+" "+options.BaseURL+"/catalog", wrapper.GetCatalog)
	m.HandleFunc(http.MethodPost+" "+options.BaseURL+"/catalog", wrapper.CreateMovie)
	m.HandleFunc(http.MethodDelete+" "+options.BaseURL+"/catalog/{movie_id}", wrapper.DeleteMovie)
	m.HandleFunc(http.MethodGet+" "+options.BaseURL+"/catalog/{movie_id}", wrapper.GetMovieById)
	m.HandleFunc(http.MethodPatch+" "+options.BaseURL+"/catalog/{movie_id}", wrapper.UpdateMovie)
	m.HandleFunc(http.MethodGet+" "+options.BaseURL+"/collection", wrapper.GetCollection)
	m.HandleFunc(http.MethodPost+" "+options.BaseURL+"/collection", wrapper.AddMovieToCollection)
	m.HandleFunc(http.MethodDelete+" "+options.BaseURL+"/collection/{movie_id}", wrapper.DeleteMovieFromCollection)
	m.HandleFunc(http.MethodGet+" "+options.BaseURL+"/collection/{movie_id}", wrapper.GetMovieFromCollection)
	m.HandleFunc(http.MethodPatch+" "+options.BaseURL+"/collection/{movie_id}", wrapper.UpdateCollectionItem)
	m.HandleFunc(http.MethodGet+" "+options.BaseURL+"/users", wrapper.GetUsers)
	m.HandleFunc(http.MethodDelete+" "+options.BaseURL+"/users/{user_id}", wrapper.DeleteUser)
	m.HandleFunc(http.MethodGet+" "+options.BaseURL+"/users/{user_id}", wrapper.GetUserById)
	m.HandleFunc(http.MethodPatch+" "+options.BaseURL+"/users/{user_id}", wrapper.UpdateUserBlockStatus)

	return m
}
