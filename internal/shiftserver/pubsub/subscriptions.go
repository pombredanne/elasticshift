/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

var (
	topicNameFormat = "gqls-%s-%s"
)

// SubscriptionHandler ..
type SubscriptionHandler interface {
	Subcribe(Connection, *Subscription) []error
	UnsubscribeAll(Connection)
	Unsubscribe(Connection, *Subscription)
	Subscriptions() Subscriptions

	Schema(schema *graphql.Schema)
}

type subscriptionHandler struct {
	schema        *graphql.Schema
	subscriptions Subscriptions
	logger        *logrus.Entry
}

// PushResponseFunc ..
type PushResponseFunc func(*SubscriptionResponse)

// Subscription ..
type Subscription struct {
	ID            string
	Query         string
	Variables     map[string]interface{}
	OperationName string
	Conenction    Connection
	Push          PushResponseFunc
	TopicID       string
}

// ConnectionSubscriptions ..
type ConnectionSubscriptions map[string]*Subscription

// Subscriptions ..
type Subscriptions map[Connection]ConnectionSubscriptions

// NewSubscriptionHandler ..
func NewSubscriptionHandler(loggr logger.Loggr) SubscriptionHandler {

	sh := &subscriptionHandler{}
	sh.logger = loggr.GetLogger("graphql/subscriptions")
	sh.subscriptions = make(Subscriptions)
	return sh
}

func (sh *subscriptionHandler) Schema(schema *graphql.Schema) {
	sh.schema = schema
}

func (sh *subscriptionHandler) Subcribe(c Connection, s *Subscription) []error {

	if errors := validate(s); len(errors) > 0 {
		return errors
	}

	doc, err := parser.Parse(parser.ParseParams{
		Source: s.Query,
	})
	if err != nil {
		return []error{err}
	}

	validation := graphql.ValidateDocument(sh.schema, doc, nil)
	if !validation.IsValid {
		return utils.GraphQLErrors(validation.Errors)
	}

	if sh.subscriptions[c] == nil {
		sh.subscriptions[c] = ConnectionSubscriptions{}

	}

	if sh.subscriptions[c][s.ID] != nil {
		return []error{errors.New("Cannot register subscription twice")}
	}

	// TODO update to use topic
	opdef := doc.Definitions[0].(*ast.OperationDefinition)
	selection := opdef.GetSelectionSet().Selections[0]

	var topic string
	switch selection.(type) {
	case *ast.Field:
		f := selection.(*ast.Field)
		subName := f.Name.Value
		id := s.Variables["id"].(string)
		topic = fmt.Sprintf(topicNameFormat, subName, id)

		c.SetTopicName(topic, subName)
	}
	s.TopicID = topic
	sh.subscriptions[c][s.ID] = s

	return nil
}

func (sh *subscriptionHandler) UnsubscribeAll(c Connection) {

	// TODO update to use topic
	sub := sh.subscriptions[c]
	if sub != nil {

		for id := range sub {
			sh.Unsubscribe(c, sub[id])
		}

		delete(sh.subscriptions, c)
	}
}

func (sh *subscriptionHandler) Unsubscribe(c Connection, s *Subscription) {

	// TODO update to use topic
	delete(sh.subscriptions[c], s.ID)
	if len(sh.subscriptions[c]) == 0 {
		delete(sh.subscriptions, c)
	}
}

func (sh *subscriptionHandler) Subscriptions() Subscriptions {
	return sh.subscriptions
}

func validate(s *Subscription) []error {

	errs := []error{}

	if s.ID == "" {
		errs = append(errs, errors.New("Subscription ID is empty"))
	}

	if s.Conenction == nil {
		errs = append(errs, errors.New("Subscription is not associated with a connection"))
	}

	if s.Query == "" {
		errs = append(errs, errors.New("Subscription query is empty"))
	}

	if s.Push == nil {
		errs = append(errs, errors.New("Subscription has no push/callback function set"))
	}

	return errs
}
