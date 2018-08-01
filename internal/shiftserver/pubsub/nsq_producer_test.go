/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"testing"

	nsq "github.com/nsqio/go-nsq"
)

func TestProducer(t *testing.T) {

	prodAddress := "127.0.0.1:4150"
	conf := nsq.NewConfig()
	ins, err := nsq.NewProducer(prodAddress, conf)
	if err != nil {
		panic(err)
	}

	m := Message{}
	m.Topic = "gqls-subscribe_build_update-5b61c0c7dc294a61ccb3d868"
	m.Payload = ""
	d, err := m.encode()
	if err != nil {
		panic(err)
	}
	err = ins.Publish("gqls-subscribe_build_update-5b61c0c7dc294a61ccb3d868", d)
	if err != nil {
		panic(err)
	}
}
