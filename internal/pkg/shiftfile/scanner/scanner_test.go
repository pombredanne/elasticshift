/*
Copyright 2018 The Elasticshift Authors.
*/
package scanner

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/elasticshift/elasticshift/internal/pkg/shiftfile/token"
)

type Expected struct {
	Line int
	Text string
}

type tokenPair struct {
	input    string
	err      error
	expected []Expected
}

var tokens = map[string][]tokenPair{
	"cache": []tokenPair{
		{`CACHE {
			- ~/.m2
		}`, nil, []Expected{
			{Line: 1, Text: "CACHE"},
			{Line: 1, Text: "{"},
			{Line: 2, Text: "~/.m2"},
			{Line: 3, Text: "}"},
		}},
	},
	"comments": []tokenPair{
		{"#", nil, []Expected{
			{Line: 1, Text: "#"},
		}},
		{"# comment", nil, []Expected{
			{Line: 1, Text: "# comment"},
		}},
		{"###### 12345 /* this is test comment */", nil, []Expected{
			{Line: 1, Text: "###### 12345 /* this is test comment */"},
		}},
		{"# All the parallel process represented by hint will run in separate thread", nil, []Expected{
			{Line: 1, Text: "# All the parallel process represented by hint will run in separate thread"},
		}},
	},
	"hints": []tokenPair{
		{"// PARALLEL:testinghint", nil, []Expected{
			{Line: 1, Text: "//"},
			{Line: 1, Text: "PARALLEL"},
			{Line: 1, Text: ":"},
			{Line: 1, Text: "testinghint"},
		}},
		{"/* PARALLEL:testinghint */", nil, []Expected{
			{Line: 1, Text: "/*"},
			{Line: 1, Text: "PARALLEL"},
			{Line: 1, Text: ":"},
			{Line: 1, Text: "testinghint"},
			{Line: 1, Text: "*/"},
		}},
		{`/*
	PARALLEL:testing
	PARALLEL:testing2
	*/`, nil, []Expected{
			{Line: 1, Text: "/*"},
			{Line: 2, Text: "PARALLEL"},
			{Line: 2, Text: ":"},
			{Line: 2, Text: "testing"},
			{Line: 3, Text: "PARALLEL"},
			{Line: 3, Text: ":"},
			{Line: 3, Text: "testing2"},
			{Line: 4, Text: "*/"},
		}},
	},
	"bool": []tokenPair{
		{"true", nil, []Expected{
			{Line: 1, Text: "true"},
		}},
		{"false", nil, []Expected{
			{Line: 1, Text: "false"},
		}},
	},
	"version": []tokenPair{
		{`VERSION "1.0"`, nil, []Expected{
			{Line: 1, Text: "VERSION"},
			{Line: 1, Text: "1.0"},
		}},
	},
	"arg": []tokenPair{
		{`@token`, nil, []Expected{
			{Line: 1, Text: "token"},
		}},
		{`@channel_name`, nil, []Expected{
			{Line: 1, Text: "channel_name"},
		}},
		{`"elasticshift/vcs", "Checking out the project" {
		checkout @checkout_url
		author (author)
		branch "develop"
		accesstoken ^1234kj12h34hh234
	}`, nil, []Expected{
			{Line: 1, Text: "elasticshift/vcs"},
			{Line: 1, Text: ","},
			{Line: 1, Text: "Checking out the project"},
			{Line: 1, Text: "{"},
			{Line: 2, Text: "checkout"},
			{Line: 2, Text: "checkout_url"},
			{Line: 3, Text: "author"},
			{Line: 3, Text: "("},
			{Line: 3, Text: "author"},
			{Line: 3, Text: ")"},
			{Line: 4, Text: "branch"},
			{Line: 4, Text: "develop"},
			{Line: 5, Text: "accesstoken"},
			{Line: 5, Text: "1234kj12h34hh234"},
			{Line: 6, Text: "}"},
		}},
	},
	"secret": []tokenPair{
		{`channel_token ^devtoken`, nil, []Expected{
			{Line: 1, Text: "channel_token"},
			{Line: 1, Text: "devtoken"},
		}},
	},
	"language": []tokenPair{
		{`LANGUAGE "java"`, nil, []Expected{
			{Line: 1, Text: "LANGUAGE"},
			{Line: 1, Text: "java"},
		}},
		{`LANGUAGE java`, nil, []Expected{
			{Line: 1, Text: "LANGUAGE"},
			{Line: 1, Text: "java"},
		}},
	},
	"image": []tokenPair{
		{`IMAGE "elasticshift/java:1.9" {
		registry "http://dockerregistry.com/elasticshift"
		SCRIPT {
			- apt-get install maven
		}
	}`, nil, []Expected{
			{Line: 1, Text: "IMAGE"},
			{Line: 1, Text: "elasticshift/java:1.9"},
			{Line: 1, Text: "{"},
			{Line: 2, Text: "registry"},
			{Line: 2, Text: "http://dockerregistry.com/elasticshift"},
			{Line: 3, Text: "SCRIPT"},
			{Line: 3, Text: "{"},
			{Line: 4, Text: "apt-get install maven"},
			{Line: 5, Text: "}"},
			{Line: 6, Text: "}"},
		}},
		{`IMAGE "elasticshift/java:1.9" {
		registry "http://dockerregistry.com/elasticshift"
		SCRIPT {
			"apt-get install maven"
		}
	}`, nil, []Expected{
			{Line: 1, Text: "IMAGE"},
			{Line: 1, Text: "elasticshift/java:1.9"},
			{Line: 1, Text: "{"},
			{Line: 2, Text: "registry"},
			{Line: 2, Text: "http://dockerregistry.com/elasticshift"},
			{Line: 3, Text: "SCRIPT"},
			{Line: 3, Text: "{"},
			{Line: 4, Text: "apt-get install maven"},
			{Line: 5, Text: "}"},
			{Line: 6, Text: "}"},
		}},
		{`IMAGE "elasticshift/java:1.9" {
		registry "http://dockerregistry.com/elasticshift"
		SCRIPT {
			"apt-get install maven && apt-get update"
		}
	}`, nil, []Expected{
			{Line: 1, Text: "IMAGE"},
			{Line: 1, Text: "elasticshift/java:1.9"},
			{Line: 1, Text: "{"},
			{Line: 2, Text: "registry"},
			{Line: 2, Text: "http://dockerregistry.com/elasticshift"},
			{Line: 3, Text: "SCRIPT"},
			{Line: 3, Text: "{"},
			{Line: 4, Text: "apt-get install maven && apt-get update"},
			{Line: 5, Text: "}"},
			{Line: 6, Text: "}"},
		}},
		{`IMAGE "elasticshift/java:1.9" {
		registry "http://dockerregistry.com/elasticshift"
		SCRIPT {
			- apt-get install maven \
			&& apt-get update
		}
	}`, nil, []Expected{
			{Line: 1, Text: "IMAGE"},
			{Line: 1, Text: "elasticshift/java:1.9"},
			{Line: 1, Text: "{"},
			{Line: 2, Text: "registry"},
			{Line: 2, Text: "http://dockerregistry.com/elasticshift"},
			{Line: 3, Text: "SCRIPT"},
			{Line: 3, Text: "{"},
			{Line: 5, Text: "apt-get install maven \\\n\t\t\t&& apt-get update"},
			{Line: 6, Text: "}"},
			{Line: 7, Text: "}"},
		}},
		{`IMAGE "elasticshift/java:1.9" {
		registry "http://dockerregistry.com/elasticshift"
		SCRIPT {
			- cd /my_folder; rm *.jar;
		}
	}`, nil, []Expected{
			{Line: 1, Text: "IMAGE"},
			{Line: 1, Text: "elasticshift/java:1.9"},
			{Line: 1, Text: "{"},
			{Line: 2, Text: "registry"},
			{Line: 2, Text: "http://dockerregistry.com/elasticshift"},
			{Line: 3, Text: "SCRIPT"},
			{Line: 3, Text: "{"},
			{Line: 4, Text: "cd /my_folder; rm *.jar;"},
			{Line: 5, Text: "}"},
			{Line: 6, Text: "}"},
		}},
	},
	"var": []tokenPair{
		{`VAR proj_url "http://github.com/ghazninattarshah/hybrid.test.runner.git"`, nil, []Expected{
			{Line: 1, Text: "VAR"},
			{Line: 1, Text: "proj_url"},
			{Line: 1, Text: "http://github.com/ghazninattarshah/hybrid.test.runner.git"},
		}},
	},
	"list": []tokenPair{
		{`cc ["ghazni.nattarshah@conspico.com", "shahm.nattarshah@conspico.com"]`, nil, []Expected{
			{Line: 1, Text: "cc"},
			{Line: 1, Text: "["},
			{Line: 1, Text: "ghazni.nattarshah@conspico.com"},
			{Line: 1, Text: ","},
			{Line: 1, Text: "shahm.nattarshah@conspico.com"},
			{Line: 1, Text: "]"},
		}},
	},
	"wordir": []tokenPair{
		{`WORKDIR "~/code"`, nil, []Expected{
			{Line: 1, Text: "WORKDIR"},
			{Line: 1, Text: "~/code"},
		}},
		{`WORKDIR ~/code`, nil, []Expected{
			{Line: 1, Text: "WORKDIR"},
			{Line: 1, Text: "~/code"},
		}},
	},
	"name": []tokenPair{
		{`NAME "elasticshift/java-maven-builder"`, nil, []Expected{
			{Line: 1, Text: "NAME"},
			{Line: 1, Text: "elasticshift/java-maven-builder"},
		}},
		{`NAME "elasticshift/ java-maven-builder"`, nil, []Expected{
			{Line: 1, Text: "NAME"},
			{Line: 1, Text: "elasticshift/ java-maven-builder"},
		}},
		{`NAME "elasticshift\n\rjava-maven-builder"`, nil, []Expected{
			{Line: 1, Text: "NAME"},
			{Line: 1, Text: "elasticshift\\n\\rjava-maven-builder"},
		}},
	},
	"realfile": []tokenPair{
		{
			`VERSION "1.0"

			# company/config-name
			# You can use..
			# FROM "elasticshift/java19-maven-builder"
			# in order to utilze the script with different build
			# The idea is to create script free builds,
			# in another word, you shall use someone else build script
			# to build your projects
			#
			NAME "elasticshift/java19-maven-builder"

			# Denotes the source code language
			LANGUAGE java

			# location of the source code and where the command starts from
			WORKDIR "~/code"

			# Variables are used when the parameter shall be changed
			# such as if you're invoking from hierarical build file.
			# this is really useful when build instructions are re-used.
			VAR proj_url "https://github.com/ghazninattarshah/hybrid.test.runner.git"

			# The container where the build is going to happen
			IMAGE "elasticshift/java:1.9" {
				registry "http://dockerregistry.com/elasticshift"
				script {
					"apt-get install maven"
					- apt-get install maven
				}
			}

			CACHE {
				- ~/.m2
				- ~/.gradle
				- ~/node-modules
			}

			#
			# Name of the plugin, description (this can be optional)
			# elasticshift - Name of the company who created this plugin
			# vcs - Name of the plugin
			#
			"elasticshift/vcs", "Checking out the project" {

				# Variable shall be used by enclosing with (..)
				checkout (proj_url)
			}

			"elasticshift/shell", "Running maven compilation" {
				- mvn clean build
			}

			"elasticshift/slack-notifier" ,"Send notification to slack channel" {
				# hint is the additional metadata to the plugin
				# which cause any process that have hint name notification will run on same group
				// PARALLEL:notification
				url "https://hooks.slack.com/services/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"
				channel "#slack-notification"
				username "shiftbot"
				icon_emoji ":ghost:"
			}

			"elasticshift/sendgrid", "send email via sendgrid" {
				# hint is the additional metadata to the plugin
				# which cause any process that have hint name notification will run on same group
				// PARALLEL:notification
				to "ghazni.nattarshah@gmail.com"
				cc ["ghazni.nattarshah@conspico.com", "shahm.nattarshah@conspico.com"]
			}

			"elasticshift/archive-sftp", "Store the build archive to sftp" {
				// PARALLEL:archive
			}

			"elasticshift/archive-s3", "Store the build archives to amazon s3" {
				// PARALLEL:archive
			}

			# All the parallel process represented by hint will run in separate thread
			# with in the group separated by hint identifier
			# Also, each process will utilize the multicore processor in order to run faster
			`, nil, []Expected{
				{Line: 1, Text: "VERSION"},
				{Line: 1, Text: "1.0"},
			},
		},
	},
	"shellfile": []tokenPair{
		{`VERSION "1.0"

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

		"shell", "checking out the project" {
			- git clone https://github.com/nshahm/hybrid.test.runner.git
		}

		"shell", "Building the project" {
			- ./gradlew clean build
		}

		`, nil, []Expected{
			{Line: 1, Text: "VERSION"},
			{Line: 1, Text: "1.0"},
		},
		}},
}

// "string": []tokenPair{
// 	input(token.STRING, " ", nil),
// },
// "number": []tokenPair{
// 	input(token.INT, "1".nil),
// 	input(token.INT, "123123123", nil),
// },
// "float": []tokenPair{
// 	input(token.FLOAT, "1.23", nil),
// },

var tokenTypes = []string{
	"cache",
	"comments",
	"hints",
	"bool",
	"string",
	"number",
	"float",
}

func testRealfile(t *testing.T) {
	toks := alltokens(tokens["realfile"][0].input)

	i := 1
	for _, tok := range toks {
		if i == tok.Position.Line {
			fmt.Print(fmt.Sprintf("%q\t", tok))
		} else {
			fmt.Print("\n")
			fmt.Print(fmt.Sprintf("%q\t", tok))
			i++
		}
	}

	//testTokenTypes(t, token.IMAGE, tokens["image"])
}

func TestShellfile(t *testing.T) {
	toks := alltokens(tokens["shellfile"][0].input)

	i := 1
	for _, tok := range toks {
		if i == tok.Position.Line {
			fmt.Print(fmt.Sprintf("%q\t", tok))
		} else {
			fmt.Print("\n")
			fmt.Print(fmt.Sprintf("%q\t", tok))
			i++
		}
	}

	//testTokenTypes(t, token.IMAGE, tokens["image"])
}

func TestCache(t *testing.T) {
	testTokenTypes(t, token.CACHE, tokens["cache"])
}

func TestImage(t *testing.T) {
	testTokenTypes(t, token.IMAGE, tokens["image"])
}

func TestVar(t *testing.T) {
	testTokenTypes(t, token.VAR, tokens["var"])
}

func TestWorkDir(t *testing.T) {
	testTokenTypes(t, token.WORKDIR, tokens["workdir"])
}

func TestLanguage(t *testing.T) {
	testTokenTypes(t, token.LANGUAGE, tokens["language"])
}

func TestName(t *testing.T) {
	testTokenTypes(t, token.NAME, tokens["name"])
}

func TestVersion(t *testing.T) {
	testTokenTypes(t, token.VERSION, tokens["version"])
}

func TestArg(t *testing.T) {
	testTokenTypes(t, token.ARGUMENT, tokens["arg"])
}

func TestSecret(t *testing.T) {
	testTokenTypes(t, token.SECRET, tokens["secret"])
}

func TestBool(t *testing.T) {
	testTokenTypes(t, token.BOOL, tokens["bool"])
}

func TestComments(t *testing.T) {
	testTokenTypes(t, token.COMMENT, tokens["comments"])
}

func TestHints(t *testing.T) {
	testTokenTypes(t, token.HINT, tokens["hints"])
}

func TestList(t *testing.T) {
	testTokenTypes(t, token.LBRACK, tokens["list"])
}

func alltokens(input string) []token.Token {

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s", input)

	errFunc := func(pos token.Position, msg string) {
		fmt.Println(fmt.Sprintf("position:%v, msg: %s", pos, msg))
	}
	s := New([]byte(buf.Bytes()), errFunc)

	outTokens := []token.Token{}
	var tok token.Token
	for s.HasMoreTokens() {

		tok = s.Scan()
		// fmt.Println(fmt.Sprintf("Out Token: %q", tok))
		outTokens = append(outTokens, tok)
	}
	return outTokens
}

func testTokenTypes(t *testing.T, ttype token.Type, tokens []tokenPair) {

	// errFunc := func(pos token.Position, msg string) {
	// 	fmt.Println(fmt.Sprintf("position:%v, msg: %s", pos, msg))
	// }

	for idx, in := range tokens {

		name := strconv.Itoa(idx + 1)
		t.Run(name, func(t *testing.T) {

			outTokens := alltokens(in.input)

			inLen := len(in.expected)
			outLen := len(outTokens)

			if inLen != outLen {
				t.Errorf("Number of token expected does not match: expected %d, actual %d", inLen, outLen)
			}

			for i := 0; i < len(in.expected); i++ {

				expectedToken := in.expected[i]
				outToken := outTokens[i]
				// fmt.Println(outToken)

				if !(expectedToken.Text == outToken.Text) {
					t.Errorf("tok = %q want %q for %q\n", outToken, expectedToken.Text, outToken.Text)
				}

				if expectedToken.Line != outToken.Position.Line {
					t.Errorf("tok = %q want in line %d, but actual is %d", outToken, expectedToken.Line, outToken.Position.Line)
				}
			}
		})
	}
}
