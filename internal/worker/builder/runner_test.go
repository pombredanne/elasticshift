/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/parser"
	"gitlab.com/conspico/elasticshift/pkg/worker/logger"
)

var testfileShellOnly = `
VERSION "1.0"
NAME "elasticshift/runner"

IMAGE "alphine:latest"

"elasticshift/shell", "Checking out the project" {
	- echo "git clone"
	- sleep 4
}

"elasticshift/shell", "Running maven compilation" {
	- echo "mvn clean build"
}

"elasticshift/shell" ,"Send notification to slack channel" {
	// PARALLEL:notification
	- echo "Notifying slack.."
	- sleep 3
}

"elasticshift/shell", "send email via sendgrid" {
	// PARALLEL:notification
	to "ghazni.nattarshah@gmail.com"
	cc ["shahm.nattarshah@gmail.com", "shahbros@conspico.com"]
	- echo "Notifying sendgrid"
	- sleep 8
}

"elasticshift/shell", "Store the build archive to sftp" {
	// PARALLEL:archive
	- echo "Archiving to sftp"
}

"elasticshift/shell", "Store the build archive to amazon s3" {
	// PARALLEL:archive
	- echo "Archiving to s3"
}
`

var testfile = `
VERSION "1.0"
NAME "elasticshift/runner"

IMAGE "alphine:latest"

"elasticshift/shell", "Checking out the project" {
	- echo "git clone"
	- sleep 4
}

"elasticshift/shell", "Running maven compilation" {
	- echo "mvn clean build"
	- sleep 5
}

"elasticshift/slack-notifier" ,"Send notification to slack channel" {
	// PARALLEL:notification
	- echo "Notifying slack.."
	- sleep 2
}

"elasticshift/sendgrid", "send email via sendgrid" {
	// PARALLEL:notification
	to "ghazni.nattarshah@gmail.com"
	cc ["shahm.nattarshah@gmail.com", "shahbros@conspico.com"]
	- echo "Notifying sendgrid"
}

"elasticshift/archive-sftp", "Store the build archive to sftp" {
	// PARALLEL:archive
	- echo "Archiving to sftp"
}

"elasticshift/archive-s3", "Store the build archive to amazon s3" {
	// PARALLEL:archive
	- echo "Archiving to s3"
}
`

func TestRunner(t *testing.T) {

	f, err := parser.AST([]byte(testfileShellOnly))
	if err != nil {
		t.Fail()
	}

	graph, err := ConstructGraph(f)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(graph.String())

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)
	b := &builder{f: f, logr: &logger.Logr{Writer: os.Stdout}}
	err = b.build(graph)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
