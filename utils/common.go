package utils

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

func RaiseError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Errorln("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
}
