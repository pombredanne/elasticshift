/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
)

type Consumers interface {
	Schema(schema *graphql.Schema)

	Add(topic string, s *Subscription) error
	Remove(topic string, s *Subscription)
}

type consumers struct {
	cfg    NSQConfig
	logger *logrus.Entry
	schema *graphql.Schema

	// topic: id[consumer]
	subscriptions map[string]map[string][]*Subscription
	topics        map[string]Consumer

	loggr logger.Loggr

	mutex sync.RWMutex
}

func NewConsumers(cfg NSQConfig, loggr logger.Loggr) Consumers {

	return &consumers{
		cfg:           cfg,
		logger:        loggr.GetLogger("pubsub/consumers"),
		subscriptions: make(map[string]map[string][]*Subscription),
		topics:        make(map[string]Consumer),
		loggr:         loggr,
	}
}

func (c *consumers) Schema(schema *graphql.Schema) {
	c.schema = schema
}

func (c *consumers) Add(topic string, s *Subscription) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.subscriptions[topic] == nil {
		c.subscriptions[topic] = make(map[string][]*Subscription)
	}

	if c.subscriptions[topic][s.OperationID] == nil {
		c.subscriptions[topic][s.OperationID] = []*Subscription{s}
	} else {
		c.subscriptions[topic][s.OperationID] = append(c.subscriptions[topic][s.OperationID], s)
	}

	go func() {
		if _, ok := c.topics[topic]; !ok {

			consumer := NewConsumer(c.cfg, c.loggr.GetLogger("pubsub/consumers/"+topic), c.schema, c)
			c.topics[topic] = consumer

			err := consumer.Subscribe(topic)
			if err != nil {
				c.logger.Errorf("Failed to subscribe : %v\n", err)
			}
		}
	}()

	return nil
}

func (c *consumers) Remove(topic string, s *Subscription) {

	c.logger.WithFields(logrus.Fields{"topic": topic, "OperationID": s.OperationID}).Infoln("Remove consumer")

	c.mutex.Lock()
	defer c.mutex.Unlock()

	subs := c.subscriptions[topic][s.OperationID]
	for i, sub := range subs {

		if s.Conenction == sub.Conenction && s.ID == sub.ID {
			subs = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	c.subscriptions[topic][s.OperationID] = subs
}
