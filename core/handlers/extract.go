package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// extract ..
type extract struct {
	ctx    context.Context
	h      http.Handler
	router *mux.Router
}

func (eh *extract) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var match mux.RouteMatch
	if eh.router.Match(r, &match) {
		ctx = context.WithValue(ctx, "params", match.Vars)
	}

	eh.h.ServeHTTP(w, r.WithContext(ctx))
}

// ExtractHandler ..
func ExtractHandler(ctx context.Context, router *mux.Router) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {

		return &extract{
			ctx:    ctx,
			h:      h,
			router: router,
		}
	}
}
