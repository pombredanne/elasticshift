/*
Copyright 2017 The Elasticshift Authors.
*/
package store

import (
	"fmt"
	"time"

	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	mgo "gopkg.in/mgo.v2"

	"github.com/sirupsen/logrus"
)

// Config ..
type Config struct {
	Server        string
	Name          string
	Username      string
	Password      string
	Timeout       time.Duration
	Monotonic     bool
	AutoReconnect bool

	// old info
	IdleConnection int
	MaxConnection  int
	Log            bool
	RetryIn        time.Duration
}

// Connect ..
// Open the database connection and returns the session
func Connect(loggr logger.Loggr, cfg Config) (*mgo.Session, error) {

	logger := loggr.GetLogger("shiftdb")

	// DB Initialization
	var session *mgo.Session
	var err error
	//tryit := true
	// for tryit {

	logger.Infoln("Connecting to database...")
	session, err = dialMongo(cfg)

	if err != nil {

		if !cfg.AutoReconnect {
			return nil, err
		}
	} else {
		// set the configurations
		session.SetMode(mgo.Monotonic, cfg.Monotonic)
	}

	// 		logger.Errorln(fmt.Sprintf("Connecting database failed, retrying in %v [", cfg.RetryIn), err, "]")
	// 		time.Sleep(cfg.RetryIn)

	// 	} else {

	// 		// Ping function checks the database connectivity
	// 		dberr := session.Ping()
	// 		if dberr != nil {
	// 			logger.Errorln(fmt.Sprintf("Ping DB failed, retrying in %v []", cfg.RetryIn), err, "]")
	// 		} else {
	// 			logger.Infoln("Database connected successfully")
	// 			tryit = false
	// 		}
	// 	}
	// }

	// starting a background process to automatically reconnect to database in case of disconnects.
	if cfg.AutoReconnect {
		autoReconnect(logger, cfg, session)
	}
	return session, nil
}

// autoReconnect ..
// Launches a separate go routine to perform below operations.
// 1. Ping the db session to ensure the conenction is live
// 2. Retries to conenect with database if connectivity broke.
func autoReconnect(logger *logrus.Entry, cfg Config, session *mgo.Session) {

	ticker := time.NewTicker(cfg.RetryIn)
	timer := time.NewTimer(cfg.Timeout)

	// timeoutCh := make(chan int, 1)

	// go func() {
	// 	for range timer.C {
	// 		timeoutCh <- 1
	// 	}
	// }()

	go func() {

		var err error
		resetTimer := true
		for {
			select {
			case <-timer.C:
				logger.Errorf("Timed out after %s \n", cfg.Timeout)
				ticker.Stop()
				timer.Stop()
				return
			case <-ticker.C:

				if session == nil {
					session, err = dialMongo(cfg)
				} else {
					// Ping function checks the database connectivity
					err = session.Ping()
				}

				if err != nil {

					if resetTimer {
						timer.Reset(cfg.Timeout)
						resetTimer = false
					}
					logger.Errorln(fmt.Sprintf("DB ping failed, something went wrong. Reconnecting in %d seconds", cfg.RetryIn))

					//Trying to refresh the db connection
					//session.Refresh()

				} else {
					logger.Infoln("Reconnected to database successfully.")
					resetTimer = true
					timer.Stop()
				}
			}
		}
	}()
}

func dialMongo(cfg Config) (*mgo.Session, error) {

	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{cfg.Server},
		Username: cfg.Username,
		Password: cfg.Password,
		Database: cfg.Name,
		Timeout:  cfg.Timeout,
	})
}
