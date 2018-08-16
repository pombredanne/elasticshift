/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

const (
	gqlConnectionInit      = "connection_init"
	gqlConnectionAck       = "connection_ack"
	gqlConnectionKeepAlive = "ka"
	gqlConnectionError     = "connection_error"
	gqlConnectionTerminate = "connection_terminate"
	gqlSubscribe           = "start"
	gqlUnsubscribe         = "stop"
	gqlData                = "data"
	gqlError               = "error"
	gqlComplete            = "complete"

	// Maximum size of incoming messages
	readLimit = 4096

	// Timeout for outgoing messages
	writeTimeout = 10 * time.Second
)

// Connection ..
type connection struct {
	id     string
	logger *logrus.Entry
	ws     *websocket.Conn
	engine Engine

	mutex  *sync.Mutex
	push   chan WebsocketMessage
	closed bool
	eh     EventHandler

	consumers map[string]Consumer
}

// Connection ..
type Connection interface {
	ID() string

	PushData(id string, response *SubscriptionResponse)

	PushError(err error)
	PushSubscriptionError(id string, err []error)
}

// EventHandler ..
type EventHandler struct {
	Subscribe   func(Connection, string, *SubscriptionRequest) []error
	Unsubscribe func(Connection, string)
	Close       func(Connection)
}

func newConnection(ws *websocket.Conn, logger *logrus.Entry, engine Engine, eh EventHandler) Connection {

	c := &connection{}
	c.id = utils.NewUUID()
	c.ws = ws
	c.eh = eh
	c.logger = logger
	c.closed = false
	c.mutex = &sync.Mutex{}
	c.engine = engine

	c.push = make(chan WebsocketMessage)
	c.consumers = make(map[string]Consumer)

	go c.readLoop()
	go c.writeLoop()

	return c
}

func (c *connection) ID() string {
	return c.id
}

func (c *connection) PushData(id string, response *SubscriptionResponse) {

	msg := WebsocketMessage{}
	msg.Type = gqlData
	msg.ID = id
	msg.Payload = response

	c.mutex.Lock()
	if !c.closed {
		c.push <- msg
	}
	c.mutex.Unlock()
}

func (c *connection) PushError(err error) {

	msg := WebsocketMessage{}
	msg.Type = gqlError
	msg.Payload = err.Error()

	c.mutex.Lock()
	if !c.closed {
		fmt.Println("Pushing error..")
		c.push <- msg
	}
	c.mutex.Unlock()
}

func (c *connection) PushSubscriptionError(id string, err []error) {

	if c.closed {
		return
	}

	msg := WebsocketMessage{}
	msg.Type = gqlError
	msg.ID = id
	msg.Payload = err

	c.mutex.Lock()
	if !c.closed {
		c.push <- msg
	}
	c.mutex.Unlock()
}

func (c *connection) close() {

	c.mutex.Lock()
	c.closed = true
	close(c.push)
	c.mutex.Unlock()

	if c.eh.Close != nil {

		c.eh.Close(c)

		if c.consumers != nil {
			for _, cons := range c.consumers {
				cons.Unsubscribe()
			}
		}
	}
}

func (c *connection) readLoop() {

	defer c.ws.Close()

	c.ws.SetReadLimit(readLimit)

	for {

		rawMsg := json.RawMessage{}
		msg := WebsocketMessage{
			Payload: &rawMsg,
		}

		err := c.ws.ReadJSON(&msg)
		if err != nil {
			c.close()
			return
		}

		switch msg.Type {
		case gqlConnectionInit:
			c.push <- WebsocketMessage{
				Type: gqlConnectionAck,
			}

		case gqlSubscribe:

			if c.eh.Subscribe != nil {
				req := SubscriptionRequest{}
				err := json.Unmarshal(rawMsg, &req)
				if err != nil {
					c.PushError(errors.New("Invalid GQL_START request"))
					return
				}
				c.logger.WithFields(logrus.Fields{"ID": msg.ID, "Request": req}).Debugln("Subscription request")

				errs := c.eh.Subscribe(c, msg.ID, &req)
				if errs != nil {
					c.PushSubscriptionError(msg.ID, errs)
					return
				}
			}

		case gqlUnsubscribe:
			if c.eh.Unsubscribe != nil {
				c.eh.Unsubscribe(c, msg.ID)
				cons := c.consumers[msg.ID]
				if cons != nil {
					cons.Unsubscribe()
				}
			}
		case gqlConnectionTerminate:
			c.logger.Info("Connection closed by client")
			c.close()
			return
		}
	}
}

func (c *connection) writeLoop() {

	defer c.ws.Close()

	for {

		select {

		case msg, ok := <-c.push:

			if !ok {
				return
			}

			c.ws.SetWriteDeadline(time.Now().Add(writeTimeout))

			err := c.ws.WriteJSON(msg)
			if err != nil {
				return
			}
		}
	}
}
