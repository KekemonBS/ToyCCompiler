package main

import (
	"fmt"
	"os"

	"kekemon.org/comp/codegen"
	"kekemon.org/comp/lexer"
	"kekemon.org/comp/parser"
	"kekemon.org/comp/semant"
)

func main() {
	path, _ := os.Getwd()
	//os.Args[1:][0]
	tokens, description, positions := lexer.Tokenize(path + "/КР-07-Golang-IV-81-Buhtiy.cpp")
	ast := parser.Parse(tokens, description, positions)
	decoast := semant.Seman(ast)
	resultingCode := codgen.CodGen(decoast)
	fmt.Print(resultingCode, "\n\n")
	out, _ := os.Create(path + "/КР-07-Golang-IV-81-Buhtiy.asm")
	out.WriteString(resultingCode)
}
