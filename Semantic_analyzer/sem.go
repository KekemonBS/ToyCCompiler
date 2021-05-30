package semant

import (
	"fmt"

	"kekemon.org/comp/parser"
)

//Seman modifie AST to get decorated AST
func Seman(ast *parser.Node) *parser.Node {
	fmt.Printf("\n[AST]\n\n")
	parser.Display(ast, 0, true)
	return ast
}
