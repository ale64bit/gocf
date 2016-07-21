package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	OK = iota
	WA
	TLE
	MLE
	RTE
)

const TestPoolDir string = "__pool__"

func ResultMsg(r int) string {
	switch r {
	case OK:
		return "OK"
	case WA:
		return "Wrong Answer"
	case TLE:
		return "Time Limit Exceeded"
	case MLE:
		return "Memory Limit Exceeded"
	case RTE:
		return "Runtime Error"
	default:
		panic("Unrecognized result code: " + strconv.Itoa(r))
	}
}

func firstAvailableId(config GocfConfig) (id int) {
	id = 1
	for FileExists(inPath(config, id)) {
		id++
	}
	return
}

func inPath(config GocfConfig, id int) string {
	return config.SessionDir + "/" + strconv.Itoa(id) + ".in"
}

func outPath(config GocfConfig, id int) string {
	return config.SessionDir + "/" + strconv.Itoa(id) + ".out"
}

func ansPath(config GocfConfig, id int) string {
	return config.SessionDir + "/" + strconv.Itoa(id) + ".ans"
}

func AddTest(config GocfConfig, input, answer []byte) int {
	id := firstAvailableId(config)
	ioutil.WriteFile(inPath(config, id), input, os.ModePerm)
	if len(answer) > 0 {
		ioutil.WriteFile(ansPath(config, id), answer, os.ModePerm)
	}
	return id
}

func AddTestFromUser(config GocfConfig) {
	session := LoadCurrentSession(config)
	fmt.Println(session.String())
	fmt.Println("\nEnter input:")
	input, _ := ioutil.ReadAll(os.Stdin)
	fmt.Println("\nEnter answer [empty if unknown]:")
	answer, _ := ioutil.ReadAll(os.Stdin)
	id := AddTest(config, input, answer)
	fmt.Println("Added test #", id)
}

func RemoveTest(config GocfConfig, id int) {
	if FileExists(inPath(config, id)) {
		os.Remove(inPath(config, id))
	} else {
		fmt.Println("No test with such id")
		return
	}
	if FileExists(ansPath(config, id)) {
		os.Remove(ansPath(config, id))
	}
	for FileExists(inPath(config, id+1)) {
		os.Rename(inPath(config, id+1), inPath(config, id))
		if FileExists(ansPath(config, id+1)) {
			os.Rename(ansPath(config, id+1), ansPath(config, id))
		}
		id++
	}
}

func CleanTestDir(config GocfConfig, session GocfSession) {
	testDir := config.SessionDir + "/" + TestPoolDir
	if FileExists(testDir) {
		err := os.RemoveAll(testDir)
		if err != nil {
			panic(err)
		}
	}
	os.MkdirAll(testDir, os.ModePerm)
}

func PopulateTestDir(config GocfConfig, session GocfSession) {
	testDir := config.SessionDir + "/" + TestPoolDir
	id := 1
	for FileExists(inPath(config, id)) {
		os.Symlink(inPath(config, id), testDir+"/"+strconv.Itoa(id)+".in")
		if FileExists(ansPath(config, id)) {
			os.Symlink(ansPath(config, id), testDir+"/"+strconv.Itoa(id)+".ans")
		}
		id++
	}
}

func Compile(config GocfConfig, session GocfSession) {
	// TODO support other languages?
	bin := config.SessionDir + "/" + TestPoolDir + "/solution"
	if FileExists(bin) {
		os.Remove(bin)
	}
	cmd := exec.Command("go", "build", "-o", bin, config.WorkFile)
	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Compilation error")
		fmt.Println(out.String())
		os.Exit(1)
	}
}

func run(cmd *exec.Cmd, session GocfSession) int {
	if err := cmd.Start(); err != nil {
		return RTE
	}
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			return RTE
		} else {
			return OK
		}
	case <-time.After(time.Duration(session.TimeLimit) * time.Millisecond):
		cmd.Process.Kill()
		return TLE
	}
}

func check(config GocfConfig, session GocfSession, id int) int {
	poolDir := config.SessionDir + "/" + TestPoolDir
	inputFile := poolDir + "/" + strconv.Itoa(id) + ".in"
	outputFile := poolDir + "/" + strconv.Itoa(id) + ".out"
	answerFile := poolDir + "/" + strconv.Itoa(id) + ".ans"
	if FileNotExist(answerFile) {
		return OK
	}

	if session.Checker == "*" {
		if err := lcmp(outputFile, answerFile); err != nil {
			return WA
		}
	} else {
		cmd := exec.Command(session.Checker, inputFile, outputFile, answerFile)
		cmd.Dir = poolDir
		if err := cmd.Run(); err != nil {
			return WA
		}
	}

	return OK
}

func TestOne(config GocfConfig, session GocfSession, id int) (int, time.Duration) {
	poolDir := config.SessionDir + "/" + TestPoolDir
	bin := poolDir + "/solution"
	cmd := exec.Command(bin)

	cmd.Dir = poolDir
	if session.Input == "*" {
		// redirect input file to process standard input
		r, _ := os.Open(poolDir + "/" + strconv.Itoa(id) + ".in")
		defer r.Close()
		cmd.Stdin = r
	} else {
		// create symlink to expected input file
		os.Link(poolDir+"/"+strconv.Itoa(id)+".in", poolDir+"/"+session.Input)
		defer os.Remove(poolDir + "/" + session.Input)
	}

	if session.Output == "*" {
		// redirect output
		out, _ := os.Create(poolDir + "/" + strconv.Itoa(id) + ".out")
		defer out.Close()
		cmd.Stdout = out
	}

	start := time.Now()
	if err := run(cmd, session); err != OK {
		elapsed := time.Since(start)
		return err, elapsed
	}
	elapsed := time.Since(start)

	if session.Output != "*" {
		os.Rename(poolDir+"/"+session.Output, poolDir+"/"+strconv.Itoa(id)+".out")
	}

	return check(config, session, id), elapsed
}

func TestAll(config GocfConfig) {
	session := LoadCurrentSession(config)
	fmt.Println("Removing test directory...")
	CleanTestDir(config, session)
	fmt.Println("Copying test files...")
	PopulateTestDir(config, session)
	fmt.Println("Compiling...")
	Compile(config, session)
	if session.Checker != "*" && FileNotExist(session.Checker) {
		fmt.Println("Checker not found:", session.Checker)
		os.Exit(1)
	}
	fmt.Println("Running...")
	var outcomes []int
	var times []time.Duration
	id := 1
	for FileExists(inPath(config, id)) {
		outcome, elapsed := TestOne(config, session, id)
		outcomes = append(outcomes, outcome)
		times = append(times, elapsed)
		id++
	}

	PrintResults(config, session, outcomes, times)
}

func PrintResults(config GocfConfig, session GocfSession, outcomes []int, times []time.Duration) {
	testDir := config.SessionDir + "/" + TestPoolDir
	testCount := len(outcomes)
	for id := 1; id <= testCount; id++ {
		inputFile := testDir + "/" + strconv.Itoa(id) + ".in"
		outputFile := testDir + "/" + strconv.Itoa(id) + ".out"
		answerFile := testDir + "/" + strconv.Itoa(id) + ".ans"
		fmt.Println("----------------------------------------------------------")
		fmt.Printf("Test #%d:\n", id)

		fmt.Println("Input:")
		input, _ := ioutil.ReadFile(inputFile)
		fmt.Println(string(input))

		fmt.Println("Expected output:")
		if FileExists(answerFile) {
			answer, _ := ioutil.ReadFile(answerFile)
			fmt.Println(string(answer))
		} else {
			fmt.Println("UNKNOWN\n")
		}

		fmt.Println("Execution output:")
		if FileExists(outputFile) {
			output, _ := ioutil.ReadFile(outputFile)
			fmt.Println(string(output))
		} else {
			fmt.Println("\n")
		}
	}

	fmt.Println("==========================================================")
	fmt.Println(" SUMMARY")
	fmt.Println("==========================================================")
	passed := 0
	for i := 0; i < testCount; i++ {
		if outcomes[i] == OK {
			passed++
		}
		fmt.Printf("  Test #%d [%.3fs]: %s\n", i+1, times[i].Seconds(), ResultMsg(outcomes[i]))
	}
	fmt.Println("----------------------------------------------------------")
	if passed == testCount {
		fmt.Println(" RESULT: All tests passed!")
	} else {
		fmt.Println(" RESULT: Some tests are failing...")
	}
	fmt.Println("==========================================================")
}

func ListTests(config GocfConfig, session GocfSession) {
	fmt.Println("TESTS")
	id := 1
	for FileExists(inPath(config, id)) {
		fmt.Println("-----------------------------------------------------------------------")
		fmt.Println("INPUT #" + strconv.Itoa(id))
		fmt.Println("-----------------------------------------------------------------------")
		b, _ := ioutil.ReadFile(inPath(config, id))
		fmt.Print(string(b))
		if FileExists(ansPath(config, id)) {
			fmt.Println("-----------------------------------------------------------------------")
			fmt.Println("ANSWER #" + strconv.Itoa(id))
			fmt.Println("-----------------------------------------------------------------------")
			b, _ := ioutil.ReadFile(ansPath(config, id))
			fmt.Print(string(b))
		}
		id++
	}
	fmt.Println("-----------------------------------------------------------------------")
}
