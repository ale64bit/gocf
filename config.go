package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const ConfigFile string = "$HOME/.gocf"
const DefaultWorkFile string = "$GOPATH/src/a.go"
const DefaultSessionDir string = "$HOME/GocfSession"
const DefaultArchiveDir string = "$HOME/GocfArchive"

type GocfConfig struct {
	WorkFile   string
	SessionDir string
	ArchiveDir string
}

func readWorkFile() string {
	var workFile string
	fmt.Print("Enter work file [default=" + DefaultWorkFile + "]: ")
	fmt.Scanf("%s\n", &workFile)
	if len(workFile) == 0 {
		workFile = DefaultWorkFile
	}
	workFile = os.ExpandEnv(workFile)
	if FileNotExist(workFile) {
		fmt.Println("Work file doesn't exist. Creating...")
		os.Create(workFile)
	}
	return workFile
}

func readSessionDir() string {
	var sessionDir string
	fmt.Print("Enter session directory path [default=" + DefaultSessionDir + "]: ")
	fmt.Scanf("%s\n", &sessionDir)
	if len(sessionDir) == 0 {
		sessionDir = DefaultSessionDir
	}
	sessionDir = os.ExpandEnv(sessionDir)
	if FileNotExist(sessionDir) {
		fmt.Println("Session directory doesn't exist. Creating...")
		os.MkdirAll(sessionDir, os.ModePerm)
	}
	return sessionDir
}

func readArchiveDir() string {
	var archiveDir string
	fmt.Print("Enter archive directory path [default=" + DefaultArchiveDir + "]: ")
	fmt.Scanf("%s\n", &archiveDir)
	if len(archiveDir) == 0 {
		archiveDir = DefaultArchiveDir
	}
	archiveDir = os.ExpandEnv(archiveDir)
	if FileNotExist(archiveDir) {
		fmt.Println("Archive directory doesn't exist. Creating...")
		os.MkdirAll(archiveDir, os.ModePerm)
	}
	return archiveDir
}

func SaveConfig(conf GocfConfig) {
	b, _ := json.Marshal(conf)
	ioutil.WriteFile(os.ExpandEnv(ConfigFile), b, os.ModePerm)
}

func LoadConfig() GocfConfig {
	configFile := os.ExpandEnv(ConfigFile)
	if FileNotExist(configFile) {
		fmt.Println("Configuration file not found.")
		workFile := readWorkFile()
		sessionDir := readSessionDir()
		archiveDir := readArchiveDir()
		conf := GocfConfig{workFile, sessionDir, archiveDir}
		SaveConfig(conf)
		return conf
	} else {
		fmt.Println("Configuration file found. Loading...")
		contents, _ := ioutil.ReadFile(configFile)
		var conf GocfConfig
		json.Unmarshal(contents, &conf)
		return conf
	}
}
