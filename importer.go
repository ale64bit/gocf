package main

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func isPropNode(n *html.Node, mark string) (bool, string) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" {
				if a.Val == mark {
					var buffer bytes.Buffer
					for u := n.FirstChild; u != nil; u = u.NextSibling {
						if u.Type == html.TextNode {
							buffer.WriteString(u.Data)
						} else if u.Type == html.ElementNode && u.Data == "br" {
							buffer.WriteString("\n")
						}
					}
					return true, buffer.String()
				}
				break
			}
		}
	}
	return false, ""
}

func isMarkNode(n *html.Node, mark string) (bool, string) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == mark {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "pre" {
						var buffer bytes.Buffer
						for u := c.FirstChild; u != nil; u = u.NextSibling {
							if u.Type == html.TextNode {
								buffer.WriteString(u.Data)
							} else if u.Type == html.ElementNode && u.Data == "br" {
								buffer.WriteString("\n")
							}
						}
						return true, buffer.String()
					}
				}
				return false, ""
			}
		}
	}
	return false, ""
}

func getMemMultiplier(unit string) (mult int, err error) {
	switch unit {
	case "bytes":
		mult = 1
	case "kilobytes":
		mult = 1 << 10
	case "megabytes":
		mult = 1 << 20
	case "terabytes":
		mult = 1 << 30
	default:
		err = errors.New("unrecognized multiplier: " + unit)
	}
	return
}

func parseMemLimit(s string) int {
	parts := strings.Split(s, " ")
	if len(parts) == 2 {
		sz, err := strconv.Atoi(parts[0])
		if err != nil {
			return DefaultMemLimit
		}
		mult, err := getMemMultiplier(parts[1])
		if err != nil {
			return DefaultMemLimit
		}
		return sz * mult
	}
	return DefaultMemLimit

}

func getTimeMultiplier(unit string) (mult int, err error) {
	switch unit {
	case "second":
		mult = 1000
	case "seconds":
		mult = 1000
	case "millis":
		mult = 1
	case "milliseconds":
		mult = 1
	default:
		err = errors.New("unrecognized multiplier: " + unit)
	}
	return
}

func parseTimeLimit(s string) int {
	parts := strings.Split(s, " ")
	if len(parts) == 2 {
		sz, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return DefaultTimeLimit
		}
		mult, err := getTimeMultiplier(parts[1])
		if err != nil {
			return DefaultTimeLimit
		}
		return int(sz * float64(mult))
	}
	return DefaultTimeLimit
}

func parseInputFile(s string) string {
	switch s {
	case "standard input":
		return "*"
	default:
		return s
	}
}

func parseAnswerFile(s string) string {
	switch s {
	case "standard output":
		return "*"
	default:
		return s
	}
}

func parseContestAndTask(s string) (string, string) {
	u, _ := url.Parse(s)
	parts := strings.Split(u.Path[1:], "/")

	contestParts := []string{"codeforces"}
	for _, s = range parts[:len(parts)-1] {
		if len(s) > 0 && s != "problem" && s != "problemset" {
			contestParts = append(contestParts, s)
		}
	}
	return strings.Join(contestParts, "/"), parts[len(parts)-1]
}

func ImportCF(url string) (session GocfSession, inputs, answers []string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return
	}
	inputs = make([]string, 0)
	answers = make([]string, 0)
	session = DefaultSession()
	contest, task := parseContestAndTask(url)
	session.Contest = contest
	session.Task = task
	var f func(*html.Node)
	f = func(n *html.Node) {
		isML, ml := isPropNode(n, "memory-limit")
		isTL, tl := isPropNode(n, "time-limit")
		isInputFile, inputFile := isPropNode(n, "input-file")
		isAnswerFile, answerFile := isPropNode(n, "output-file")
		isInput, inputData := isMarkNode(n, "input")
		isOutput, answerData := isMarkNode(n, "output")
		switch {
		case isML:
			session.MemLimit = parseMemLimit(ml)
		case isTL:
			session.TimeLimit = parseTimeLimit(tl)
		case isInputFile:
			session.Input = parseInputFile(inputFile)
		case isAnswerFile:
			session.Output = parseAnswerFile(answerFile)
		case isInput:
			inputs = append(inputs, inputData)
		case isOutput:
			answers = append(answers, answerData)
		default:
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)
	return
}
