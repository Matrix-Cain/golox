package tests

import (
	"golox/VM"
	"testing"
)

func TestCodeSnippet1(t *testing.T) {
	snippet := `// Your first Lox program!
/* nested
comment
*/
var a = "你好世界！!";
print a;`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}

func TestCodeSnippet2(t *testing.T) {
	snippet := `// Your first Lox program!
/* nested
comment
*/
(1+2)/3*6;
`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}

/* no `;` at end of line */
func TestMalformedCodeSnippet1(t *testing.T) {
	snippet := `// Your first Lox program!
/* nested
comment
*/
var a = "你好世界！!"
print a;`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}
