/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"github.com/Sirupsen/logrus"
	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
)

// Producer ..
type Producer interface {
	Publish(topic string, payload interface{}) error
}

type producer struct {
	logger logrus.Logger
	sh     SubscriptionHandler
	cfg    NSQConfig

	ins *nsq.Producer
}

// NewProducer ..
func NewProducer(cfg NSQConfig, sh SubscriptionHandler, logger logrus.Logger) (Producer, error) {

	// TODO connect to topic broker
	p := &producer{}
	p.sh = sh
	p.logger = logger
	p.cfg = cfg

	conf := nsq.NewConfig()
	ins, err := nsq.NewProducer(cfg.Producer.Address, conf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init nsq producer.")
	}
	p.ins = ins
	return p, nil
}

func (p *producer) Publish(topic string, payload interface{}) error {

	msg := Message{Topic: topic, Payload: payload}

	encoded, err := msg.encode()
	if err != nil {
		return err
	}

	return p.ins.Publish(topic, encoded)
}

func (p *producer) PublishAsync(topic string, payload interface{}) error {

	msg := Message{Topic: topic, Payload: payload}
	encoded, err := msg.encode()
	if err != nil {
		return err
	}

	resChan := make(chan *nsq.ProducerTransaction, 1)
	go func(resChan chan *nsq.ProducerTransaction) {

		for {
			trans, ok := <-resChan
			if ok && trans.Error != nil {
				p.logger.Fatalf(trans.Error.Error())
			}
		}
	}(resChan)

	return p.ins.PublishAsync(topic, encoded, resChan, "")
}
