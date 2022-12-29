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
var a = 3;
a = a + 1;
print a;
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

func TestIf1(t *testing.T) {
	snippet := `
var a = 10;
if(a>0){
print a;
}else{
print "not possible";
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)

}

func TestElse1(t *testing.T) {
	snippet := `
var a = -3;
if(a>0){
print "not possible";
}else{
print a;
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)

}

func TestWhile1(t *testing.T) {
	snippet := `
var a = -100;
while(a<0){
print a;
a = a+1;
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)

}
func TestInfiniteWhile(t *testing.T) {
	snippet := `
var a = 0;
while(a==0){
print a;
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}

func TestFibonacciCodeSnippet(t *testing.T) {
	snippet := `
var a = 0;
var temp;

for (var b = 1; a < 10000; b = temp + b) {
  print a;
  temp = a;
  a = b;
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}

func TestForLoop(t *testing.T) {
	snippet := `
for (var a=0; a<10; a = a+1){
	print a;
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)

}

func TestForLoopBreak(t *testing.T) {
	snippet := `
for (var a=0; a<10;){
	print a;
	break;
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)

}

func TestForLoopContinue(t *testing.T) {
	snippet := `
for (var a=0; a<10; a++){
	print a;
	if(a < 5){
		continue;
	}else{
		break;
	}
}
`
	vm := &VM.VM{}
	vm.RunStr(snippet)

}

func TestBasicNoReturnFunc(t *testing.T) {
	snippet := `
fun sayHi(first, last) {
  print "Hi, " + first + " " + last + "!";
}

sayHi("Dear", "Reader");
`
	vm := &VM.VM{}
	vm.RunStr(snippet)

}

func TestBasicWithReturnFunc(t *testing.T) {
	snippet := `
fun fib(n) {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
}

for (var i = 0; i < 20; i = i + 1) {
  print fib(i);
}
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
var a = "你好世界！!";
a = "我是练习时长两年半的练习生";
print a;
a=2; // no
print a;`
	vm := &VM.VM{}
	vm.RunStr(snippet)
}
