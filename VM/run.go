package VM

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golox/lox/ast"
	"golox/lox/lexer"
	"golox/lox/parser"
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
	tokens, lexerError := scanner.ScanTokens()
	if lexerError.HasError {
		utils.RaiseError(lexerError.Line, lexerError.Reason)
		v.hadError = true
	}

	parser0 := parser.NewParser(tokens)
	expression, parseError := parser0.Parse()

	if parseError.HasError {
		v.hadError = true
	}

	if v.hadError {
		return
	}
	printer := ast.AstPrinter{}
	str, _ := printer.Print(expression)

	log.Println(str.(string))

}

func (v *VM) SetError(error bool) {
	v.hadError = error
}
