package main

func FileTemplate(config GocfConfig, session GocfSession) string {
	return `package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	fi, _ := os.Open("` + session.Input + `")
	fo, _ := os.Create("` + session.Output + `")
	scanner = bufio.NewScanner(fi)
	scanner.Split(bufio.ScanWords)
	writer = bufio.NewWriter(fo)
	defer fi.Close()
	defer fo.Close()
	defer writer.Flush()

	// TODO your code here
}

/******************/
/* IO boilerplate */
/******************/

var scanner *bufio.Scanner
var writer *bufio.Writer

func NextInt() int {
	ret, _ := strconv.Atoi(Next())
	return ret
}

func Next() string {
	scanner.Scan()
	return scanner.Text()
}

func Print(s string) {
	writer.WriteString(s)
}

func Println(s string) {
	writer.WriteString(s)
	writer.WriteByte('\n')
}
`
}