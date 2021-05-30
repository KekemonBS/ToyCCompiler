package lexer

import (
	"io/ioutil"
	"log"
	_ "reflect"
)

func Tokenize(path string) []string {

	var file, err = ioutil.ReadFile(path)
	var tokens []string

	if err != nil {
		log.Fatal("Wrong path")
	}

	//~ lexems := []string{"int", "float", "main", "return"}
	blanks := []string{" ", "\n", "\t"}
	//~ special_characters := []string{"(", ")", "{", "}", ";", "="}

	token := ""
	for i := 0; i < len(string(file)); i++ {
		char := string(file[i])
		if token == "" && in(char, blanks) {
			continue
		} else if in(char, blanks) && char != "" {
			tokens = append(tokens, token)
			token = ""
		//~ } else if умова {
		//	щось
		//~ } else if in(token, lexems) || in(token, special_characters) {
			//~ tokens = append(tokens, token)
			//~ token = char
		} else {
			token += char
		}
	}
	return tokens
}

func in(char string, arr []string) bool {
	for _, i := range arr {
		if i == char {
			return true
		}
	}
	return false
}
