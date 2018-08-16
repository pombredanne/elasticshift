/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
)

// NewHandler ..
func NewGraphqlWSHandler(engine Engine, loggr logger.Loggr) http.Handler {

	logger := loggr.GetLogger("graphql/wshandler")

	upgrader := websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{"graphql-ws"},
	}

	subscriptionHandler := engine.SubscriptionHandler()

	var connections = make(map[Connection]bool)

	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}

			if ws.Subprotocol() != "graphql-ws" {
				ws.Close()
				return
			}

			conn := newConnection(ws, logger, engine, EventHandler{
				Subscribe: func(c Connection, id string, req *SubscriptionRequest) []error {

					logger.WithFields(logrus.Fields{
						"connection": c.ID(),
						"id":         id,
					}).Debug("Start..")

					s := Subscription{
						ID:            id,
						OperationName: req.OperationName,
						Query:         req.Query,
						Variables:     req.Variables,
						Conenction:    c,
						Push: func(res *SubscriptionResponse) {
							c.PushData(id, res)
						},
					}

					errs := subscriptionHandler.Subcribe(c, &s)
					return errs
				},
				Unsubscribe: func(c Connection, id string) {
					subscriptionHandler.Unsubscribe(c, &Subscription{ID: id})
				},
				Close: func(c Connection) {

					subscriptionHandler.UnsubscribeAll(c)
					delete(connections, c)
				},
			})

			connections[conn] = true
		},
	)
}
