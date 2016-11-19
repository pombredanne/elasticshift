package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/esh/core/auth"
)

var (
	unAuthorized    = "Unauthorized"
	errUnauthorized = errors.New(unAuthorized)

	// AuthTokenCookieName ..
	AuthTokenCookieName = "__at"
)

type security struct {
	ctx      context.Context
	h        http.Handler
	logger   *logrus.Logger
	signer   interface{}
	verifier interface{}
}

func (sh *security) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	cookie, err := r.Cookie(AuthTokenCookieName)
	if err != nil {

		logError(sh.logger, err, w)
		return
	}

	if "" == cookie.Value {
		logError(sh.logger, err, w)
		return
	}

	token, err := auth.VefifyToken(sh.verifier, cookie.Value)
	if err != nil || !token.Valid {
		logError(sh.logger, err, w)
		return
	}

	c := context.WithValue(ctx, "token", auth.GetToken(token))

	// Refresh the token
	refreshtoken(sh.logger, token, sh.signer, w)

	sh.h.ServeHTTP(w, r.WithContext(c))
}

func refreshtoken(logger *logrus.Logger, token *jwt.Token, signer interface{}, w http.ResponseWriter) {

	signedTok, err := auth.RefreshToken(signer, token)
	if err != nil {
		logger.Errorln("Failed to refresh the token.", err)
	}

	cookie := &http.Cookie{
		Name:     AuthTokenCookieName,
		Value:    signedTok,
		Expires:  time.Now().Add(time.Minute * 15),
		HttpOnly: true,
		Path:     "/",
		//Secure : true, // TODO enable this to ensure the cookie is passed only with https
	}
	http.SetCookie(w, cookie)
}

func logError(logger *logrus.Logger, err error, w http.ResponseWriter) {

	logger.Errorln(err)
	http.Error(w, unAuthorized, http.StatusUnauthorized)
}

// SecurityHandler ..
func SecurityHandler(ctx context.Context, logger *logrus.Logger, signer interface{}, verifier interface{}) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {

		return &security{
			ctx:      ctx,
			signer:   signer,
			verifier: verifier,
			h:        h,
			logger:   logger,
		}
	}
}
