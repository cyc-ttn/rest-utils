package controller;

// Package to act as a 'controller' which helps with mapping
// methods to router-set.
//
// Uses gorilla mux

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/urfave/negroni"
)

// Main Router
var router = initRouter();

// Get - retrieves the router
func Get() *mux.Router{
  return router;
}

// initRouter initializes the router
func initRouter() *mux.Router {
  return mux.NewRouter();
}

// PrefixRequest returns a subrouter prefixed to the value
func PrefixRequest(value string) *mux.Router {
  return router.PathPrefix(value).Subrouter().StrictSlash(true)
}

// PrefixRequestWithHandler prefixes the request with Middleware
func PrefixRequestWithHandler(value string, h http.Handler){
  router.PathPrefix(value).Handler( h )
}

// PrefixRequestWithHandlerFunc prefixes the request with Middleware
func PrefixRequestWithHandlerFunc(value string, h http.HandlerFunc){
  router.PathPrefix(value).HandlerFunc(h)
}

// CreateSubrouterWithMiddleware provides a subrouter for use for middleware
func CreateSubrouterWithMiddleware(
  value string,
  h... negroni.Handler,
) * mux.Router {

  r := mux.NewRouter().PathPrefix(value).Subrouter().StrictSlash(true)

  handlers := make([]negroni.Handler, len(h)+1)
  copy(handlers, h)
  handlers[ len(h) ] = negroni.Wrap(r)

  PrefixRequestWithHandler(
    value,
    negroni.New(handlers...),
  )

  return r
}

// MapRequestToSubRouter maps a request to the sub router
func MapRequestToSubRouter(
  r *mux.Router,
  path string,
  method string,
  handler http.Handler,
) *mux.Route {
  return r.Handle(path, handler).Methods(method)
}

// MapRequest maps a request to the main router
func MapRequest(
  path string,
  method string,
  handler http.Handler,
) *mux.Route {
  return MapRequestToSubRouter(router, path, method, handler)
}

// MapRequestToSubRouterFunc maps a request to the sub router
func MapRequestToSubRouterFunc(
  r *mux.Router,
  path string,
  method string,
  f http.HandlerFunc,
) *mux.Route {
  return r.HandleFunc(path, f).Methods(method)
}

// MapRequestFunc maps a request to the main router
func MapRequestFunc(
  path string,
  method string,
  f http.HandlerFunc,
) *mux.Route {
  return MapRequestToSubRouterFunc(router, path, method, f)
}
