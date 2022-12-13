package VM

import (
	"bufio"
	"fmt"
	"golox/lox/lexer"
	"golox/utils"
	"os"
)

type VM struct {
	hadError bool
}

func (v *VM) RunFile(path string) {
	fileBytes, _ := os.ReadFile(path)
	v.run(string(fileBytes[:]))
	// Indicate an error in the exit code.
	if v.hadError {
		os.Exit(65)
	}
}

func (v *VM) RunStr(code string) {
	v.run(code)
	// Indicate an error in the exit code.
	if v.hadError {
		os.Exit(65)
	}
}

func (v *VM) RunPrompt() {
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
		v.run(line)
		v.hadError = false
	}
}

func (v *VM) run(source string) {
	scanner := lexer.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err.HasError {
		utils.RaiseError(err.Line, err.Reason)
		v.hadError = true
	}

	// For now, just print the tokens.
	for _, v := range tokens {
		fmt.Println(v)
		//log.Info(v)
	}
}

func (v *VM) SetError(error bool) {
	v.hadError = error
}
