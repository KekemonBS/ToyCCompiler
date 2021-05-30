package lexer

import (
	_ "fmt" //used to debug sometimes
	"io/ioutil"
	"log"
	_ "reflect" //i dont know why i imported that
	"strconv"
)

/*Tokenize splits code into tokens that can be analyzed by Parser*/
func Tokenize(path string) ([]string, []string, []string) {

	var file, err = ioutil.ReadFile(path)
	var tokens []string
	var description []string

	var line int = 1
	//var posInLine int
	var positions []string //line for number of line and posInLine
	//for position in line

	if err != nil {
		log.Fatal("\n  Wrong path")
	}

	keywords := []string{"int", "float", "return"}
	blanks := []string{" ", "\n", "\t"}
	delimiters := []string{"+", "-", "*", "/", ",", ";", "%",
		">", "<", "=", "(", ")", "[", "]", "{", "}", "!", "~", "?"} //, "!" not sure

	//~ operators := []string{"+", "-", "*", "/", "<", ">", "="}
	//~ identifier := []string{}

	previous := ""
	token := ""
	next := ""
	ignoreLine := false
	inString := false
	for i := 0; i < len(string(file)); i++ {
		//posInLine++ //додає до позиції в лінії
		//кожен нюлайн в "Рахує номер рядка" скидає цю змінну

		//Вибирає попередній, поточний і кінцевий символ
		if i > 1 {
			previous = string(file[i-1])
		}

		char := string(file[i])

		if i < len(string(file))-1 {
			next = string(file[i+1])
		}

		//for "\r\n" handling
		if char == "\r" && next == "\n" {
			char = "\n"
		} else if char == "\n" && previous == "\r" {
			continue
		}
		
		//Пропускає коментарі
		if char == "/" && next == "/" {
			ignoreLine = true
		} else if ignoreLine && char == "\n" {
			ignoreLine = false
		}
		//Формує рядки
		if !ignoreLine && !inString && char == "\"" {
			inString = true
			tokens = append(tokens, char)
			positions = append(positions, strconv.Itoa(line))
			continue
		} else if inString && char == "\"" && previous != "\\" {
			inString = false
			tokens = append(tokens, token, char)
			token = ""
			continue
		}
		if !inString {
			if !ignoreLine {
				//Пропускає лишні пробіли
				if previous != "" && in(char, blanks) && token != "" {
					tokens = append(tokens, token)
					positions = append(positions, strconv.Itoa(line))

					token = ""
				}
				if token == "" && in(char, blanks) {
					//Рахує номер рядка
					if char == "\n" {
						line += 1
					}
					continue

					//Виділяє ключові слова
				} else if in(token, keywords) && in(next, blanks) {
					tokens = append(tokens, token)
					positions = append(positions, strconv.Itoa(line))
					token = ""
					tokens = append(tokens, char)
					positions = append(positions, strconv.Itoa(line))
					//Виділяє роздільники
				} else if in(char, delimiters) {
					if token != "" {
						tokens = append(tokens, token)
						positions = append(positions, strconv.Itoa(line))

					}
					tokens = append(tokens, char)
					positions = append(positions, strconv.Itoa(line))
					token = ""
				} else {
					token += char
				}
			}
		} else {
			token += char
		}
	}
	description = analyze(tokens)
	return tokens, description, positions
}

func in(char string, arr []string) bool {
	for _, i := range arr {
		if i == char {
			return true
		}
	}
	return false
}

func inStr(char string, arr string) bool {
	for _, i := range arr {
		if string(i) == char {
			return true
		}
	}
	return false
}

func analyze(tokens []string) []string {
	var analyzed []string

	keywordsEquiv := map[string]string{"int": "TYPE", "float": "TYPE", "return": "RET"}

	delimitersEquiv := map[string]string{"{": "OPNBRC", "}": "CLSBRC", "[": "OPNBRKT", "]": "CLSBRKT", "(": "OPNPRTH",
		")": "CLSPRTH", ",": "COMMA", ";": "SEMCOL", "%": "PERCENT", "+": "ADD",
		"-": "SUB", "*": "MUL", "!": "NOT", "/": "DIV", "~": "BITCOMPL", ">": "GR",
		"<": "LS", "=": "EQ", "?": "QM", ":": "COL", "break": "BR", "continue": "CONT", "for": "FOR", "while": "WHL",
		"do": "DO", "||": "OR"} // BITCOMPL = flip all bits, QM = question mark

	opnq := false
	for ind, i := range tokens {

		//~ fmt.Printf("\n%d -- %d", ind, ind-1)
		if keywordsEquiv[i] != "" {
			analyzed = append(analyzed, keywordsEquiv[i])
		} else if delimitersEquiv[i] != "" {
			analyzed = append(analyzed, delimitersEquiv[i])
		} else if i == "\"" && opnq == false {
			analyzed = append(analyzed, "OPNQ")
			opnq = true
		} else if ind-1 >= 0 && analyzed[ind-1] == "OPNQ" {
			analyzed = append(analyzed, "STR")
		} else if tokens[ind] == "\"" && opnq == true {
			analyzed = append(analyzed, "CLSQ")
			opnq = false
		} else if isInteger(i) {
			analyzed = append(analyzed, whichInteger(i))
		} else /*if  {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } else /*if () {

		  } */{
			analyzed = append(analyzed, "ID")
		}
	}
	return analyzed
}

//Finnaly its done

//~ // Returns 'true' if the character is a DELIMITER.
//~ isDelimiter()

//~ // Returns 'true' if the character is an OPERATOR.
//~ isOperator()

//~ // Returns 'true' if the string is a VALID IDENTIFIER.
//~ validIdentifier()

//~ // Returns 'true' if the string is a KEYWORD.
//~ isKeyword()

// Returns 'true' if the string is an INTEGER.
func isInteger(token string) bool {
	hex := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "a", "b", "c", "d", "e", "f", "x"}
	integers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "."}
	oct := []string{"0", "1", "2", "3", "4", "5", "6", "7"}
	isint := true
	for _, i := range token {
		if len(token) < 2 && in(string(i), integers) {
			continue
		} else if in(string(i), integers) && token[0] != '0' && token[1] != 'x' {
			continue
		} else if token[0] == '0' && in(string(i), oct) {
			continue
		} else if token[0] == '0' && token[1] == 'x' && in(string(i), hex) {
			continue
		} else {
			isint = false
		}
	}
	return isint
}

func whichInteger(number string) string {
	result := "none"
	if len(number) > 1 && number[0] == '0' && number[1] == 'x' {
		result = "HEX"
	} else if number[0] == '0' && len(number) > 1 {
		result = "OCT"
	} else if inStr(".", number) {
		result = "FLOAT"
	} else {
		result = "INT"
	}
	return result
}

//~ // Returns 'true' if the string is a REAL NUMBER.
//~ isRealNumber()

//~ // Extracts the SUBSTRING.
//~ subString()
