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
			name: "Invalid Issuer",
			configFunc: func(c *Config) {
				c.Issuer = "http://192.168.0.%31/"
			},
			err: true,
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
	idTokenExpiry, _ := time.ParseDuration("24h")
	signingKeysExpiry, _ := time.ParseDuration("6h")

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &logrus.TextFormatter{DisableColors: true},
		Level:     logrus.DebugLevel,
	}

	c := Config{

		Issuer: "http://armor.elasticshift.com",

		Store: store.Config{

			Server:        "127.0.0.1",
			Name:          "armor",
			Username:      "armor",
			Password:      "armorpazz",
			Monotonic:     true,
			Timeout:       timeout,
			AutoReconnect: false,
		},

		IDTokensLifeSpan:   idTokenExpiry,
		SignerKeysLifeSpan: signingKeysExpiry,

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
