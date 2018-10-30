/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
)

var (
	topichHttpUrl  = "http://%s/topic/create"
	topichHttpsUrl = "https://%s/topic/create"
)

type engine struct {
	loggr     logger.Loggr
	logger    *logrus.Entry
	sh        SubscriptionHandler
	schema    *graphql.Schema
	conf      NSQConfig
	producers map[string]Producer
	consumers Consumers
}

// Engine ..
type Engine interface {
	Producer() (Producer, error)
	Consumers() Consumers
	SubscriptionHandler() SubscriptionHandler

	Publish(topic string, payload interface{}) error

	Schema(schema *graphql.Schema)
}

// NewEngine ..
func NewEngine(loggr logger.Loggr, sh SubscriptionHandler, conf NSQConfig, cons Consumers) Engine {

	l := loggr.GetLogger("pubsub/engine")
	return &engine{
		logger:    l,
		sh:        sh,
		conf:      conf,
		loggr:     loggr,
		producers: make(map[string]Producer),
		consumers: cons,
	}
}

func (e *engine) Schema(schema *graphql.Schema) {
	e.schema = schema
	e.consumers.Schema(schema)
	if e.sh != nil {
		e.sh.Schema(schema)
	}
}

func (e *engine) Producer() (Producer, error) {
	return NewProducer(e.conf, e.sh, e.loggr)
}

func (e *engine) Consumers() Consumers {
	return e.consumers
}

// func (e *engine) Consumer() Consumer {
// 	mh := NewMessageHandler(e.consumers)
// 	return NewConsumer(e.conf, mh, e.loggr, e.schema)
// }

func (e *engine) SubscriptionHandler() SubscriptionHandler {
	return e.sh
}

func (e *engine) Publish(topic string, payload interface{}) error {

	// topic := fmt.Sprintf(topicNameFormat, subscriptionName, payload)
	prod, err := e.producer(topic)
	if err != nil {
		return fmt.Errorf("Failed to get producer, can't publish the to topic [%s] : %v", topic, err)
	}
	return prod.Publish(topic, payload)
}

func (e *engine) producer(topic string) (Producer, error) {

	var prod Producer
	var err error

	// check if the producer exist for the topic
	//prod = e.producers[topic]
	if prod == nil {

		// create producer
		prod, err = e.Producer()
		if err != nil {
			return nil, err
		}
		e.producers[topic] = prod
	}
	return prod, nil
}

// func (e *engine) CreateTopic(name, id string) error {

// 	// TODO run with https
// 	r := dispatch.NewPostRequestMaker(fmt.Sprintf(topichHttpUrl, e.conf.Producer.Address))
// 	r.QueryParam("topic", fmt.Sprintf(topicNameFormat, name, id))
// 	r.SetLogger(e.logger)
// 	r.Verbose(true)
// 	r.UnescapeQueryParams(true)

// 	return r.Dispatch()
// }
