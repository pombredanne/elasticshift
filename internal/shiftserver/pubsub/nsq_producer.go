/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
)

// Producer ..
type Producer interface {
	Publish(topic string, payload interface{}) error
}

type producer struct {
	logger *logrus.Entry
	sh     SubscriptionHandler
	cfg    NSQConfig

	ins *nsq.Producer
}

// NewProducer ..
func NewProducer(cfg NSQConfig, sh SubscriptionHandler, loggr logger.Loggr) (Producer, error) {

	// TODO connect to topic broker
	p := &producer{}
	p.sh = sh
	p.logger = loggr.GetLogger("pubsub/producer")
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

	p.logger.Infof("Publishig to topic '%s': [%v]\n", topic, payload)

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
