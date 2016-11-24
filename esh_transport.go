package esh

import (
	"net/http/pprof"

	ghandlers "github.com/gorilla/handlers"
	"github.com/justinas/alice"
	"gitlab.com/conspico/esh/core/handlers"
)

// MakeHandlers ..
// Create application specific handlers
func MakeHandlers(ctx AppContext) {

	corsOpts := ghandlers.AllowedOrigins([]string{"*"})
	corsHandler := ghandlers.CORS(corsOpts)
	recoveryHandler := ghandlers.RecoveryHandler()

	commonChain := alice.New(recoveryHandler, corsHandler)

	extractHandler := handlers.ExtractHandler(ctx.Context, ctx.Router)
	ctx.PublicChain = commonChain.Extend(alice.New(extractHandler))

	secureHandler := handlers.SecurityHandler(ctx.Context, ctx.Logger, ctx.Signer, ctx.Verifier)
	ctx.SecureChain = commonChain.Extend(alice.New(secureHandler, extractHandler))

	MakeTeamHandler(ctx)
	MakeUserHandler(ctx)
	MakeVCSHandler(ctx)
	MakeRepoHandler(ctx)

	// pprof
	ctx.Router.HandleFunc("/debug/pprof", pprof.Index)
	ctx.Router.HandleFunc("/debug/symbol", pprof.Symbol)
	ctx.Router.HandleFunc("/debug/profile", pprof.Profile)
	ctx.Router.Handle("/debug/heap", pprof.Handler("heap"))
	ctx.Router.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	ctx.Router.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	ctx.Router.Handle("/debug/block", pprof.Handler("block"))
}

// MakeTeamHandler ..
func MakeTeamHandler(ctx AppContext) {

	/** Team **/
	createTeamHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeCreateTeamRequest,
		ProcessFunc: makeCreateTeamEdge(ctx.TeamService),
		EncodeFunc:  encodeCreateTeamResponse,
		Logger:      ctx.Logger,
	}
	ctx.Router.Handle("/api/teams", ctx.PublicChain.Then(createTeamHandler)).Methods("POST")
}

// MakeUserHandler ..
func MakeUserHandler(ctx AppContext) {

	r := ctx.Router

	/** User **/
	signUpHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeSignUpRequest,
		EncodeFunc:  encodeSignInResponse,
		ProcessFunc: makeSignupEdge(ctx.UserService),
		Logger:      ctx.Logger,
	}
	r.Handle("/api/users/signup", ctx.PublicChain.Then(signUpHandler)).Methods("POST")

	signInHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeSignInRequest,
		EncodeFunc:  encodeSignInResponse,
		ProcessFunc: makeSignInEdge(ctx.UserService),
		Logger:      ctx.Logger,
	}
	r.Handle("/api/users/signin", ctx.PublicChain.Then(signInHandler)).Methods("POST")

	signOutHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeSignOutRequest,
		EncodeFunc:  encodeSignOutResponse,
		ProcessFunc: makeSignOutEdge(ctx.UserService),
		Logger:      ctx.Logger,
	}
	r.Handle("/api/users/signout", signOutHandler).Methods("POST")

}

// MakeVCSHandler ..
func MakeVCSHandler(ctx AppContext) {

	r := ctx.Router

	/** VCS **/
	authorizeHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeAuthorizeRequest,
		EncodeFunc:  encodeAuthorizeResponse,
		ProcessFunc: makeAuthorizeEdge(ctx.VCSService),
		Logger:      ctx.Logger,
	}
	r.Handle("/api/auth/{provider}", ctx.SecureChain.Then(authorizeHandler)).Methods("GET")

	authorizedHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeAuthorizedRequest,
		EncodeFunc:  encodeAuthorizeResponse,
		ProcessFunc: makeAuthorizedEdge(ctx.VCSService),
		Logger:      ctx.Logger,
	}
	r.Handle("/api/auth/{provider}/callback", ctx.PublicChain.Then(authorizedHandler)).Methods("GET")

	getVCSHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeGetVCSRequest,
		EncodeFunc:  encodeGetVCSResponse,
		ProcessFunc: makeGetVCSEdge(ctx.VCSService),
		Logger:      ctx.Logger,
	}
	r.Handle("/api/vcs", ctx.SecureChain.Then(getVCSHandler)).Methods("GET")

	syncVCSHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeSyncVCSRequest,
		EncodeFunc:  encodeSyncVCSResponse,
		ProcessFunc: makeSyncVCSEdge(ctx.VCSService),
		Logger:      ctx.Logger,
	}
	r.Handle("/api/vcs/sync/{id}", ctx.SecureChain.Then(syncVCSHandler)).Methods("GET")
}

// MakeRepoHandler ..
func MakeRepoHandler(ctx AppContext) {

	r := ctx.Router

	/** Repo **/
	getRepoHandler := &handlers.RequestHandler{
		DecodeFunc:  decodeGetRepoRequest,
		EncodeFunc:  encodeGetRepoResponse,
		ProcessFunc: makeGetRepoEdge(ctx.RepoService),
		Logger:      ctx.Logger,
	}

	r.Handle("/api/repos", ctx.SecureChain.Then(getRepoHandler)).Methods("GET")
	r.Handle("/api/vcs/{id}/repos", ctx.SecureChain.Then(getRepoHandler)).Methods("GET")
}
