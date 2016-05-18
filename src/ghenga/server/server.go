package server

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// NewRouter returns a new router with the complete ghenga API already attached.
func NewRouter(ctx context.Context, env *Env) *mux.Router {
	router := mux.NewRouter()
	PeopleHandler(ctx, env, router)
	LoginHandler(ctx, env, router)
	SearchHandler(ctx, env, router)
	UserHandler(ctx, env, router)
	return router
}
