package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Структура для хранения информации о лексемах
type Token struct {
	Type   string
	Lexeme string
	Pos    [2]int
}

// Функция для проверки, находится ли позиция внутри кавычек
func isInQuotes(pos int, quotedMatches [][]int) bool {
	for _, qMatch := range quotedMatches {
		if pos >= qMatch[0] && pos <= qMatch[1] {
			return true
		}
	}
	return false
}

// Функция для анализа строки и выделения лексем
func analyze(text string) []Token {
	// Регулярные выражения для поиска элементов
	coutRegex := regexp.MustCompile(`\bcout\b`)
	endlRegex := regexp.MustCompile(`\bendl\b`)
	quotedRegex := regexp.MustCompile(`"[^"]*"`)
	variableRegex := regexp.MustCompile(`\b[_a-zA-Z][_a-zA-Z0-9]*\b`)
	operatorRegex := regexp.MustCompile(`<<|\+|-|/|\*|\%`)
	numberRegex := regexp.MustCompile(`\b\d+\b`)
	endRegex := regexp.MustCompile(`;`)

	// Поиск всех совпадений
	coutMatches := coutRegex.FindAllStringIndex(text, -1)
	endlMatches := endlRegex.FindAllStringIndex(text, -1)
	quotedMatches := quotedRegex.FindAllStringIndex(text, -1)
	variableMatches := variableRegex.FindAllStringIndex(text, -1)
	operatorMatches := operatorRegex.FindAllStringIndex(text, -1)
	numberMatches := numberRegex.FindAllStringIndex(text, -1)
	endMatches := endRegex.FindAllStringIndex(text, -1)

	// Создание списка лексем
	var tokens []Token

	// Добавление найденных лексем в список
	for _, match := range coutMatches {
		if !isInQuotes(match[0], quotedMatches) {
			tokens = append(tokens, Token{"KEY_WORD", text[match[0]:match[1]], [2]int{match[0], match[1]}})
		}
	}
	for _, match := range endlMatches {
		if !isInQuotes(match[0], quotedMatches) {
			tokens = append(tokens, Token{"ENDL_WORD", text[match[0]:match[1]], [2]int{match[0], match[1]}})
		}
	}
	for _, match := range quotedMatches {
		tokens = append(tokens, Token{"TEXT", text[match[0]:match[1]], [2]int{match[0], match[1]}})
	}
	for _, match := range variableMatches {
		word := text[match[0]:match[1]]
		if !coutRegex.MatchString(word) && !endlRegex.MatchString(word) {
			if !isInQuotes(match[0], quotedMatches) {
				tokens = append(tokens, Token{"VARIABLE", word, [2]int{match[0], match[1]}})
			}
		}
	}
	for _, match := range operatorMatches {
		if !isInQuotes(match[0], quotedMatches) {
			tokens = append(tokens, Token{"OPERATOR", text[match[0]:match[1]], [2]int{match[0], match[1]}})
		}
	}
	for _, match := range numberMatches {
		if !isInQuotes(match[0], quotedMatches) {
			tokens = append(tokens, Token{"NUMBER", text[match[0]:match[1]], [2]int{match[0], match[1]}})
		}
	}
	for _, match := range endMatches {

		tokens = append(tokens, Token{"SEMICOLON", text[match[0]:match[1]], [2]int{match[0], match[1]}})

	}

	// Сортировка всех лексем по их позициям
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].Pos[0] < tokens[j].Pos[0]
	})

	// Добавление неизвестных токенов
	var unknownTokens []Token
	lastEnd := 0
	for _, token := range tokens {
		if token.Pos[0] > lastEnd {
			unknownLexeme := text[lastEnd:token.Pos[0]]
			if strings.TrimSpace(unknownLexeme) != "" {
				unknownTokens = append(unknownTokens, Token{"UNKNOWN", unknownLexeme, [2]int{lastEnd, token.Pos[0]}})
			}
		}
		lastEnd = token.Pos[1]
	}
	if lastEnd < len(text) {
		unknownLexeme := text[lastEnd:]
		if strings.TrimSpace(unknownLexeme) != "" {
			unknownTokens = append(unknownTokens, Token{"UNKNOWN", unknownLexeme, [2]int{lastEnd, len(text)}})
		}
	}
	tokens = append(tokens, unknownTokens...)
	return tokens
}

func leksem(text string) string {

	// Анализ строки кода и получение лексем с их позициями
	tokens := analyze(text)
	var builder strings.Builder
	// Вывод списка лексем с их типами и позициями
	builder.WriteString(fmt.Sprintf("Лексема\t\tТип\t\tПозиция в тексте\n"))
	for _, token := range tokens {
		builder.WriteString(fmt.Sprintf("%s\t\t%s\t\t[%d,%d]\n", token.Lexeme, token.Type, token.Pos[0], token.Pos[1]-1))

	}
	return builder.String()
}
