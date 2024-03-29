package main

import (
	"fmt"
	"regexp"
)

type Token struct {
	TokenType string // Тип лексемы ("число", "идентификатор", "знак" и т.д.)
	Value     string // Значение лексемы (например, число или идентификатор)
	Position  int    // Позиция лексемы в тексте
}

func main() {
	text := "x = 10 + y * 5"

	tokens := Lexer(text)
	for _, token := range tokens {
		fmt.Printf("Тип: %s, Значение: %s, Позиция: %d\n", token.TokenType, token.Value, token.Position)
	}
}

func Lexer(text string) []Token {
	var tokens []Token

	// Регулярные выражения для поиска чисел, идентификаторов и знаков
	numberRegex := regexp.MustCompile(`\d+`)
	identifierRegex := regexp.MustCompile(`[a-zA-Z_]\w*`)
	operatorRegex := regexp.MustCompile(`[+\-*\/=]`)

	// Поиск и выделение чисел
	numberMatches := numberRegex.FindAllStringIndex(text, -1)
	for _, match := range numberMatches {
		tokens = append(tokens, Token{TokenType: "число", Value: text[match[0]:match[1]], Position: match[0]})
	}

	// Поиск и выделение идентификаторов
	identifierMatches := identifierRegex.FindAllStringIndex(text, -1)
	for _, match := range identifierMatches {
		tokens = append(tokens, Token{TokenType: "идентификатор", Value: text[match[0]:match[1]], Position: match[0]})
	}

	// Поиск и выделение знаков
	operatorMatches := operatorRegex.FindAllStringIndex(text, -1)
	for _, match := range operatorMatches {
		tokens = append(tokens, Token{TokenType: "знак", Value: text[match[0]:match[1]], Position: match[0]})
	}

	return tokens
}
