package codgen

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"kekemon.org/comp/parser"
)

var resultingCode string

//CodGen traverses AST to generate needed assembly code
func CodGen(decoast *parser.Node) string {
	fmt.Printf("\n[Generated code]\n\n")
	//display(decoast, 0, false)
	resultingCode += `.386
.model flat, stdcall
option casemap :none

include C:\masm32\include\kernel32.inc
include C:\masm32\include\masm32rt.inc

includelib C:\masm32\lib\kernel32.lib
includelib C:\masm32\lib\masm32rt.lib

.code
`

	isfunid = true
	generatecode(decoast)

	resultingCode += `
start:
	call main
	fn MessageBox,0,str$(eax),"Buhtiy",MB_OK
	invoke ExitProcess, 0
end start`

	//fmt.Println(names)
	return resultingCode
}

var tmp string    //temporary retrieved element
var tmppos string //its position

var typeid string //type storage
var id string     //id

var funcargs string

var code string //code

var funid string       //function identifier
var isfunid bool       //identifier of function
var isinitialized bool //indicates that variable is already initialized

var condExprCounterDepth int //counts how deep we inside nested conditions
var condExprCounter int      //counts how many not nested ternary operators are there

var relopCounter int = 1 //counts how many relational operations are there

var iterCounterDepth int //counts how deep we inside nested iterartion stateents
var iterCounter int      //counts how many not nested iteration statements is there

/*tracking ((offset)) from EBP (winapi stores its variables in places
that were told in tutorial, so i will just step from EBP a little further)*/
var offset int = 4    //offset to left	[<rest_of_stack>....esp....<local>....ebp....<args>....<something_before>]
var argOffset int = 8 //offset to right (first four bytes is for return adress)

//-------------------------------Stack------------------------------------
//first element on stack stores records of var_name:mem_address(offset from EBP) to
//to know where in stack variable is located
//for each new scope new table is being pushed on top of stack
type Stack struct {
	namesStack []map[string]int
	argStack   []map[string]int
}

//--------------------for_names----------------------
func (stack *Stack) push(new map[string]int) {
	stack.namesStack = append(stack.namesStack, new)
}

func (stack *Stack) pop() (map[string]int, error) {

	empty := make(map[string]int)
	length := len(stack.namesStack)
	if length == 0 {
		return empty, errors.New("Empty Stack")
	}

	res := stack.namesStack[length-1]
	stack.namesStack = stack.namesStack[:length-1]
	return res, nil
}

//--------------------for_args----------------------
func (stack *Stack) argPush(new map[string]int) {
	stack.argStack = append(stack.argStack, new)
}

func (stack *Stack) argPop() (map[string]int, error) {

	empty := make(map[string]int)
	length := len(stack.argStack)
	if length == 0 {
		return empty, errors.New("Empty Stack")
	}

	res := stack.argStack[length-1]
	stack.namesStack = stack.argStack[:length-1]
	return res, nil
}

var names Stack // used stack
//-------------------------------Stack------------------------------------

var countDeclaredParams int   //made to check if func called with right count of params
var countPassedParams bool    //indicates that following items is params and should be counted
var counterOfPassedParams int //will count passed params to erase them when function returns

var definedFunctions map[string]bool = make(map[string]bool)   // map of defined functions
var howManyParamsFuncHas map[string]int = make(map[string]int) // [funcName]howManyArgsItShouldRecieve
//made for detecting when wrong amount of params is passed

/*inspectNode puts everithyng unspecified to tmp variable
and than dependingly on val does something to it*/
func inspectNode(node *parser.Node) {

	//fmt.Println(node.Value)
	val := node.Value
	pos := node.Position

	if val == "identifier" {
		id = tmp
		tmp = ""
	} else if val == "funcIdentifier" {
		funid = tmp
		definedFunctions[funid] = true
		tmp = ""
	} else if val == "returnStmt" {
		if !isInteger(tmp) && names.namesStack[len(names.namesStack)-1][tmp] != 0 {
			code += "\n	pop	eax\n	mov	eax, dword ptr [ebp - " + strconv.Itoa(names.namesStack[len(names.namesStack)-1][tmp]) + "]\n"
		} else if isInteger(tmp) || names.namesStack[len(names.namesStack)-1][tmp] == 0 {
			code += "\n	pop	eax\n"
		}
		tmp = ""
	} else if val == "funDeclaration" {
		resultingCode += funid + ":\n	push	ebp\n	mov	ebp, esp\n	sub	esp, 256\n\n"
		resultingCode += code + "\n	mov	esp, ebp\n	pop	ebp\n	ret	0\n"
		isfunid = true

		code = ""
		offset = 4

		_, _ = names.argPop()
		argOffset = 8

		countDeclaredParams = 0
	} else if isInteger(val) {
		code += "	push    " + val + "\n"
		if countPassedParams {
			counterOfPassedParams++
		}
	} else if val == "un-" {
		code += "	pop	ebx\n	neg	ebx\n	push	ebx\n\n"

	} else if val == "/=" {
		code += "	mov	edx, 0\n	pop	eax\n	pop	ecx\n	idiv	ecx\n"

		found := false
		for i := len(names.namesStack) - 1; i >= 0; i-- {
			//fmt.Println(names.namesStack[i], "  ", val)
			if names.namesStack[i][tmp] != 0 {
				found = true
				code += "	mov	dword ptr [ebp - " + strconv.Itoa(names.namesStack[i][tmp]) + "], eax\n"
				break
			}
		}
		if !found {
			fmt.Printf("\nVariable %s is not declared at line : %s\n", tmp, tmppos)
			os.Exit(1)
		}
	} else if val == "/" {
		code += "	mov	edx, 0\n	pop	ecx\n	pop	eax\n	idiv	ecx\n	push	eax\n\n"
	} else if val == "*" {
		code += "	mov	edx, 0\n	pop	ecx\n	pop	eax\n	imul	ecx\n	push	eax\n\n"
	} else if val == "+" {
		code += "	mov	edx, 0\n	pop	ecx\n	pop	eax\n	add	eax, ecx\n	push	eax\n\n"
	} else if val == "-" {
		code += "	mov	edx, 0\n	pop	ecx\n	pop	eax\n	sub	eax, ecx\n	push	eax\n\n"
	} else if val == ">" { // ось тут може трохи збитись розуміння але я навиворіт зробив бо попаю в зворотньому порядку операнди
		code += "	pop	eax\n	pop	ecx\n	cmp	eax, ecx\n"
		code += "	jl	lemit" + strconv.Itoa(relopCounter) + "\n	push	0\n	jmp	outmit" + strconv.Itoa(relopCounter) + "\n"
		code += "lemit" + strconv.Itoa(relopCounter) + ":\n	push	1\noutmit" + strconv.Itoa(relopCounter) + ":\n\n"
		relopCounter++
	} else if val == "<" {
		code += "	pop	eax\n	pop	ecx\n	cmp	eax, ecx\n"
		code += "	jg	lemit" + strconv.Itoa(relopCounter) + "\n	push	0\n	jmp	outmit" + strconv.Itoa(relopCounter) + "\n"
		code += "lemit" + strconv.Itoa(relopCounter) + ":\n	push	1\noutmit" + strconv.Itoa(relopCounter) + ":\n\n"
		relopCounter++
	} else if val == ">=" {
		code += "	pop	eax\n	pop	ecx\n	cmp	eax, ecx\n"
		code += "	jle	lemit" + strconv.Itoa(relopCounter) + "\n	push	0\n	jmp	outmit" + strconv.Itoa(relopCounter) + "\n"
		code += "lemit" + strconv.Itoa(relopCounter) + ":\n	push	1\noutmit" + strconv.Itoa(relopCounter) + ":\n\n"
		relopCounter++
	} else if val == "<=" {
		code += "	pop	eax\n	pop	ecx\n	cmp	eax, ecx\n"
		code += "	jge	lemit" + strconv.Itoa(relopCounter) + "\n	push	0\n	jmp	outmit" + strconv.Itoa(relopCounter) + "\n"
		code += "lemit" + strconv.Itoa(relopCounter) + ":\n	push	1\noutmit" + strconv.Itoa(relopCounter) + ":\n\n"
		relopCounter++
	} else if val == "==" {
		code += "	pop	eax\n	pop	ecx\n	cmp	eax, ecx\n"
		code += "	je	lemit" + strconv.Itoa(relopCounter) + "\n	push	0\n	jmp	outmit" + strconv.Itoa(relopCounter) + "\n"
		code += "lemit" + strconv.Itoa(relopCounter) + ":\n	push	1\noutmit" + strconv.Itoa(relopCounter) + ":\n\n"
		relopCounter++
	} else if val == "!=" {
		code += "	pop	eax\n	pop	ecx\n	cmp	eax, ecx\n"
		code += "	jne	lemit" + strconv.Itoa(relopCounter) + "\n	push	0\n	jmp	outmit" + strconv.Itoa(relopCounter) + "\n"
		code += "lemit" + strconv.Itoa(relopCounter) + ":\n	push	1\noutmit" + strconv.Itoa(relopCounter) + ":\n\n"
		relopCounter++
	} else if val == "simpleExpression" {
		code += "	pop	eax\n	pop	ecx\n"
		code += "	or	eax, ecx\n	push	eax\n"
	} else if val == "typeSpecifier" {
		typeid = tmp
		tmp = ""
	} else if val == "varDeclId" {
		isinitialized = false
		if names.namesStack[len(names.namesStack)-1][tmp] != 0 {
			fmt.Printf("\nDouble declaration of %s\n", tmp)
			os.Exit(1)
		}
		id = tmp
	} else if val == "scopedVarDeclaration" {
		if !isinitialized {
			code += "	mov	dword ptr [ebp - " + strconv.Itoa(offset) + "], 0\n"
			names.namesStack[len(names.namesStack)-1][id] = offset
			offset += 4
			id = ""
		}
	} else if val == "varDeclInitialize" {
		isinitialized = true
		code += "	pop	eax\n	mov	dword ptr [ebp - " + strconv.Itoa(offset) + "], eax\n"
		names.namesStack[len(names.namesStack)-1][id] = offset
		offset += 4
	} else if val == "=" {
		found := false
		for i := len(names.namesStack) - 1; i >= 0; i-- {
			//fmt.Println(names.namesStack[i], "  ", val)
			if names.namesStack[i][tmp] != 0 {
				found = true
				code += "	pop	eax\n	pop	eax\n"
				code += "	mov	dword ptr [ebp - " + strconv.Itoa(names.namesStack[i][tmp]) + "], eax\n"
				break
			}
		}
		if !found {
			fmt.Printf("\nVariable %s is not declared at line : %s\n", tmp, tmppos)
			os.Exit(1)
		}

	} else if val == "condition" {
		condExprCounterDepth++
		code += "	pop	eax\n	cmp	eax, 0\n	jz	false" + strconv.Itoa(condExprCounter) + strconv.Itoa(condExprCounterDepth) + "\n"
	} else if val == "booltrue" {
		code += "	jmp	out" + strconv.Itoa(condExprCounter) + strconv.Itoa(condExprCounterDepth) + "\n	false" +
			strconv.Itoa(condExprCounter) + strconv.Itoa(condExprCounterDepth) + ":\n"
	} else if val == "conditionalExpression" {
		code += "	out" + strconv.Itoa(condExprCounter) + strconv.Itoa(condExprCounterDepth) + ":\n"
		condExprCounterDepth--
		if condExprCounterDepth == 0 {
			condExprCounter++
		}
	} else if val == "in" { //--------------------------inside compound---------------------------------
		newNames := make(map[string]int)
		names.push(newNames)
	} else if val == "compoundStmt" { //--------------------------inside compound---------------------------------
		_, _ = names.pop()
	} else if val == "incall" {
		countPassedParams = true
	} else if val == "args" {
		countPassedParams = false
	} else if val == "call" {
		if definedFunctions[tmp] == false {
			fmt.Println("Function called before declaration at line : ", pos)
			os.Exit(1)
		}
		if counterOfPassedParams != howManyParamsFuncHas[tmp] {
			fmt.Println("Function called with wrong amounts of args at line : ", pos)
			fmt.Println("Needed :", howManyParamsFuncHas[tmp])
			fmt.Println("Got :", counterOfPassedParams)
			os.Exit(1)
		}
		//"push eax" -- because we need result to appear on stack for next manipulations ALSO *4 conted params fo the right esp alignment
		code += "	call	" + tmp + "\n	add	esp, " + strconv.Itoa(counterOfPassedParams*4) + "\n	push	eax\n"
		counterOfPassedParams = 0
	} else if val == "param" {
		/*if !isinitialized {*/
		names.argStack[len(names.argStack)-1][tmp] = argOffset
		argOffset += 4
		countDeclaredParams++
		/*}*/
	} else if val == "infunc" {
		newArgNames := make(map[string]int)
		names.argPush(newArgNames)
	} else if val == "initeration" {
		iterCounterDepth++
		code += "iter" + strconv.Itoa(iterCounter) + strconv.Itoa(iterCounterDepth) + ":\n"
	} else if val == "iterationStmt" {
		code += "	pop	eax\n	cmp	eax, 1\n	je	iter" + strconv.Itoa(iterCounter) + strconv.Itoa(iterCounterDepth) +
			"\nbreakiter" + strconv.Itoa(iterCounter) + strconv.Itoa(iterCounterDepth) + ":\n"
		iterCounterDepth--
		if iterCounterDepth == 0 {
			iterCounter++
		}
	} else if val == "iterCondition" {
		code += "continueiter" + strconv.Itoa(iterCounter) + strconv.Itoa(iterCounterDepth) + ":\n"

	} else if val == "continue" {
		if iterCounterDepth == 0 {
			fmt.Println("Continue used outside iterationStatement at line : ", pos)
			os.Exit(1)
		}
		code += "	jmp	continueiter" + strconv.Itoa(iterCounter) + strconv.Itoa(iterCounterDepth) + "\n"
	} else if val == "break" {
		if iterCounterDepth == 0 {
			fmt.Println("Break used outside iterationStatement at line : ", pos)
			os.Exit(1)
		}
		code += "	jmp	breakiter" + strconv.Itoa(iterCounter) + strconv.Itoa(iterCounterDepth) + "\n"
	} else if val == "paramSepList" {
		howManyParamsFuncHas[funid] = countDeclaredParams
	} else /* if () {

	} else /* if () {

	} else /* if () {

	} else /* if () {

	} else /* if () {

	} else /* if () {

	} else */if !isInteger(val) && len(names.argStack) != 0 && names.argStack[len(names.argStack)-1][val] != 0 {
		code += "	push	dword ptr [ebp + " + strconv.Itoa(names.argStack[len(names.argStack)-1][val]) + "]\n"
	} else if !isInteger(val) && len(names.namesStack) != 0 /*&& names.namesStack[len(names.namesStack)-1][val] != 0 */ {
		for i := len(names.namesStack) - 1; i >= 0; i-- {
			//fmt.Println(names.namesStack[i], "  ", val)
			if names.namesStack[i][val] != 0 {
				code += "	push	dword ptr [ebp - " + strconv.Itoa(names.namesStack[i][val]) + "]\n"
				break
			}
		}
		if countPassedParams && val != "" && val != "argList" {
			counterOfPassedParams++
		}
		tmp = node.Value //оце гівнокод навіть по моїх мірках
		tmppos = node.Position
	} else /* if */ {
		tmp = node.Value
		tmppos = node.Position
	}
}

//From bottom to top, from left to right
func generatecode(node *parser.Node) {
	for _, i := range node.Children {
		generatecode(i)
	}
	inspectNode(node)
	//fmt.Println(node.Value)
}

//isInteger checks if string is an aniteger
func isInteger(token string) bool {
	integers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "."}
	isint := true
	if token == "" {
		return false
	}
	for _, i := range token {
		if !in(string(i), integers) {
			isint = false
		}
	}
	return isint
}

//in is supplementary function that i could import from other file
//but i wanted some independency between files
func in(char string, arr []string) bool {
	for _, i := range arr {
		if i == char {
			return true
		}
	}
	return false
}
