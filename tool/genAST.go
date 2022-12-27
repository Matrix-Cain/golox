package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Error("Usage: go run genAST.go <output directory>")
		os.Exit(64)
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Binary   :  Left Expr, Operator Token, Right Expr",
		"Grouping :  Expression Expr",
		"Literal  :  Type TokenType// golang is static typed language cache type to avoid unnecessary type switch cost, Value interface{}",
		"Unary    :  Operator Token, Right Expr",
		"Ternary  :  ConditionalExpr Expr, ThenExpr Expr, ElseExpr Expr"})

	path := outputDir + "/" + "expr.go"
	err := exec.Command("go", "fmt", path).Run()
	if err != nil {
		log.Error(err)
		return
	}

	//defineAst(outputDir, "Stmt", []string{
	//	"Expression :  expression Expr",
	//	"Print      :  expression Expr"})
}

func defineAst(outputDir string, baseName string, types []string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error("Cannot create or open file with error " + err.Error())
		os.Exit(64)
	}
	w := bufio.NewWriter(fd)
	w.WriteString("package ast\n")
	w.WriteString("\n")

	w.WriteString("import .\"golox/lox/lexer\"")
	w.WriteString("\n")

	w.WriteString("\n")

	w.WriteString("type " + baseName + " interface {\n")
	w.WriteString("    Accept(v Visitor) (interface{}, error)\n")
	w.WriteString("}\n")

	w.WriteString("\n")

	defineVisitor(w, baseName, types)

	for _, type0 := range types {
		typeName := strings.TrimSpace(strings.Split(type0, ":")[0])
		fields := strings.TrimSpace(strings.Split(type0, ":")[1])
		fieldList := strings.Split(fields, ",")

		defineType(w, baseName, typeName, fieldList)
		defineTypeMethod(w, baseName, typeName)
	}
	w.WriteString("\n")
	w.Flush()
	fd.Close()

}

func defineVisitor(w *bufio.Writer, baseName string, types []string) {
	w.WriteString("type Visitor interface {\n")

	for _, typeName := range types {
		typeName = strings.TrimSpace(strings.Split(typeName, ":")[0])

		w.WriteString("    Visit" + typeName + baseName + "(expr *" + typeName + ") (interface{}, error)")
		w.WriteString("\n")
	}
	w.WriteString("}\n")
	w.WriteString("\n")
}

func defineType(w *bufio.Writer, baseName string, className string, fieldList []string) {
	w.WriteString("type " + className + " struct {\n")
	for _, field := range fieldList {
		field = strings.TrimSpace(field)
		w.WriteString("    " + field + "\n")
	}
	w.WriteString("}\n")
	w.WriteString("\n")
}

func defineTypeMethod(w *bufio.Writer, baseName string, className string) {
	w.WriteString("func (t *" + className + ") Accept(v Visitor) (interface{}, error) {\n")
	w.WriteString("    return v.Visit" + className + baseName + "(t)\n")
	w.WriteString("}\n")
	w.WriteString("\n")
}
