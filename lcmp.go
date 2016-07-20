package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// This is the default checker if none is specified
// @see https://github.com/MikeMirzayanov/testlib/blob/master/checkers/lcmp.cpp

func skipEmpty(s []string) []string {
	var ret []string
	for _, si := range s {
		if len(si) > 0 {
			ret = append(ret, si)
		}
	}
	return ret
}

func compareWords(s, t string) bool {
	vs := skipEmpty(strings.Split(s, " "))
	vt := skipEmpty(strings.Split(t, " "))
	if len(vs) != len(vt) {
		return false
	}
	for i := 0; i < len(vs); i++ {
		if vs[i] != vt[i] {
			return false
		}
	}
	return true
}

func lcmp(outputFile, answerFile string) error {
	fout, _ := os.Open(outputFile)
	defer fout.Close()
	ouf := bufio.NewScanner(fout)
	fans, _ := os.Open(answerFile)
	defer fans.Close()
	ans := bufio.NewScanner(fans)
	n := 0
	for {
		eof := !ans.Scan()
		j := ans.Text()
		if j == "" && eof {
			break
		}
		ouf.Scan()
		p := ouf.Text()
		if !compareWords(j, p) {
			return errors.New(fmt.Sprintf("line %d differ - expected: %s, found: %s", n, j, p))
		}
		n++
	}
	return nil
}
