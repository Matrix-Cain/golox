package utils

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

func RaiseError(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	log.Errorln("[line " + strconv.Itoa(line) + "] Error " + where + ": " + message)
}
