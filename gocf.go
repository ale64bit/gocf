package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
)

func CreateSession(config GocfConfig) {
	session := LoadCurrentSession(config)
	if Yes("Do you want to archive current session?") {
		session.Archive(config, false)
	}

	os.RemoveAll(config.SessionDir)
	os.MkdirAll(config.SessionDir, os.ModePerm)

	contest := ReadDefault("Enter contest name", DefaultContest)
	task := ReadDefault("Enter task name", DefaultTask)
	input := ReadDefault("Enter input file name", DefaultInput)
	output := ReadDefault("Enter output file name", DefaultOutput)
	tl := ReadDefault("Enter time limit", strconv.Itoa(DefaultTimeLimit))
	ml := ReadDefault("Enter memory limit", strconv.Itoa(DefaultMemLimit))
	checker := ReadDefault("Enter task checker", DefaultChecker)

	timeLimit, _ := strconv.Atoi(tl)
	memLimit, _ := strconv.Atoi(ml)
	session = GocfSession{contest, task, input, output, timeLimit, memLimit, checker}
	session.Save(config)

	ioutil.WriteFile(config.WorkFile, []byte(GoTemplate(session)), os.ModePerm)
	fmt.Println("done")
}

func ImportSession(config GocfConfig, s string) {
	session := LoadCurrentSession(config)
	if Yes("Do you want to archive current session?") {
		session.Archive(config, false)
	}

	os.RemoveAll(config.SessionDir)
	os.MkdirAll(config.SessionDir, os.ModePerm)

	u, _ := url.Parse(s)
	switch u.Host {
	case "codeforces.com":
		session, inputs, answers, err := ImportCF(s)
		if err != nil {
			panic(err)
		}
		session.Save(config)
		ioutil.WriteFile(config.WorkFile, []byte(GoTemplate(session)), os.ModePerm)

		nrOfTests := len(inputs)
		for id := 1; id <= nrOfTests; id++ {
			ioutil.WriteFile(inPath(config, id), []byte(inputs[id-1]), os.ModePerm)
			ioutil.WriteFile(ansPath(config, id), []byte(answers[id-1]), os.ModePerm)
		}
		fmt.Println("import successful")

	default:
		fmt.Println("Unsupported source:", u.Host)
	}
}

func ArchiveSession(config GocfConfig) {
	session := LoadCurrentSession(config)
	if session.NotArchived(config) {
		session.Archive(config, false)
	} else {
		if Yes("This session is already archived. Do you want to overwrite it?") {
			session.Archive(config, true)
		}
	}
}

func ListSession(config GocfConfig) {
	session := LoadCurrentSession(config)
	fmt.Println(session.String())
	ListTests(config, session)
}

func PrintUsage() {
	fmt.Println(`
-----------------------------
   ____        ____ _____
  / ___| ___  / ___|  ___|
 | |  _ / _ \| |   | |_
 | |_| | (_) | |___|  _|
  \____|\___/ \____|_|

-----------------------------
Usage: gocf <cmd> [args...]

where <cmd> is one of:
  create                   - create a new session
  import <url>             - create a new session from a supported url (e.g. Codeforces, Timus)
  test                     - compile and run work file againts current tests
  add                      - add a new test to current session
  rm <id>                  - remove the test #id from current session
  archive                  - archive current session
  restore <contest> <task> - restore an archived session 
  ls                       - list current session properties and tests
`)
}

func CheckArgCount(exp int) {
	if len(os.Args)-2 != exp {
		PrintUsage()
		os.Exit(1)
	}
}

func main() {

	if len(os.Args) == 1 {
		PrintUsage()
		os.Exit(0)
	}

	config := LoadConfig()
	cmd := os.Args[1]
	switch cmd {
	case "create":
		CheckArgCount(0)
		CreateSession(config)
	case "import":
		CheckArgCount(1)
		ImportSession(config, os.Args[2])
	case "test":
		TestAll(config)
	case "add":
		CheckArgCount(0)
		AddTestFromUser(config)
	case "rm":
		CheckArgCount(1)
		id, _ := strconv.Atoi(os.Args[2])
		RemoveTest(config, id)
	case "archive":
		CheckArgCount(0)
		ArchiveSession(config)
	case "restore":
		CheckArgCount(2)
		RestoreSession(config, os.Args[2], os.Args[3])
	case "ls":
		ListSession(config)
	default:
		PrintUsage()
	}
}
