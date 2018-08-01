/*

Package pubsub ...

Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"encoding/json"

	nsq "github.com/nsqio/go-nsq"
)

// SubscriptionRequest ..
type SubscriptionRequest struct {
	OperationName string                 `json:"operation_name,omitempty"`
	Query         string                 `json:"query,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// SubscriptionResponse ..
type SubscriptionResponse struct {
	Data   interface{} `json:"data"`
	Errors []error     `json:"errors"`
}

// WebsocketMessage ..
type WebsocketMessage struct {
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func (msg WebsocketMessage) String() string {
	s, _ := json.Marshal(msg)
	if s != nil {
		return string(s)
	}
	return ""
}

// NSQConfig ..
type NSQConfig struct {
	Consumer struct {
		Address string
		Config  *nsq.Config
	}
	Producer struct {
		Address string
		Config  *nsq.Config
	}
}

// CosumerHandleFunc ..
type CosumerHandleFunc func(m Message) error

// ConsumerConfig ..
type ConsumerConfig struct {
	NSQConfig   NSQConfig
	HandlerFunc CosumerHandleFunc
	Topic       string
	Channel     string
}

// Message ..
type Message struct {
	//nsq.Message
	Topic   string      `json:"topic,omitempty"` // topic to publish or consume, as well as subscription name
	Payload interface{} `json:"payload,omitempty"`
}

func (m *Message) encode() ([]byte, error) {
	return json.Marshal(m)
}

// NewMessage ..
func NewMessage(msg *nsq.Message) (*Message, error) {
	var m Message
	err := json.Unmarshal(msg.Body, &m)
	return &m, err
}
