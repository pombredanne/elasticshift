package store

import (
	"testing"
	"time"

	"github.com/elasticshift/elasticshift/pkg/testhelper"
)

func TestConnectFail(t *testing.T) {

	l, _ := testhelper.GetLoggr()
	duration, _ := time.ParseDuration("3s")
	retryDur, _ := time.ParseDuration("1s")
	c := Config{
		"127.0.0.1",
		"test",
		"test",
		"test",
		duration,
		true,
		false,
		5,
		5,
		false,
		retryDur,
	}

	_, err := Connect(l, c)
	if err == nil {
		t.Log("Suppose to fail during connect.")
		t.Fail()
	}
}

func TestConnectAutoReconnectTimeout(t *testing.T) {

	l, _ := testhelper.GetLoggr()
	duration, _ := time.ParseDuration("5s")
	retryDur, _ := time.ParseDuration("1s")
	c := Config{
		"127.0.0.1",
		"test",
		"test",
		"test",
		duration,
		true,
		true,
		5,
		5,
		false,
		retryDur,
	}

	Connect(l, c)

	time.Sleep(6 * time.Second)
}
