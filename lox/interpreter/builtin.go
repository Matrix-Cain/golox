package interpreter

import "time"

type Clock struct {
}

func (t *Clock) Arity() int {
	return 0
}

func (t *Clock) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return time.Now().UnixMicro(), nil
}

func (t *Clock) String() string {
	return "<native fn>"
}
