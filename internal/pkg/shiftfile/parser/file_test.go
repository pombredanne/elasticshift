package parser

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/keys"
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

	t.Run("vars", func(t *testing.T) {
		testVars(f, t)
	})

	t.Run("image", func(t *testing.T) {
		testImage(f, t)
	})

	t.Run("imagenames", func(t *testing.T) {
		testImageNames(f, t)
	})

	t.Run("blocks", func(t *testing.T) {
		testBlock(f, t)
	})
}

func testVars(f *ast.File, t *testing.T) {

	vars := f.Vars()

	assertString(t, "java_builder", vars["proj_name"])
	assertString(t, "https://github.com/ghazninattarshah/hybrid.test.runner.git", vars["proj_url"])
}

func testImage(f *ast.File, t *testing.T) {

	imap := f.Image()
	assertString(t, "elasticshift/java:1.9", imap[keys.NAME].(string))

	assertString(t, "http://dockerregistry.com/elasticshift", imap["registry"].(string))
	assertString(t, "testuser", imap["username"].(string))
	assertString(t, "isdf1i41i23iu12i", imap["token"].(string))
}

func testImageNames(f *ast.File, t *testing.T) {

	imgn := f.ImageNames()

	assertString(t, "elasticshift/java:1.9", imgn[0])
}

func testBlock(f *ast.File, t *testing.T) {

	type Expected struct {
		Key   string
		Value interface{}
	}

	items := map[int][]Expected{
		1: []Expected{
			{"name", "elasticshift/vcs"},
			{"description", "Checking out the project"},
			{"checkout", "https://github.com/ghazninattarshah/hybrid.test.runner.git"},
		},
		2: []Expected{
			{"name", "elasticshift/shell"},
			{"description", "Running maven compilation"},
			{"command", []string{"mvn clean build"}},
		},
		3: []Expected{
			{"name", "elasticshift/slack-notifier"},
			{"description", "Send notification to slack channel"},
			{"url", "https://hooks.slack.com/services/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"},
			{"channel", "#slack-notification"},
			{"username", "shiftbot"},
			{"icon_emoji", ":ghost:"},
			{"hint", map[string]string{"PARALLEL": "notification"}},
		},
		4: []Expected{
			{"name", "elasticshift/sendgrid"},
			{"description", "send email via sendgrid"},
			{"to", "ghazni.nattarshah@gmail.com"},
			{"cc", []string{"shahm.nattarshah@gmail.com", "shahbros@conspico.com"}},
			{"hint", map[string]string{"PARALLEL": "notification"}},
		},
		5: []Expected{
			{"name", "elasticshift/archive-sftp"},
			{"description", "Store the build archive to sftp"},
			{"hint", map[string]string{"PARALLEL": "archive"}},
		},
		6: []Expected{
			{"name", "elasticshift/archive-s3"},
			{"description", "Store the build archive to amazon s3"},
			{"hint", map[string]string{"PARALLEL": "archive"}},
		},
	}

	count := 1
	for f.HasMoreBlocks() {

		blk := f.NextBlock()

		num := blk[keys.BLOCK_NUMBER].(int)
		expected := items[num]

		for _, item := range expected {

			assertInt(t, count, num)

			switch item.Value.(type) {
			case string:
				assertString(t, item.Value.(string), blk[item.Key].(string))
			case map[string]string, []string:
				assertEqual(t, item.Value, blk[item.Key])
			}
		}
		count++
	}
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func assertString(t *testing.T, expected string, actual string) {

	if !strings.EqualFold(expected, actual) {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

func assertInt(t *testing.T, expected int, actual int) {

	if expected != actual {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
