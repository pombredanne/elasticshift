/*
Copyright 2017 The Elasticshift Authors.
*/
package user

import (
	"errors"
	"fmt"
	"strings"

	"net/mail"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/dex"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var (
	errInvalidEmail              = errors.New("Invalid email address")
	errBlackPassword             = errors.New("Password cannot be empty")
	errBlankFullname             = errors.New("Fullname cannot be empty")
	errUserCreationFailed        = errors.New("Failed to create user")
	errInvalidEmailOrPassword    = errors.New("Invalid email or password")
	errUsernameOrPasswordIsEmpty = errors.New("User or password is empty ")
	errUserAlreadyExist          = errors.New("User already exists")
)

type server struct {
	logger logrus.FieldLogger
	dex    dex.DexClient
}

// NewServer ..
// Implementation of api.UserServer
func NewServer(logger logrus.FieldLogger, dex dex.DexClient) api.UserServer {
	return &server{
		logger: logger,
		dex:    dex,
	}
}

func (s server) SignUp(ctx context.Context, req *api.SignUpReq) (*api.SignUpRes, error) {

	fmt.Println("SignUp Request", req)

	resp := &api.SignUpRes{Created: false}

	// email validation
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return resp, errInvalidEmail
	}

	// password check
	if req.Password == "" {
		return resp, errBlackPassword
	}

	if req.Fullname == "" {
		return resp, errBlankFullname
	}

	in := &dex.CreatePasswordReq{}
	in.Password = &dex.Password{}
	// strip username from email
	in.Password.UserId = strings.Split(req.Email, "@")[0]

	// generate hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resp, errUserCreationFailed
	}
	in.Password.Hash = hashedPwd
	in.Password.Email = req.Email
	in.Password.Username = req.Fullname

	out, err := s.dex.CreatePassword(ctx, in)
	if err != nil {

		if out != nil && out.AlreadyExists {
			return nil, errUserAlreadyExist
		} else if strings.ContainsAny(err.Error(), "UNIQUE constraint") {
			return nil, errUserAlreadyExist
		}

		return nil, errUserCreationFailed
	}

	resp.Created = true

	// TODO send email for verification
	return resp, err
}

// SignIn ..
func (s server) SignIn(ctx context.Context, req *api.SignInReq) (*api.SignInRes, error) {

	// validate username and password
	if req.Username == "" || req.Password == "" {
		return nil, errUsernameOrPasswordIsEmpty
	}

	// u, err := s.store.GetUser(req.Username, req.Team)
	// if err != nil {
	// 	return nil, errInvalidEmailOrPassword
	// }

	// err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	// if err != nil {
	// 	return nil, errInvalidEmailOrPassword
	// }

	return &api.SignInRes{Valid: true}, nil
}
