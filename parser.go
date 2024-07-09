package main

import (
	"fmt"
	"unicode"
)

type ParseError struct {
	Unccorect     string
	Start         int
	ExceptionDesc string
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isSign(ch rune) bool {
	return ch != '"'
}

func Parsing(input string, offset int) []ParseError {
	id := ""
	count := offset
	var errors []ParseError

	states := map[int]func(rune) int{
		1: func(ch rune) int {
			if isLetter(ch) {
				id += string(ch)
				if id == "cout" {
					id = ""
					return 2
				}
				return 1
			} else if unicode.IsSpace(ch) {
				return 1
			} else {
				errors = append(errors, reportError(input, count, "Пропущено ключевое слово 'cout'"))
				return 0
			}
		},
		2: func(ch rune) int {
			if unicode.IsSpace(ch) {
				return 2
			}
			if ch == '<' {
				id += string(ch)
				return 3
			}
			errors = append(errors, reportError(input, count, "Пропущен оператор '<<'"))

			return 0
		},
		3: func(ch rune) int {
			fmt.Println(string(ch))
			if (id == "<") && ch != '<' {
				errors = append(errors, reportError(input, count, "Пропущен оператор '<<'"))
				id = ""
				return 0
			} else if ch == '<' && (id == "<") {
				id += string(ch)
				if id == "<<" {
					id = ""
					return 3
				}
			}
			if ch == ' ' {
				return 3
			}
			if ch == '+' || ch == '-' {
				return 4
			}
			if isDigit(ch) {
				return 5
			}
			if ch == 'e' {
				id += "e"
				return 14
			}
			if ch == '"' {
				return 11
			}
			if isLetter(ch) {
				return 10
			}
			errors = append(errors, reportError(input, count, "Пропущено выражение после оператора"))
			return 0
		},
		4: func(ch rune) int {
			if unicode.IsSpace(ch) {
				return 4
			}
			if isLetter(ch) {
				return 10
			}
			if isDigit(ch) {
				return 5
			}
			errors = append(errors, reportError(input, count, "Пропущено выражение после оператора"))
			return 0
		},
		5: func(ch rune) int {
			if ch == '<' {
				id += "<"
				return 3
			}
			if isDigit(ch) {
				return 5
			}
			if ch == '.' {
				return 6
			}
			if ch == ';' {
				return 16
			}
			if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
				return 8
			}
			if unicode.IsSpace(ch) {
				return 9
			}
			errors = append(errors, reportError(input, count, "Неверный символ"))
			return 0
		},
		6: func(ch rune) int {
			if isDigit(ch) {
				return 7
			}
			errors = append(errors, reportError(input, count, "Пропущено число после точки"))
			return 0
		},
		7: func(ch rune) int {
			if isDigit(ch) {
				return 7
			}
			if ch == ';' {
				return 16
			}
			if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
				return 8
			}
			if unicode.IsSpace(ch) {
				return 9
			}
			errors = append(errors, reportError(input, count, "Неверная запись числа"))
			return 0
		},
		8: func(ch rune) int {

			if unicode.IsSpace(ch) {
				return 8
			}
			if ch == '+' || ch == '-' {
				return 4
			}
			if isLetter(ch) {
				return 10
			}
			if isDigit(ch) {
				return 5
			}

			errors = append(errors, reportError(input, count, "Пропущено выражение после оператора"))

			return 0
		},
		9: func(ch rune) int {
			if unicode.IsSpace(ch) {
				return 9
			}
			if ch == ';' {
				return 16
			}
			if ch == '<' {
				id += "<"
				return 3
			}
			if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
				return 8
			}
			if ch == ';' {
				return 16
			}
			errors = append(errors, reportError(input, count, "Неверный символ"))
			return 0
		},
		10: func(ch rune) int {
			if isLetter(ch) || isDigit(ch) {
				return 10
			}
			if unicode.IsSpace(ch) {
				return 9
			}
			if ch == ';' {
				return 16
			}
			if ch == '<' {
				id += "<"
				return 3
			}
			if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
				return 8
			}

			errors = append(errors, reportError(input, count, "Пропущен оператор <<"))
			return 0
		},
		11: func(ch rune) int {
			if isSign(ch) {
				return 12
			} else {
				return 13
			}
		},
		12: func(ch rune) int {
			if isSign(ch) {
				return 12
			}
			if ch == '"' {
				return 13
			}

			return 0
		},
		13: func(ch rune) int {
			if ch == '<' {
				id += "<"
				return 3
			}
			if ch == ';' {
				return 14
			}
			if unicode.IsSpace(ch) {
				return 13
			}
			errors = append(errors, reportError(input, count, "Пропущен оператор после текста"))
			return 0
		},
		14: func(ch rune) int {
			if ch == '<' {
				id += "<"
				return 3
			}
			if ch == ';' {
				return 16
			}
			if isLetter(ch) {
				id += string(ch)
				if id == "endl" {
					id = ""
					return 15
				} else if len(id) >= 4 {
					return 10
				}
				return 14
			}
			if isDigit(ch) {
				return 10
			}
			errors = append(errors, reportError(input, count, "Неверный символ"))
			return 0
		},
		15: func(ch rune) int {

			if ch == '<' {
				id += "<"
				return 3
			}
			if ch == ';' {
				return 16
			}
			if unicode.IsSpace(ch) {
				return 15
			}
			errors = append(errors, reportError(input, count, "Неверный символ"))
			return 0
		},
		16: func(ch rune) int {
			return 16
		},
		0: func(ch rune) int {
			/*fmt.Println(id)
			if id != "" && id != "<" {
				errors = append(errors, reportError(input, count, "Пропущено ключевое слово 'cout'"))
				return 1
			}*/
			if id != "<" {
				errors = append(errors, reportError(input, count, "Пропущен оператор '<<'"))
			}
			id = "<"
			if ch == '<' {
				id += "<"
				return 3
			}

			if ch == ';' {
				return 16
			}

			return 3
		},
	}
	state := 1

	runeInput := []rune(input)

	for i := 0; i < len(runeInput); i++ {
		nextState := states[state](runeInput[i])
		if nextState == 0 {
			if i != 0 {
				i--
			}
			runeInput = append(runeInput[:i], runeInput[i+1:]...)
		}
		state = nextState
		count++
		//fmt.Println(count)
		fmt.Println(state)
	}
	if state == 12 {
		errors = append(errors, reportError(input, count, "Пропущена '\"'"))
	}
	if state == 3 && (string(input[len(input)-2]) != " " || string(input[len(input)-1]) != " ") {
		errors = append(errors, reportError(input, count, "Пропущена выражение после оператора"))
	}
	// Проверяем, завершилась ли строка на состоянии 16
	if string(input[len(input)-1]) != ";" {
		errors = append(errors, reportError(input, count, "Пропущена ';' в конце"))
	}

	return errors

}

func reportError(input string, count int, message string) ParseError {
	return ParseError{
		Unccorect:     input,
		Start:         count,
		ExceptionDesc: message,
	}
}

/*func parser(input string) []ParseError {
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Enter the input string: ")
	//input, _ := reader.ReadString('\n')
	//fmt.Println(input)
	//input := `cout << endl;`
	var errors []ParseError
	offset := 0
	parts := strings.Split(input, ";")
	for i, part := range parts {
		if i == len(parts)-1 {
			if input[len(input)-1] == ';' && part != "" {
				part += ";"
			}
		} else if part != "" && part != "\n" {
			part += ";"
		}

		parts[i] = part
	}
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			errors := Parsing(part, offset)
			partErrors := Parsing(part, offset)
			errors = append(errors, partErrors...)
		}
		offset += len(part)
	}
	return errors
}*/
