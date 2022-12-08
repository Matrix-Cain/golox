package tests

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestZH(t *testing.T) {
	in := "你1好1世1界"
	log.Infof(string(in[0]))
	log.Infof(string(in[1]))
	log.Infof(string(in[2]))
	log.Infof(string(in[3]))
	log.Infof(string(in[4]))
	log.Infof(string(in[5]))
	log.Infof(string(in[6]))
}
