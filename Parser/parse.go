package parser

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

var tokens []string
var description []string
var position []string

//Mnode is main tree that will be formed and shrinked
var Mnode Node

/*Parse function transforms token array into Parse tree
than shrinks it to Abstract syntax tree*/
func Parse(tokenArr []string, descriptionArr []string, positionArr []string) *Node {
	fmt.Println("\n[Tokens]")
	for i := 0; i < len(tokenArr); i++ {
		fmt.Printf(" \n %-30s   ------------    %7s %10s", tokenArr[i], descriptionArr[i], positionArr[i])
	}
	fmt.Printf("\n\n")
	tokens = tokenArr
	description = descriptionArr
	position = positionArr

	program(tokens, description, position)
	fmt.Printf("\n[Parse tree]\n\n")
	Display(&Mnode, 0, true)

	transformTree(&Mnode, nil)
	return &Mnode
}

//Node is part of syntax tree ,make it private (change to node (will cause headache))
type Node struct {
	Value    string
	Position string
	Children []*Node
}

func program(tok []string, desc []string, pos []string) { //1. program → declarationList (обовязково)
	Mnode.Value = "program"
	declarationList(&Mnode, tok, desc, pos)
}

func declarationList(n *Node, slice []string, descslice []string, posslice []string) { //2. declarationList → declarationList declaration | declaration(1:поки тільки останнє)
	var node Node
	node.Value = "declarationList"
	n.Children = append(n.Children, &node)

	var innode Node
	innode.Value = "in"
	node.Children = append(node.Children, &innode)

	var open int
	var close int
	var inprth bool
	var depth int

	for i := 0; i < len(slice); i++ {
		if descslice[i] == "OPNPRTH" || descslice[i] == "OPNBRC" {
			depth++
		} else if descslice[i] == "CLSPRTH" || descslice[i] == "CLSBRC" {
			depth--
		}
		if depth > 0 {
			inprth = true
		} else {
			inprth = false
		}

		if descslice[i] == "TYPE" && !inprth {
			open = i
		} else if (descslice[i] == "CLSBRC" || descslice[i] == "SEMCOL") && !inprth {
			close = i
			declaration(&node, slice[open:close+1], descslice[open:close+1], posslice[open:close+1])
		}
	}
}

func declaration(n *Node, slice []string, descslice []string, posslice []string) { //3. declaration → varDeclaration | funDeclaration(1:поки тільки останнє)
	var node Node
	node.Value = "declaration"
	n.Children = append(n.Children, &node)

	if descslice[0] == "TYPE" && descslice[len(descslice)-1] == "CLSBRC" {
		funDeclaration(&node, slice, descslice, posslice)
	} else {
		varDeclaration(&node, slice, descslice, posslice)
	}
}

func varDeclaration(n *Node, slice []string, descslice []string, posslice []string) { //4. varDeclaration → typeSpecifier varDeclList ;
	var node Node
	node.Value = "varDeclaration"
	n.Children = append(n.Children, &node)

	if descslice[0] == "TYPE" || descslice[len(descslice)-1] == "SEMCOL" {
		typeSpecifier(&node, slice[0], descslice[0], posslice[0])
		varDeclList(&node, slice[1:len(descslice)-1], descslice[1:len(descslice)-1], posslice[1:len(descslice)-1])
	}

}

func scopedVarDeclaration(n *Node, slice []string, descslice []string, posslice []string) { //5. scopedVarDeclaration → static typeSpecifier varDeclList ; | typeSpecifier varDeclList(це) ;
	var node Node
	node.Value = "scopedVarDeclaration"
	n.Children = append(n.Children, &node)
	//Add distinction someday
	if descslice[0] == "TYPE" || descslice[len(descslice)-1] == "SEMCOL" {
		typeSpecifier(&node, slice[0], descslice[0], posslice[0])
		varDeclList(&node, slice[1:len(descslice)-1], descslice[1:len(descslice)-1], posslice[1:len(descslice)-1])
	}
}

func varDeclList(n *Node, slice []string, descslice []string, posslice []string) { //6. varDeclList → varDeclList , varDeclInitialize | varDeclInitialize (тільки це поки)
	var node Node
	node.Value = "varDeclList"
	n.Children = append(n.Children, &node)

	varDeclInitialize(&node, slice, descslice, posslice)
}

func varDeclInitialize(n *Node, slice []string, descslice []string, posslice []string) { //7. varDeclInitialize → varDeclId | varDeclId = simpleExpression
	var node Node
	node.Value = "varDeclInitialize"
	n.Children = append(n.Children, &node)

	var foundeq bool
	for i := 0; i < len(slice); i++ {
		if descslice[i] == "EQ" {
			foundeq = true
			varDeclID(&node, slice[:i], descslice[:i], posslice[:i])
			simpleExpression(&node, slice[i+1:], descslice[i+1:], posslice[i+1:])
		}
	}
	if !foundeq {
		varDeclID(&node, slice, descslice, posslice)
	}
}

func varDeclID(n *Node, slice []string, descslice []string, posslice []string) { //8. varDeclID → ID | ID [ NUMCONST ]
	var node Node
	node.Value = "varDeclId"
	n.Children = append(n.Children, &node)

	//Add someday distinction betwen two possibilities
	var id Node
	id.Value = slice[0]
	node.Children = append(node.Children, &id)
}

func typeSpecifier(n *Node, item string, description string, position string) { //9. typeSpecifier → int | bool | char (узнать який тип поставить в дерево)
	var node Node
	node.Value = "typeSpecifier"
	n.Children = append(n.Children, &node)

	var specconstnode Node
	specconstnode.Value = item
	node.Children = append(node.Children, &specconstnode)
}

func funDeclaration(n *Node, slice []string, descslice []string, posslice []string) { //10. funDeclaration → typeSpecifier ID ( params ) statement (1:поки тільки перше) | ID ( params ) statement
	//if first is ID then 1 case :
	var node Node
	node.Value = "funDeclaration"
	n.Children = append(n.Children, &node)

	var funcnode Node
	funcnode.Value = "infunc"
	node.Children = append(node.Children, &funcnode)

	var open int
	var close int
	var inprth bool

	//hadparams := false
	for ind, i := range descslice {
		if i == "OPNPRTH" {
			inprth = true
			open = ind + 1
		} else if i == "CLSPRTH" {
			close = ind
			params(&node, slice[open:close], descslice[open:close], posslice[open:close])
			inprth = false
			//~ hadparams = true
			break
		} else if i == "TYPE" && !inprth {
			typeSpecifier(&node, slice[ind], i, posslice[ind])
		} else if i == "ID" && !inprth {
			var idnode Node
			idnode.Value = "funcIdentifier"
			node.Children = append(node.Children, &idnode)

			var idspecnode Node
			idspecnode.Value = slice[ind]
			idspecnode.Position = posslice[ind]
			idnode.Children = append(idnode.Children, &idspecnode)

		}
	}
	if !inprth {
		statement(&node, slice[close+1:], descslice[close+1:], posslice[close+1:])

	}
}

func params(n *Node, slice []string, descslice []string, posslice []string) { //11. params → paramList | nothing (nothing)
	var node Node
	node.Value = "params"
	n.Children = append(n.Children, &node)

	if len(slice) == 0 {
		var eNode Node
		eNode.Value = ""
		node.Children = append(node.Children, &eNode)
	} else {
		paramList(&node, slice, descslice, posslice)
	}
}

func paramList(n *Node, slice []string, descslice []string, posslice []string) { //12. paramList → paramList ; paramTypeList | paramTypeList
	var node Node
	node.Value = "paramList"
	n.Children = append(n.Children, &node)

	var found bool
	for i := len(descslice) - 1; i >= 0; i-- {
		if descslice[i] == "SEMCOL" {
			found = true
			paramSepList(&node, slice[:i], descslice[:i], posslice[:i])
			paramList(&node, slice[i+1:], descslice[i+1:], posslice[i+1:])
		}
	}
	if !found {
		paramSepList(&node, slice, descslice, posslice)
	}
}

func paramSepList(n *Node, slice []string, descslice []string, posslice []string) { //13. paramSepList → param {{ , param }}
	var node Node
	node.Value = "paramSepList"
	n.Children = append(n.Children, &node)

	var foundcomma bool
	var open int
	var close int
	for i := 0; i < len(descslice); i++ {
		if descslice[i] == "COMMA" {
			close = i
			param(&node, slice[open:close], descslice[open:close], posslice[open:close])
			open = close
			foundcomma = true
			continue
		}
		// if descslice[i] == "COMMA" && close != 0 {
		// 	close = i
		// 	param(&node, slice[open+1:close], descslice[open+1:close], posslice[open+1:close])
		// 	open = close
		// 	foundcomma = true
		// }
	}
	if open != 0 && close == open {
		param(&node, slice[open+1:], descslice[open+1:], posslice[open+1:])

	}
	if !foundcomma {
		param(&node, slice, descslice, posslice)
	}

}

func param(n *Node, slice []string, descslice []string, posslice []string) { //14. param → typeSpecifier paramId
	var node Node
	node.Value = "param"
	n.Children = append(n.Children, &node)

	for ind, i := range descslice {
		if i == "TYPE" {
			typeSpecifier(&node, slice[ind], descslice[ind], posslice[ind])
			paramID(&node, slice[ind+1:], descslice[ind+1:], posslice[ind+1:])
		}
	}
}

/* func paramIDList(n *Node, slice []string, descslice []string, posslice []string ) { //14. paramIdList → paramIdList , paramId | paramId
	var node Node
	node.Value = "paramIDList"
	n.Children = append(n.Children, &node)

	for i := len(descslice) - 1; i >= 0; i-- {
		if descslice[i] == "COMMA" {
			paramID(&node, slice[i+1:], descslice[i+1:])
			paramIDList(&node, slice[:i], descslice[:i])

		}
	}
} */

func paramID(n *Node, slice []string, descslice []string, posslice []string) { //15. paramId → ID | ID []
	var node Node
	node.Value = "paramID"
	n.Children = append(n.Children, &node)

	if len(slice) == 1 {
		var id Node
		id.Value = slice[0]
		node.Children = append(node.Children, &id)
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------
func statement(n *Node, slice []string, descslice []string, posslice []string) { //16. statement → expressionStmt | compoundStmt(1_це) | selectionStmt(2_і це) | iterationStmt | returnStmt (1_і це) | breakStmt
	var node Node
	node.Value = "statement"
	n.Children = append(n.Children, &node)

	var incompound bool
	var compoundDepth int
	var inIterationStmt bool

	var found bool
	var open int
	var close int

	for ind, i := range descslice {
		if compoundDepth > 0 {
			incompound = true
		} else {
			incompound = false
		}

		if !incompound {
			//numberdescr := []string{"HEX", "OCT", "FLOAT", "INT", "BIN"}
			if i == "RET" {
				returnStmt(&node, slice, descslice, posslice)
				found = true
			} else if i == "ID" && !found && !inIterationStmt && !incompound /*|| in(i, numberdescr)*/ {
				expressionStmt(&node, slice, descslice, posslice)
				break
			} else if i == "DO" {
				inIterationStmt = true
				open = ind
			} else if inIterationStmt && !incompound && i == "SEMCOL" {
				close = ind
				iterationStmt(&node, slice[open:close+1], descslice[open:close+1], posslice[open:close+1])
				inIterationStmt = false
			}
		}
		if i == "OPNBRC" {
			compoundDepth++
			if !incompound && !inIterationStmt {
				open = ind
			}
		} else if i == "CLSBRC" {
			compoundDepth--
			if compoundDepth == 0 && !inIterationStmt {
				close = ind
				compoundStmt(&node, slice[open:close+1], descslice[open:close+1], posslice[open:close+1])
			}
		}
	}
	if open > close {
		fmt.Println("Unclosed ")
		if incompound {
			fmt.Print("compound statement tahat starts at line : " + posslice[open])
			fmt.Println("Expecting closing braces")
			os.Exit(1)
		} else if inIterationStmt {
			fmt.Print("iteration statement tahat starts at line : " + posslice[open])
			fmt.Println("Expecting semicolon")
			os.Exit(1)
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------
func expressionStmt(n *Node, slice []string, descslice []string, posslice []string) { //17. expressionStmt → expression ; | ;
	var node Node
	node.Value = "expressionStmt"
	n.Children = append(n.Children, &node)

	if descslice[0] == "SEMCOL" {
		var exprnode Node
		exprnode.Value = slice[0]
		node.Children = append(node.Children, &exprnode)
	} else if len(slice) != 0 {
		expression(&node, slice[:len(slice)-1], descslice[:len(slice)-1], posslice[:len(slice)-1])
		//var exprnode Node
		//exprnode.Value = slice[len(slice)-1]
		//node.Children = append(node.Children, &exprnode)
	} else {
		fmt.Println("Error somewhere in expressionStmt")
	}
}

func compoundStmt(n *Node, slice []string, descslice []string, posslice []string) { //18. compoundStmt → { localDeclarations statementList(statement) }
	var node Node
	node.Value = "compoundStmt"
	n.Children = append(n.Children, &node)
	//if len(slice) > 0 {
	//	n.Children = append(n.Children, &node)
	//	statementList(&node, slice, descslice, posslice)
	//}
	var innode Node
	innode.Value = "in"
	node.Children = append(node.Children, &innode)

	slice = slice[1 : len(slice)-1]
	descslice = descslice[1 : len(descslice)-1]
	posslice = posslice[1 : len(posslice)-1]

	var locdecl bool
	var statementdecl bool
	var nesting int
	var open int
	var close int
	//localDeclarations(&node, slice[open:close+1], descslice[open:close+1], posslice[open:close+1])
	//statement(&node, slice[open:close+1], descslice[open:close+1], posslice[open:close+1])
	for ind, val := range descslice {
		if val == "TYPE" && !locdecl && nesting == 0 {
			locdecl = true
			open = ind
		} else if !statementdecl && !locdecl {
			statementdecl = true
			open = ind
		}
		if val == "OPNBRC" || val == "OPNPRTH" {
			nesting++
		} else if val == "CLSBRC" || val == "CLSPRTH" {
			nesting--
		}
		if nesting == 0 {
			if val == "CONT" {
				var contnode Node
				contnode.Value = "continue"
				contnode.Position = posslice[ind]
				node.Children = append(node.Children, &contnode)
			} else if val == "BR" {
				var brnode Node
				brnode.Value = "break"
				brnode.Position = posslice[ind]
				node.Children = append(node.Children, &brnode)
			}
		}
		if nesting == 0 && val == "SEMCOL" {
			if locdecl {
				close = ind
				localDeclarations(&node, slice[open:close+1], descslice[open:close+1], posslice[open:close+1])
				locdecl = false
			} else {
				close = ind
				statement(&node, slice[open:close+1], descslice[open:close+1], posslice[open:close+1])
				statementdecl = false
			}
		}
	}
}

func localDeclarations(n *Node, slice []string, descslice []string, posslice []string) { //19. localDeclarations → localDeclarations <scopedVarDeclaration>(тільки те що в <>) | nothing
	var node Node
	node.Value = "localDeclarations"
	n.Children = append(n.Children, &node)
	scopedVarDeclaration(&node, slice, descslice, posslice)

}

func statementList(n *Node, slice []string, descslice []string, posslice []string) { //20. statementList → statementList statement | nothing
	var node Node
	node.Value = "statementList"

	if len(slice) > 0 {
		n.Children = append(n.Children, &node)
		statement(&node, slice, descslice, posslice)
	}
	fmt.Println(slice)
	//Add detection of both possibilities

}

func a21() { //21.

}

//-------------------------------------------------------------------------------------------------------------------------------------
func iterationStmt(n *Node, slice []string, descslice []string, posslice []string) { //22. iterationStmt → while ( simpleExpression ) statement | for ( ID in ID ) statement
	var node Node
	node.Value = "iterationStmt"
	n.Children = append(n.Children, &node)

	var iternode Node
	iternode.Value = "initeration"
	node.Children = append(node.Children, &iternode)

	var incompound bool
	var compoundDepth int

	var inDoWhile bool
	var do int
	var while int
	var semcol int

	for ind, val := range descslice {
		if compoundDepth > 0 {
			incompound = true
		} else {
			incompound = false
		}

		if val == "DO" && !inDoWhile {
			inDoWhile = true
			do = ind
		} else if val == "WHL" && inDoWhile {
			while = ind
		} else if val == "SEMCOL" && inDoWhile && !incompound {
			if do >= while {
				fmt.Println("Iteration statement lacks \"while\" after \"do\" at line : " + posslice[do])
				os.Exit(1)
			}

			semcol = ind

			if semcol-while <= 3 {
				fmt.Println("No expression after \"while\" at line : " + posslice[while])
				os.Exit(1)
			}
			statement(&node, slice[do+1:while], descslice[do+1:while], posslice[do+1:while])

			var itercondnode Node
			itercondnode.Value = "iterCondition"
			node.Children = append(node.Children, &itercondnode)

			simpleExpression(&node, slice[while+2:semcol-1], descslice[while+2:semcol-1], posslice[while+2:semcol-1])
		}

		if val == "OPNBRC" {
			compoundDepth++
		} else if val == "CLSBRC" {
			compoundDepth--
		}

	}
}

//-------------------------------------------------------------------------------------------------------------------------------------

func returnStmt(n *Node, slice []string, descslice []string, posslice []string) { //23. returnStmt → return ; | return expression ;(поки так : 2)
	var node Node
	node.Value = "returnStmt"
	n.Children = append(n.Children, &node)

	var open int
	var close int
	if descslice[len(descslice)-1] != "SEMCOL" {
		fmt.Println("No semicolon in return statement.")
		os.Exit(1)
	} else if descslice[0] == "RET" && descslice[1] == "SEMCOL" {
		var retnode Node
		retnode.Value = slice[0]
		node.Children = append(node.Children, &retnode)
	} else {
		for ind, i := range descslice {
			if i == "RET" {
				open = ind + 1
			} else if i == "SEMCOL" {
				close = ind
				expression(&node, slice[open:close], descslice[open:close], posslice[open:close])
			}
		}
	}
}

func a() { //24.

}

func expression(n *Node, slice []string, descslice []string, posslice []string) { /*  25. expression → mutable = expression  | mutable += expression | mutable −= expression
	 *  | mutable ∗= expression | mutable /= expression | mutable ++ | mutable −−
	    | simpleExpression (ось це)  */
	var node Node
	node.Value = "expression"
	n.Children = append(n.Children, &node)

	var operator string
	var found bool
	for ind, i := range descslice {
		if i == "EQ" && (descslice[ind-1] == "ID" && descslice[ind+1] != "EQ") { //ideally second condition should be ismutable() or similar
			found = true

			var eq Node
			eq.Value = slice[ind]
			node.Children = append(node.Children, &eq)

			expression(&eq, slice[ind+1:], descslice[ind+1:], posslice[ind+1:])
			mutable(&eq, slice[:ind], descslice[:ind], posslice[:ind])
		}
		if i == "DIV" || i == "EQ" {
			operator += slice[ind]
		}
		operators := []string{"+=", "-=", "*=", "/="}
		if in(operator, operators) {
			found = true

			var opernode Node
			opernode.Value = operator
			node.Children = append(node.Children, &opernode)

			expression(&opernode, slice[ind+1:], descslice[ind+1:], posslice[ind+1:])
			mutable(&opernode, slice[:ind-1], descslice[:ind-1], posslice[:ind-1])
			break
		}
	}

	if !found {
		simpleExpression(&node, slice, descslice, posslice)
	}
	//add distinction one day

}

func simpleExpression(n *Node, slice []string, descslice []string, posslice []string) { //26. simpleExpression → simpleExpression ‘|’ andExpression | andExpression(це)
	var node Node
	node.Value = "simpleExpression"
	n.Children = append(n.Children, &node)

	var found bool
	for i := len(descslice) - 1; i >= 0; i-- {
		if descslice[i] == "OR" {
			found = true
			simpleExpression(&node, slice[:i], descslice[:i], posslice[:i])
			andExpression(&node, slice[i+1:], descslice[i+1:], posslice[i+1:])
			break
		}
	}
	if !found {
		andExpression(&node, slice, descslice, posslice)
	}
}

func andExpression(n *Node, slice []string, descslice []string, posslice []string) { //27. andExpression → andExpression & unaryRelExpression | unaryRelExpression(це)
	var node Node
	node.Value = "andExpression"
	n.Children = append(n.Children, &node)
	unaryRelExpression(&node, slice, descslice, posslice)
}

func unaryRelExpression(n *Node, slice []string, descslice []string, posslice []string) { //28. unaryRelExpression → ! unaryRelExpression | relExpression(це)
	var node Node
	node.Value = "unaryRelExpression"
	n.Children = append(n.Children, &node)

	if descslice[0] == "NOT" {
		var notnode Node
		notnode.Value = "!"
		node.Children = append(node.Children, &notnode)
		unaryRelExpression(&node, slice[1:], descslice[1:], posslice[1:])
	} else {
		relExpression(&node, slice, descslice, posslice)
	}
}

func relExpression(n *Node, slice []string, descslice []string, posslice []string) { //29. relExpression → sumExpression relop sumExpression | sumExpression(це)
	var node Node
	node.Value = "relExpression"
	n.Children = append(n.Children, &node)

	var symbol string
	var open int
	var close int
	var inprth bool
	var prthdepth int

	for i := 0; i < len(slice); i++ {
		if descslice[i] == "OPNPRTH" {
			prthdepth++
		} else if descslice[i] == "CLSPRTH" {
			prthdepth--
		}
		if prthdepth > 0 {
			inprth = true
		} else {
			inprth = false
		}

		elem := descslice[i]
		if (elem == "NOT" || elem == "EQ" || elem == "LS" || elem == "GR") &&
			!inprth {
			symbol += slice[i]
			if open != 0 {
				close = i
				break
			}
			open = i
		}
	}
	if open != 0 && close == 0 {
		close = open
	}
	relops := []string{">", "<", ">=", "<=", "==", "!="}
	if in(symbol, relops) {
		sumExpression(&node, slice[:open], descslice[:open], posslice[:open])
		relop(&node, symbol, descslice[open:close+1], posslice[open:close+1])
		sumExpression(&node, slice[close+1:], descslice[close+1:], posslice[close+1:])
	}
	if open == 0 {
		sumExpression(&node, slice, descslice, posslice)
	}

}

func relop(n *Node, slice string, descslice []string, posslice []string) { //30. relop → <= | < | > | >= | == | ! =
	n.Value = slice

	/* 	var node Node
	node.Value = slice
	n.Children = append(n.Children, &node) */
}

func sumExpression(n *Node, slice []string, descslice []string, posslice []string) { //31. sumExpression → sumExpression sumop mulExpression | mulExpression(це)
	var node Node
	node.Value = "sumExpression"
	n.Children = append(n.Children, &node)

	var inprth bool
	var prthdepth int
	var found bool
	var where int
	for i := len(slice) - 1; i >= 0; i-- {
		if descslice[i] == "CLSPRTH" {
			prthdepth++
		} else if descslice[i] == "OPNPRTH" {
			prthdepth--
		}
		if prthdepth > 0 {
			inprth = true
		} else {
			inprth = false
		}
		operations := []string{"SUB", "BITCOMPL"}
		if (descslice[i] == "ADD" || descslice[i] == "SUB") && !inprth &&
			i > 0 && (!in(descslice[i-1], operations) || (i > 0 && descslice[i-1] == "")) {
			found = true
			where = i
			if len(slice[:where]) > 0 && len(slice[where+1:]) > 0 {
				sumExpression(&node, slice[:where], descslice[:where], posslice[:where])
				sumop(&node, slice[where], descslice[where], posslice[where])
				mulExpression(&node, slice[where+1:], descslice[where+1:], posslice[where+1:])
				break
			} else {
				unaryExpression(&node, slice, descslice, posslice)
				break
			} //
		}
	}
	if !found {
		mulExpression(&node, slice, descslice, posslice)
	}
}

func sumop(n *Node, slice string, descslice string, position string) { //32. sumop → + | −
	n.Value = slice

	/* 	var node Node
	node.Value = slice
	n.Children = append(n.Children, &node) */
}

func mulExpression(n *Node, slice []string, descslice []string, posslice []string) { //33. mulExpression → mulExpression mulop unaryExpression | unaryExpression(це)
	var node Node
	node.Value = "mulExpression"
	n.Children = append(n.Children, &node)

	var inprth bool
	var prthdepth int
	var found bool
	var where int
	for i := len(slice) - 1; i >= 0; i-- {
		if descslice[i] == "CLSPRTH" {
			prthdepth++
		} else if descslice[i] == "OPNPRTH" {
			prthdepth--
		}
		if prthdepth > 0 {
			inprth = true
		} else {
			inprth = false
		}

		if (descslice[i] == "MUL" || descslice[i] == "DIV") && !inprth {
			found = true
			where = i
			mulExpression(&node, slice[:where], descslice[:where], posslice[:where])
			mulop(&node, slice[where], descslice[where], posslice[where])
			unaryExpression(&node, slice[where+1:], descslice[where+1:], posslice[where+1:])
			break

		}
	}
	if !found {
		unaryExpression(&node, slice, descslice, posslice)
	}

}

func mulop(n *Node, slice string, descslice string, position string) { //34.mulop → ∗ | / | %
	n.Value = slice

	/* 	var node Node
	node.Value = slice
	n.Children = append(n.Children, &node) */
}

func unaryExpression(n *Node, slice []string, descslice []string, posslice []string) { //35. unaryExpression → unaryop unaryExpression | factor(це)
	var node Node
	node.Value = "unaryExpression"
	n.Children = append(n.Children, &node)

	if descslice[0] == "BITCOMPL" || descslice[0] == "SUB" { //someday add support for left unary operators if needed
		unaryop(&node, slice[0], descslice[0], posslice[0])
		unaryExpression(&node, slice[1:], descslice[1:], posslice[1:])

	} else {
		conditionalExpression(&node, slice, descslice, posslice)
	}

}

func unaryop(n *Node, slice string, descslice string, position string) { //36.unaryop → − | ∗ | ?
	n.Value = "un" + slice

	/* 	var node Node
	node.Value = slice
	n.Children = append(n.Children, &node) */
}

func conditionalExpression(n *Node, slice []string, descslice []string, posslice []string) { //35_1. conditionalExpression → simpleExpression ? expression : expression | factor(це)
	var node Node
	node.Value = "conditionalExpression"
	n.Children = append(n.Children, &node)

	var qmind int
	var colind int

	var exprdepth int
	for i := len(slice) - 1; i >= 0; i-- {
		if descslice[i] == "COL" {
			if exprdepth == 0 {
				colind = i
			}
			exprdepth++
		} else if descslice[i] == "QM" {
			exprdepth--
			if exprdepth == 0 {
				qmind = i
			}
		}
	}
	if colind-qmind == 1 {
		fmt.Println("No expression that matches true condition.")
		os.Exit(1)
	} else if qmind == 0 && colind != 0 {
		fmt.Println("Syntax error in conditional expression.")
		os.Exit(1)
	}
	if qmind != 0 {
		if colind != 0 {
			var cond Node
			cond.Value = "condition"
			node.Children = append(node.Children, &cond)

			simpleExpression(&cond, slice[:qmind], descslice[:qmind], posslice[:qmind])
			//condop(&node, slice[qmind], descslice[qmind], posdescslice[colind])
			var truenode Node
			truenode.Value = "booltrue"
			node.Children = append(node.Children, &truenode)

			expression(&truenode, slice[qmind+1:colind], descslice[qmind+1:colind], posslice[qmind+1:colind])
			//condop(&node, slice[colind], descslice[colind], posdescslice[colind])
			var falsenode Node
			falsenode.Value = "boolfalse"
			node.Children = append(node.Children, &falsenode)

			expression(&falsenode, slice[colind+1:], descslice[colind+1:], posslice[colind+1:])
			qmind = 0
			colind = 0
		} else {
			var cond Node
			cond.Value = "condition"
			node.Children = append(node.Children, &cond)

			simpleExpression(&cond, slice[:qmind], descslice[:qmind], posslice[:qmind])
			//condop(&node, slice[qmind], descslice[qmind])
			var truenode Node
			truenode.Value = "booltrue"
			node.Children = append(node.Children, &truenode)

			expression(&truenode, slice[qmind+1:colind], descslice[qmind+1:colind], posslice[qmind+1:colind])
			qmind = 0
		}
	} else {
		factor(&node, slice, descslice, posslice)
	}
}

//func condop(n *Node, slice string, descslice string, position string) { //36.condop → ? | :
//	n.Value = "cond" + slice
//
//	/* 	var node Node
//	node.Value = slice
//	n.Children = append(n.Children, &node) */
//}

func factor(n *Node, slice []string, descslice []string, posslice []string) { //37. factor → immutable | mutable
	var node Node
	node.Value = "factor"
	n.Children = append(n.Children, &node)
	if descslice[0] == "ID" && len(slice) == 1 {
		mutable(&node, slice, descslice, posslice)
	} else {
		immutable(&node, slice, descslice, posslice)
	}
}

func mutable(n *Node, slice []string, descslice []string, posslice []string) { //38. mutable → ID | mutable [ expression ]
	var node Node
	node.Value = "mutable"
	n.Children = append(n.Children, &node)

	//Add someday distinction betwen two possibilities
	var id Node
	id.Value = slice[0]
	id.Position = posslice[0]
	node.Children = append(node.Children, &id)

}

func immutable(n *Node, slice []string, descslice []string, posslice []string) { //39. immutable → ( expression ) | call | constant
	var node Node
	node.Value = "immutable"
	n.Children = append(n.Children, &node)

	if len(descslice) > 1 && descslice[0] == "OPNPRTH" && descslice[len(slice)-1] == "CLSPRTH" {
		expression(&node, slice[1:len(slice)-1], descslice[1:len(slice)-1], posslice[1:len(slice)-1])
	} else if len(descslice) > 1 && descslice[0] == "ID" && descslice[1] == "OPNPRTH" {
		call(&node, slice, descslice, posslice)
	} else {
		constant(&node, slice, descslice, posslice)
	}
}

func call(n *Node, slice []string, descslice []string, posslice []string) { //40. call → ID ( args )
	var node Node
	node.Value = "call"
	node.Position = posslice[1]
	n.Children = append(n.Children, &node)

	var callnode Node
	callnode.Value = "incall"
	node.Children = append(node.Children, &callnode)

	var inprth bool
	var open int
	var close int
	for i := len(slice) - 1; i >= 0; i-- {
		if descslice[i] == "ID" && !inprth {
			var idnode Node
			idnode.Value = slice[i]
			idnode.Position = posslice[i]
			node.Children = append(node.Children, &idnode)
		} else if descslice[i] == "CLSPRTH" {
			close = i
			inprth = true
		} else if descslice[i] == "OPNPRTH" {
			open = i
			args(&node, slice[open+1:close], descslice[open+1:close], posslice[open+1:close])
			inprth = false
		}

	}
}

func args(n *Node, slice []string, descslice []string, posslice []string) { //41. args → argList | nothing
	var node Node
	node.Value = "args"
	n.Children = append(n.Children, &node)

	if len(descslice) != 0 {
		argList(&node, slice, descslice, posslice)
	} else {
		var argsnode Node
		argsnode.Value = ""
		node.Children = append(node.Children, &argsnode)
	}

}

func argList(n *Node, slice []string, descslice []string, posslice []string) { //42. argList → argList , expression | expression
	var node Node
	node.Value = "argList"
	n.Children = append(n.Children, &node)

	var found bool
	var where int
	for i := len(slice) - 1; i >= 0; i-- {
		if descslice[i] == "COMMA" {
			found = true
			where = i
			expression(&node, slice[where+1:], descslice[where+1:], posslice[where+1:])
			argList(&node, slice[:where], descslice[:where], posslice[:where])
			break
		}
	}
	if !found {
		expression(&node, slice, descslice, posslice)
	}
}

func constant(n *Node, slice []string, descslice []string, posslice []string) { //43. constant → NUMCONST | CHARCONST | STRINGCONST | true | false
	var node Node
	node.Value = "constant"
	n.Children = append(n.Children, &node)

	if descslice[0] == "INT" {
		var numconstnode Node
		numconstnode.Value = slice[0]
		node.Children = append(node.Children, &numconstnode)
	} else if descslice[0] == "HEX" {
		var numconstnode Node
		val, _ := strconv.ParseInt(slice[0], 0, 64)
		numconstnode.Value = strconv.Itoa(int(val))
		node.Children = append(node.Children, &numconstnode)

	} else if descslice[0] == "OCT" {
		var numconstnode Node
		val, _ := strconv.ParseInt(slice[0], 0, 64)
		numconstnode.Value = strconv.Itoa(int(val))
		node.Children = append(node.Children, &numconstnode)
	} else if descslice[0] == "FLOAT" {
		var numconstnode Node
		val, _ := strconv.ParseFloat(slice[0], 64)
		numconstnode.Value = strconv.Itoa(int(math.Round(val)))
		node.Children = append(node.Children, &numconstnode)
	} else {
		fmt.Println("Unsupported return type at line : " + posslice[0])
		os.Exit(1)
	}

}

var idented map[int]bool = make(map[int]bool)

//Display displays any tree
func Display(n *Node, ident int, last bool) {
	if len(n.Children) > 1 {
		idented[ident] = true
	}
	ident++
	s := n.Value
	if ident > 0 {
		for i := 0; i < ident-1; i++ {
			if i < ident-2 {
				if idented[i] {
					fmt.Print("\u2502 ")
				} else {
					fmt.Print("   ")
				}
			} else if last == false {
				fmt.Print("\u251c\u2500\u2500 ")
			} else {
				fmt.Print("\u2514\u2500\u2500 ")
				idented[i] = false
			}
		}
		fmt.Print("")
	}
	fmt.Printf("%s\n", s)
	for i := 0; i < len(n.Children); i++ {
		if i == len(n.Children)-1 {
			Display(n.Children[i], ident, true)
		} else {
			Display(n.Children[i], ident, false)
		}
	}
}

//transformTree rewires node if it has <= 1 children and it is not in exeptions (shrinks tree)
func transformTree(node *Node, Parent *Node) {
	exceptions := []string{"typeSpecifier", "varDeclId", "funcIdentifier", "returnStmt", "program",
		"identifier", "un~", "un-", "condition", "booltrue", "boolfalse", "args", "paramSepList"}

	for _, i := range node.Children {
		if len(node.Children) > 1 || in(node.Value, exceptions) {
			transformTree(i, node)
		} else if len(node.Children) <= 1 {
			node.Value = i.Value
			node.Position = i.Position
			node.Children = i.Children
			transformTree(node, nil)
		}
	}
}

func in(char string, arr []string) bool {
	for _, i := range arr {
		if i == char {
			return true
		}
	}
	return false
}
