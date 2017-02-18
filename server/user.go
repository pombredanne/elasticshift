// Package server ..
// Author Ghazni Nattarshah
// Date: 1/11/17
package server

import (
	"errors"
	"fmt"
	"strings"

	"net/mail"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/pb"
	"gitlab.com/conspico/elasticshift/store"
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
)

type userServer struct {
	store  store.UserStore
	cache  *Cache
	logger logrus.FieldLogger
}

// NewUserServer ..
// Implementation of pb.UserServer
func NewUserServer(s *Server) pb.UserServer {
	return &userServer{
		store:  store.NewUserStore(s.Store),
		cache:  s.Cache,
		logger: s.Logger,
	}
}

func (s userServer) SignUp(ctx context.Context, req *pb.SignUpReq) (*pb.SignUpRes, error) {

	fmt.Println("SignUp Request", req)

	resp := &pb.SignUpRes{Created: false}

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

	// strip username from email
	userName := strings.Split(req.Email, "@")[0]

	// generate hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resp, errUserCreationFailed
	}

	user := &store.User{
		Fullname:      req.Fullname,
		Username:      userName,
		Email:         req.Email,
		Password:      string(hashedPwd[:]),
		Locked:        false,
		Active:        true,
		BadAttempt:    0,
		EmailVefified: false,
		Team:          req.Team,
	}

	err = s.store.Insert(user)

	resp.Created = err == nil

	// TODO send email for verification

	return resp, err
}

// SignIn ..
func (s userServer) SignIn(ctx context.Context, req *pb.SignInReq) (*pb.SignInRes, error) {

	// validate username and password
	if req.Username == "" || req.Password == "" {
		return nil, errUsernameOrPasswordIsEmpty
	}

	u, err := s.store.GetUser(req.Username, req.Team)
	if err != nil {
		return nil, errInvalidEmailOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return nil, errInvalidEmailOrPassword
	}

	return &pb.SignInRes{Valid: true}, nil
}
