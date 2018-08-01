/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/pkg/dispatch"
)

var (
	topichHttpUrl  = "http://%s/topic/create"
	topichHttpsUrl = "https://%s/topic/create"
)

type engine struct {
	logger    logrus.Logger
	sh        SubscriptionHandler
	schema    *graphql.Schema
	conf      NSQConfig
	producers map[string]Producer
}

// Engine ..
type Engine interface {
	Producer() (Producer, error)
	Consumer() Consumer
	SubscriptionHandler() SubscriptionHandler

	Publish(subscriptionName, id string) error

	CreateTopic(name, id string) error
	Schema(schema *graphql.Schema)
}

// NewEngine ..
func NewEngine(logger logrus.Logger, sh SubscriptionHandler, conf NSQConfig) Engine {

	return &engine{
		logger:    logger,
		sh:        sh,
		conf:      conf,
		producers: make(map[string]Producer),
	}
}

func (e *engine) Schema(schema *graphql.Schema) {
	e.schema = schema
	e.sh.Schema(schema)
}

func (e *engine) Producer() (Producer, error) {
	return NewProducer(e.conf, e.sh, e.logger)
}

func (e *engine) Consumer() Consumer {
	return NewConsumer(e.conf, e.sh, e.logger, e.schema)
}

func (e *engine) SubscriptionHandler() SubscriptionHandler {
	return e.sh
}

func (e *engine) Publish(subscriptionName, id string) error {

	topic := fmt.Sprintf(topicNameFormat, subscriptionName, id)
	prod, err := e.producer(topic)
	if err != nil {
		return fmt.Errorf("Failed to get producer, can't publish the to topic [%s] : %v", topic, err)
	}
	return prod.Publish(topic, "")
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

func (e *engine) CreateTopic(name, id string) error {

	// TODO run with https
	r := dispatch.NewPostRequestMaker(fmt.Sprintf(topichHttpUrl, e.conf.Producer.Address))
	r.QueryParam("topic", fmt.Sprintf(topicNameFormat, name, id))
	r.SetLogger(e.logger)
	r.Verbose(true)
	r.UnescapeQueryParams(true)

	return r.Dispatch()
}
