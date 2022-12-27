package tests

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestZH(t *testing.T) {
	const nihongo = "日本語"
	for i, w := 0, 0; i < len(nihongo); i += w {
		runeValue, width := utf8.DecodeRuneInString(nihongo[i:])
		fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
		w = width
	}
}

func TestPrint(t *testing.T) {
	fmt.Println(nil)
}
