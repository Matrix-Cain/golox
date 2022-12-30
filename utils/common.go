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

func InterfaceToFloat64(a interface{}) (float64, bool) {
	switch v := a.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	}

	return 0, false
}
