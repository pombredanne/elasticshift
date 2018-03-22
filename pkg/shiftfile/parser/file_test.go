package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/ast"
)

func TestFile(t *testing.T) {

	testfileDir := "./testfiles"

	buf, e := ioutil.ReadFile(filepath.Join(testfileDir, "file.shift"))
	if e != nil {
		t.Fatalf("err: %s", e)
	}

	p := New(buf)

	f, err := p.Parse()

	if err != nil {
		t.Fatalf("Failed %v", err)
	}

	// fmt.Printf("File : %#v", f.Node)
	testVars(f, t)
}

func testVars(f *ast.File, t *testing.T) {

	vars := f.Vars()
	fmt.Printf("Variables : %v", vars)
}
