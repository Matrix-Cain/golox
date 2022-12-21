package utils

import (
	log "github.com/sirupsen/logrus"
	"golox/lox/lexer"
	"strconv"
)

func RaiseError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Errorln("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
}

func Error(token lexer.Token, message string) {
	if token.Type0 == lexer.EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}
