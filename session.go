package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type GocfSession struct {
	Contest   string
	Task      string
	Input     string
	Output    string
	TimeLimit int // milliseconds
	MemLimit  int // bytes
	Checker   string
}

const SessionFileName string = "/session.json"
const DefaultContest string = "practice"
const DefaultTask string = "task"
const DefaultInput string = "*"
const DefaultOutput string = "*"
const DefaultTimeLimit string = "1000"
const DefaultMemLimit string = "67108864"
const DefaultChecker string = "*"
const SolutionFile string = "__solution__"

var DefaultSession GocfSession = GocfSession{
	Contest:   "practice",
	Task:      "task",
	Input:     "*", // stdin
	Output:    "*", // stdout
	TimeLimit: 1000,
	MemLimit:  64 * (1 << 20),
	Checker:   "*", // default checker
}

func (session GocfSession) String() string {
	return "Session description:\n" +
		"  Contest:    " + session.Contest + "\n" +
		"  Task:       " + session.Task + "\n" +
		"  Input:      " + session.Input + "\n" +
		"  Output:     " + session.Output + "\n" +
		"  Time limit: " + strconv.Itoa(session.TimeLimit) + " [ms]\n" +
		"  Mem limit:  " + strconv.Itoa(session.MemLimit/(1<<20)) + " [MiB]\n" +
		"  Checker:    " + session.Checker + "\n"
}

func (session GocfSession) Save(config GocfConfig) {
	sessionFile := config.SessionDir + SessionFileName
	b, _ := json.Marshal(session)
	ioutil.WriteFile(sessionFile, b, os.ModePerm)
}

func (session GocfSession) archivePath(config GocfConfig) string {
	return config.ArchiveDir + "/" + session.Contest + "/" + session.Task
}

func (session GocfSession) NotArchived(config GocfConfig) bool {
	return FileNotExist(session.archivePath(config))
}

func (session GocfSession) Archive(config GocfConfig, overwrite bool) {
	testDir := config.SessionDir + "/" + TestPoolDir
	if FileExists(testDir) {
		os.RemoveAll(testDir)
	}
	if session.NotArchived(config) || overwrite {
		d := session.archivePath(config)
		if FileExists(d) {
			os.RemoveAll(d)
		}
		os.MkdirAll(d, os.ModePerm)
		filepath.Walk(config.SessionDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() { // ignore directories
				return nil
			}
			CopyFile(path, d+"/"+info.Name())
			return nil
		})
		CopyFile(config.WorkFile, d+"/"+SolutionFile)
	} else {
		panic("Archive path already exists and overwrite is false: " + session.archivePath(config))
	}
}

func LoadCurrentSession(config GocfConfig) GocfSession {
	sessionFile := config.SessionDir + SessionFileName
	var session GocfSession
	if FileNotExist(sessionFile) {
		session = DefaultSession
		session.Save(config)
	} else {
		contents, _ := ioutil.ReadFile(sessionFile)
		json.Unmarshal(contents, &session)
	}
	return session
}

func RestoreSession(config GocfConfig, contest, task string) {
	session := LoadCurrentSession(config)
	if session.NotArchived(config) {
		if Yes("Current session is not archived. Do you want to archive it?") {
			session.Archive(config, true)
		}
	}

	archiveDir := config.ArchiveDir + "/" + contest + "/" + task
	if FileNotExist(archiveDir) {
		fmt.Println("There is no archived session for contest " + contest + " and task " + task)
		return
	}
	os.RemoveAll(config.SessionDir)
	os.Mkdir(config.SessionDir, os.ModePerm)
	filepath.Walk(archiveDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Name() == SolutionFile { // ignore directories and solution file
			return nil
		}
		CopyFile(path, config.SessionDir+"/"+info.Name())
		return nil
	})
	CopyFile(archiveDir+"/"+SolutionFile, config.WorkFile)
	fmt.Println("Session restored")
}
