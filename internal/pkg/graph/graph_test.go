/*
Copyright 2018 The Elasticshift Authors.
*/
package graph

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/parser"
)

var file = `
VERSION "1.0"

NAME "elasticshift/java18-gradle-builder"

LANGUAGE java

WORKDIR "~/code"

#comment
VAR proj_url "https://github.com/nshahm/hybrid.test.runner.git"

# The container where the build is going to happen
IMAGE "openjdk:7"

CACHE {
	- ~/.gradle
}

#
# Name of the plugin, description (this can be optional)
# elasticshift - Name of the company who created this plugin
# vcs - Name of the plugin
#
"elasticshift/vcs", "Checking out the project" {
	url (proj_url)
	branch "master"
	directory "~/code"
}

"shell", "the build" {
	- mvn clean
	- mvn compile
	- mvn test
}
`

var file2 = `
VERSION "1.0"
NAME "elasticshift/java19-maven-builder"
LANGUAGE java
WORKDIR "~/code"
VAR proj_name "java_builder"
VAR proj_url "https://github.com/ghazninattarshah/hybrid.test.runner.git"
IMAGE "elasticshift/java:1.9" {
	registry "http://dockerregistry.com/elasticshift"
		- apt-get install maven
}
"elasticshift/vcs", "Checking out the project" {
	checkout (proj_url)
}
"shell", "Running maven compilation" {
	- mvn clean build
}
"elasticshift/slack-notifier" ,"Send notification to slack channel" {
	// PARALLEL:notification
	url "https://hooks.slack.com/services/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"
		channel "#slack-notification"
		username "shiftbot"
		icon_emoji ":ghost:"
}

"elasticshift/sendgrid", "send email via sendgrid" {
	// PARALLEL:notification
	to "ghazni.nattarshah@gmail.com"
		cc ["shahm.nattarshah@gmail.com", "shahbros@conspico.com"]
}

"elasticshift/archive-sftp", "Store the build archive to sftp" {
	// PARALLEL:archive
}

"elasticshift/archive-s3", "Store the build archive to amazon s3" {
	// PARALLEL:archive
}
`
var file3 = `VERSION "1.0"

NAME "elasticshift/java18-gradle-builder"

LANGUAGE java

WORKDIR "~/code"

#comment
VAR proj_url "https://github.com/nshahm/hybrid.test.runner.git"

# The container where the build is going to happen
IMAGE "openjdk:7" 

#
# Name of the plugin, description (this can be optional)
# elasticshift - Name of the company who created this plugin
# vcs - Name of the plugin
#
#"elasticshift/vcs", "Checking out the project" {
#	url (proj_url)
#  branch "master"
#  directory "~/code"
#}

CACHE {
	- ~/.gradle
}

"shell", "checking out the project" {
	- git clone https://github.com/nshahm/hybrid.test.runner.git ~/code
}

"shell", "echo 1" {
	// PARALLEL:echogroup
	- echo "fan1"
	- sleep 5
}

"shell", "echo 2" {
	// PARALLEL:echogroup
	- echo "fan2"
	- sleep 5
}

"shell", "echo 3" {
	// PARALLEL:echogroup
	- echo "fan3"
	- sleep 5
}

"shell", "Building the project" {
	- ./gradlew clean build
}

`

func TestDuration(t *testing.T) {

	//t1 := time.Date(2016, time.August, 15, 0, 20, 15, 125, time.UTC)
	t1 := time.Date(2017, time.February, 16, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2017, time.February, 16, 0, 0, 0, 123267565, time.UTC)

	d := t2.Sub(t1)

	fmt.Println(d.String())

	// text := ""
	var text string

	// fmt.Printf("Duration: %dh %dm %ds %dns \n", int(t3.Hours()), int(t3.Minutes()), int(t3.Seconds()), int(t3.Nanoseconds()))

	h := int(d.Hours())
	m := int(d.Minutes()) - (h * 60)
	s := int(d.Seconds()) - (int(d.Minutes()) * 60)

	if h > 0 {
		text += fmt.Sprintf("%.02dh ", h)
	}

	if m > 0 {
		text += fmt.Sprintf("%.02dm ", m)
	}

	if s > 0 {
		text += fmt.Sprintf("%.02ds ", s)
	} else {

		durStr := d.String()

		var finalDur string
		var stripLen int
		stripLen = 3
		if strings.HasSuffix(durStr, "ms") || strings.HasSuffix(durStr, "ns") {
			stripLen = 2
		}

		idx := len(durStr) - stripLen
		dur := durStr[:idx]
		notation := durStr[idx:]

		dotIdx := strings.Index(dur, ".")
		if dotIdx > 0 {

			befrDec := dur[:dotIdx+1]
			aftrDec := dur[dotIdx+1:]
			if len(aftrDec) > 3 {
				finalDur = befrDec + aftrDec[:3] + notation
			} else {
				finalDur = befrDec + aftrDec + notation
			}
		} else {
			finalDur = d.String()
		}

		text += finalDur
	}
}

func TestNodeLevel(t *testing.T) {

	f, err := parser.AST([]byte(file))
	if err != nil {
		t.Fail()
	}

	graph, err := Construct(f)
	fmt.Println(graph.String())

	f, err = parser.AST([]byte(file2))
	if err != nil {
		t.Fail()
	}

	graph, err = Construct(f)
	fmt.Println(graph.String())

	f, err = parser.AST([]byte(file3))
	if err != nil {
		t.Fail()
	}

	graph, err = Construct(f)
	fmt.Println(graph.String())
}

func TestGraph(t *testing.T) {

	f, err := parser.AST([]byte(file))
	if err != nil {
		t.Fail()
	}

	graph, err := Construct(f)
	// assertString(t, `(1) START
	// (2) elasticshift/vcs
	// (3) elasticshift/shell
	// (4) END
	// `, graph.String())

	f, err = parser.AST([]byte(file2))
	if err != nil {
		t.Fail()
	}

	graph, err = Construct(f)

	fmt.Println(graph.JSON())
	assertString(t, `(1) START
(2) elasticshift/vcs
(3) elasticshift/shell
(4) FANOUT-notification
(4) - elasticshift/slack-notifier
(4) - elasticshift/sendgrid
(5) FANIN-notification
(6) FANOUT-archive
(6) - elasticshift/archive-sftp
(6) - elasticshift/archive-s3
(7) FANIN-archive
(8) END
`, graph.String())

	fmt.Println(graph.JSON())
}

func assertString(t *testing.T, expected string, actual string) {

	if !strings.EqualFold(expected, actual) {
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}
