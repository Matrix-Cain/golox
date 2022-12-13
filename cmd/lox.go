package main

import (
	log "github.com/sirupsen/logrus"
	"golox/VM"
	"os"
)

var vm *VM.VM

func init() {
	vm = &VM.VM{}
}

func main() {
	if len(os.Args) > 2 {
		log.Error("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		vm.RunFile(os.Args[1])
	} else {
		vm.RunPrompt()
	}

}
