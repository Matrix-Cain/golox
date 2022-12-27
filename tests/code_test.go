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
print "带专";
`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}

func TestCodeSnippet3(t *testing.T) {
	snippet := `
"2"+3
1+1
2>1:4/3.1415926?3
`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}

func TestCodeSnippet4(t *testing.T) {
	snippet := `
true+true
`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}

func TestCodeSnippet5(t *testing.T) {
	snippet := `
var a = "global a";
var b = "global b";
var c = "global c";
{
  var a = "outer a";
  var b = "outer b";
  {
    var a = "inner a";
    print a;
    print b;
    print c;
  }
  print a;
  print b;
  print c;
}
print a;
print b;
print c;`
	vm := VM.VM{}
	vm.RunStr(snippet)
}

/* no `;` at end of line */
func TestMalformedCodeSnippet1(t *testing.T) {
	snippet := `// Your first Lox program!
/* nested
comment
*/
var a = "你好世界！!";
a = "帽子戏法";
print a;
a=2 // no ;
print a;`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}
