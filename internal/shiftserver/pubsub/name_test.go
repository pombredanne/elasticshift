/*
Copyright 2018 The Elasticshift Authors.
*/
package pubsub

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

func testName(t *testing.T) {

	query := `subscription test {
			buildUpdateSubscibe (id : "12345"){
				id
				graph
				status
			}
		}`

	doc, err := parser.Parse(parser.ParseParams{
		Source: query,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("kind = %s\n", doc.GetKind())
	fmt.Printf("defs = %#v\n", doc.Definitions)

	opdef := doc.Definitions[0].(*ast.OperationDefinition)

	fmt.Printf("Name = %s\n", opdef.GetName().Value)

	selections := opdef.GetSelectionSet()
	fmt.Printf("Kind = %s, selections = %#v\n", selections.Kind, selections)

	s := selections.Selections[0]
	switch s.(type) {
	case *ast.Field:
		f := s.(*ast.Field)
		fmt.Printf("field.name = %#v\n", f.Name.Value)
		fmt.Printf("Arguments-name = %#v\n", f.Arguments[0].Name.Value)
		fmt.Printf("Arguments-value = %#v", f.Arguments[0].Value.GetValue())
	}

	// selections2 := selections.Selections[0].GetSelectionSet()
	// fmt.Printf("Kind = %s, selections2 = %#v\n", selections2.Kind, selections2)

	// printSelection(opdef.GetSelectionSet())
}

func printSelection(selections *ast.SelectionSet) {

	for k, v := range selections.Selections {
		fmt.Printf("Kind = %d, Type= %q, selections = %v\n", k, k, k)
		printSelection(v.GetSelectionSet())
	}

}
