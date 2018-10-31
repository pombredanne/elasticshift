package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/keys"
)

func TestImageBlock(t *testing.T) {

	type Expected struct {
		Key   string
		Value string
	}

	tests := []struct {
		Filename string
		Expected []Expected
	}{
		{
			"image.shift",
			[]Expected{
				{keys.NAME, "elasticshift/java:1.9"},
				{"registry", "http://dockerregistry.com/elasticshift"},
				{"secret", "askjdflkjahsdlkjfhlkjs"},
			},
		},
		{
			"image2.shift",
			[]Expected{
				{keys.NAME, "elasticshift/java:1.9"},
			},
		},
		{
			"image3.shift",
			[]Expected{
				{keys.NAME, "elasticshift/java,elasticshift/java6,elasticshift/java7"},
			},
		},
		{
			"image4.shift",
			[]Expected{
				{keys.NAME, "openjdk:7,openjdk:8"},
			},
		},
	}

	for _, test := range tests {

		t.Run(test.Filename, func(t *testing.T) {

			f := load(test.Filename, t)
			img := f.Image()

			for _, exp := range test.Expected {
				assertString(t, exp.Value, img[exp.Key].(string))
			}
		})
	}
}

func load(name string, t *testing.T) *ast.File {

	testfileDir := "./testfiles"

	buf, e := ioutil.ReadFile(filepath.Join(testfileDir, name))
	if e != nil {
		t.Fatalf("err: %s", e)
	}

	p := New(buf)

	f, err := p.Parse()

	if err != nil {
		t.Fatalf("Failed %v", err)
	}
	return f
}
