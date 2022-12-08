package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var hadError = false

func main() {
	if len(os.Args) > 2 {
		log.Error("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}

}

func runFile(path string) {
	fileBytes, _ := os.ReadFile(path)
	run(string(fileBytes[:]))
	// Indicate an error in the exit code.
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	var line string
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("> ")
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			break
		} else {
			line = string(lineBytes[:])
			if line == "" {
				break
			}
		}
		run(line)
		hadError = false
	}
}
func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.scanTokens()

	// For now, just print the tokens.
	for _, v := range tokens {
		fmt.Println(v)
		//log.Info(v)
	}
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Errorln("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
	hadError = true
}
