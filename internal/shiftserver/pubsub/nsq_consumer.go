/*
Package pubsub ..

Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"github.com/nsqio/go-nsq"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

const (
	// CONCURRENCY ..
	// max concurreny the listener should be
	CONCURRENCY = 200
)

// Consumer ..
type Consumer interface {
	Subscribe(topic, channel string) error
	Unsubscribe()
}

type consumer struct {
	logger logrus.Logger
	sh     SubscriptionHandler
	cfg    NSQConfig
	schema *graphql.Schema

	ins  *nsq.Consumer
	conf *nsq.Config

	topic   string
	channel string
}

// NewConsumer ..
func NewConsumer(cfg NSQConfig, sh SubscriptionHandler, logger logrus.Logger, schema *graphql.Schema) Consumer {

	// TODO connect to topic broker
	c := &consumer{}
	c.sh = sh
	c.logger = logger
	c.cfg = cfg
	c.schema = schema

	return c
}

func (c *consumer) Subscribe(topic, channel string) error {

	c.topic = topic
	c.channel = channel

	conf := nsq.NewConfig()
	c.conf = conf

	consumer, err := nsq.NewConsumer(topic, channel, conf)
	if err != nil {
		fmt.Printf("Errror when creating new nsq consumer: %v", err)
		return err
	}
	c.ins = consumer

	consumer.AddHandler(c.HandleMessage())
	err = consumer.ConnectToNSQLookupd(c.cfg.Consumer.Address)
	if err != nil {
		fmt.Printf("Error when trying to launch consumer: %v\n", err)
	}

	<-consumer.StopChan

	return nil
}

func (c *consumer) Unsubscribe() {
	if c.ins != nil {
		c.ins.Stop()
	}
}

func (c *consumer) HandleMessage() nsq.Handler {

	return nsq.HandlerFunc(func(msg *nsq.Message) error {

		m, err := NewMessage(msg)
		if err != nil {
			return err
		}

		c.logger.Infoln("Incoming message: id=%s, payload=%s", m.Topic, m.Payload)

		// Get all subscripts from handler
		subscriptions := c.sh.Subscriptions()

		// push the results to the subscribers.
		for conn, _ := range subscriptions {

			for _, subscription := range subscriptions[conn] {

				if subscription.TopicID == m.Topic {

					// Prepare an execution context for running the query
					ctx := context.Background()

					fmt.Println("Query = ", subscription.Query)
					var name string
					name = subscription.OperationName
					if name == "" {
						name = c.channel
					}
					fmt.Println("OperationName = ", subscription.OperationName)
					fmt.Println("Chanel name = ", name)
					fmt.Println("Vars = ", subscription.Variables)

					// Re-execute the subscription query
					params := graphql.Params{
						Schema:         *c.schema, // The GraphQL schema
						RequestString:  subscription.Query,
						VariableValues: subscription.Variables,
						OperationName:  name,
						Context:        ctx,
					}
					result := graphql.Do(params)

					// Send query results back to the subscriber at any point
					data := SubscriptionResponse{
						// Data can be anything (interface{})
						Data: result.Data,
						// Errors is optional ([]error)
						Errors: utils.GraphQLErrors(result.Errors),
					}
					subscription.Push(&data)
				}
			}
		}
		// invoke graphql and send websocket response

		fmt.Println("Finished callback execution.")
		return nil
	})
}
