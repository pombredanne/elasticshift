package server

import (
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/store"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestServer(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name       string
		configFunc func(c *Config)
		err        bool
	}{
		{
			name:       "Positive Test",
			configFunc: nil,
			err:        false,
		},
		{
			name: "No logger",
			configFunc: func(c *Config) {
				c.Logger = nil
			},
			err: true,
		},
	}

	for _, test := range tests {

		_, err := newTestServer(ctx, t, test.configFunc)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func newTestServer(ctx context.Context, t *testing.T, updateFunc func(conf *Config)) (*Server, error) {

	timeout, _ := time.ParseDuration("10s")

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &logrus.TextFormatter{DisableColors: true},
		Level:     logrus.DebugLevel,
	}

	c := Config{

		Store: store.Config{

			Server:        "127.0.0.1",
			Name:          "esh",
			Username:      "esh",
			Password:      "eshpazz",
			Monotonic:     true,
			Timeout:       timeout,
			AutoReconnect: false,
			RetryIn:       timeout,
		},

		Dex: Dex{
			HostAndPort: "127.0.0.1:5557",
			Issuer:      "http://127.0.0.1:5556/dex",
			ID:          "yyjw66rn2hso6wriuzlic62jiy",
			Secret:      "l77r6wixjjtgmo4iym2kmk3jcuuxetj3afnqaw5w3rnl5nu5hehu",
			RedirectURI: "http://127.0.0.1:5050/login/callback",
		},

		Logger: logger,
	}

	session, err := store.Connect(c.Logger, c.Store)
	if err != nil {

		logger.Errorf("Failed connect to database : %v", err)
		return nil, err
	}
	c.Session = session

	if updateFunc != nil {
		updateFunc(&c)
	}

	s, err := NewServer(ctx, c)
	if err != nil {

		logger.Errorf("Failed initialize the server : %v", err)
		return nil, err
	}

	return s, nil
}
