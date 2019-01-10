/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"fmt"
	"testing"

	"github.com/elasticshift/shiftfile/parser"
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

CACHE {
	- ~/.gradle
}

"shell", "Checking out the project" {
	- echo "git clone"
	#- sleep 4
}

"shell", "Running maven compilation" {
	- echo "mvn clean build"
	#- sleep 5
}

"elasticshift/slack-notifier" ,"Send notification to slack channel" {
	// PARALLEL:notification
	- echo "Notifying slack.."
	#- sleep 2
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

var shellonly = `
VERSION "1.0"

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
	- git clone https://github.com/nshahm/hybrid.test.runner.git
}

"shell", "Building the project" {
	- ./gradlew clean build
}`

func TestRunner(t *testing.T) {

	f, err := parser.AST([]byte(shellonly))
	if err != nil {
		t.Fail()
	}

	fmt.Printf("%#v\n", f)
	for f.HasMoreBlocks() {
		fmt.Println(f.NextBlock())
	}

	// graph, err := ConstructGraph(f)
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fail()
	// }

	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)
	// b := &builder{f: f, logr: &logger.Logr{Writer: os.Stdout}, g: graph}
	// err = b.build(graph)
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fail()
	// }
	// fmt.Println(graph.Json())
}
