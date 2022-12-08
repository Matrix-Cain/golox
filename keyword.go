package main

var KeyWords map[string]TokenType

func init() {
	KeyWords = make(map[string]TokenType, 20)
	KeyWords["and"] = AND
	KeyWords["class"] = CLASS
	KeyWords["else"] = ELSE
	KeyWords["false"] = FALSE
	KeyWords["for"] = FOR
	KeyWords["fun"] = FUN
	KeyWords["if"] = IF
	KeyWords["nil"] = NIL
	KeyWords["or"] = OR
	KeyWords["print"] = PRINT
	KeyWords["return"] = RETURN
	KeyWords["super"] = SUPER
	KeyWords["this"] = THIS
	KeyWords["true"] = TRUE
	KeyWords["var"] = VAR
	KeyWords["while"] = WHILE
}
