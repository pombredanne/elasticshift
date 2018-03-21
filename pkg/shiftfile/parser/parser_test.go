package parser

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {

	tests := []struct {
		Filename    string
		ExpectedErr bool
	}{
		{
			"version.shift",
			false,
		},
		{
			"name.shift",
			false,
		},
		{
			"variables.shift",
			false,
		},
		{
			"image.shift",
			false,
		},
		{
			"block.shift",
			false,
		},
		{
			"hint.shift",
			false,
		},
		{
			"command.shift",
			false,
		},
		{
			"list.shift",
			false,
		},
		{
			"file.shift",
			false,
		},
		{
			"filenoc.shift",
			false,
		},
	}

	testfileDir := "./testfiles"

	for _, test := range tests {

		t.Run(test.Filename, func(t *testing.T) {
			buf, e := ioutil.ReadFile(filepath.Join(testfileDir, test.Filename))
			if e != nil {
				t.Fatalf("err: %s", e)
			}

			p := New(buf)

			_, err := p.Parse()

			if (err != nil) != test.ExpectedErr {
				t.Fatalf("Input: %s\n\nError: %s\n\nAST: %#v", test.Filename, err, p)
			}

			// fmt.Println(fmt.Sprintf("File node: %#v", f.Node))

			// if f.Version != nil {
			// 	fmt.Printf("Version %s\n", f.Version.Value)
			// }

			// if f.Name != nil {
			// 	fmt.Printf("Name %s\n", f.Name.Value)
			// }

			// if len(f.Variables) > 0 {
			// 	fmt.Printf("Variables %v\n", f.Variables)
			// }

			// fmt.Println(Nodes())

		})
	}
}
