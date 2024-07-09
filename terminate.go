package main

import (
	"fmt"

	//"strconv"
	"strings"
	"unicode"
)

type PhoneNumberInfoDKA struct {
	Number    string
	Start     int
	End       int
	Statepath string
}

func Process(input string, offset int) PhoneNumberInfoDKA {
	var phoneNumber PhoneNumberInfoDKA
	count := 0
	c := 0
	start := 0
	end := 0
	states := map[string]func(rune) string{
		"Start": func(ch rune) string {
			switch {
			case ch == '+':
				return "S0"
			case ch == '8':
				return "CC"
			default:
				return "Start"
			}
		},
		"S0": func(ch rune) string {
			switch {
			case ch == '7':
				return "CC"
			default:
				return "Start"
			}
		},
		"CC": func(ch rune) string {
			switch {
			case ch == '(':
				return "S3"
			case ch == ' ':
				return "S7"
			case ch == '9':
				return "S1"
			default:
				return "Start"
			}
		},
		"S1": func(ch rune) string {

			if unicode.IsDigit(ch) {
				c++
				if c == 9 {
					return "End"
				}
				return "S1"
			}
			return "Start"
		},
		"S3": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S4"
			}
			return "Start"
		},
		"S4": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S5"
			}
			return "Start"
		},
		"S5": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S6"
			}
			return "Start"
		},
		"S6": func(ch rune) string {
			if ch == ')' {
				return "AC"
			}
			return "Start"
		},
		"S7": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S8"
			}
			return "Start"
		},
		"S8": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S9"
			}
			return "Start"
		},
		"S9": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S10"
			}
			return "Start"
		},
		"S10": func(ch rune) string {
			if ch == ' ' {
				return "AC"
			}
			return "Start"
		},
		"AC": func(ch rune) string {

			if unicode.IsDigit(ch) {
				count++
				if count == 3 {
					return "EC"
				}
				return "AC"
			}
			return "Start"
		},
		"EC": func(ch rune) string {
			switch {
			case ch == '-':
				return "S11"
			case ch == ' ':
				return "S16"
			default:
				return "Start"
			}
		},
		"S11": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S12"
			}
			return "Start"
		},
		"S12": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S13"
			}
			return "Start"
		},
		"S13": func(ch rune) string {
			if ch == '-' {
				return "S14"
			}
			return "Start"
		},
		"S14": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S15"
			}
			return "Start"
		},
		"S15": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "End"
			}
			return "Start"
		},
		"S16": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S17"
			}
			return "Start"
		},
		"S17": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S18"
			}
			return "Start"
		},
		"S18": func(ch rune) string {
			if ch == ' ' {
				return "S19"
			}
			return "Start"
		},
		"S19": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "S20"
			}
			return "Start"
		},
		"S20": func(ch rune) string {
			if unicode.IsDigit(ch) {
				return "End"
			}
			return "Start"
		},
	}

	// Начальное состояние
	currentState := "Start"

	var statepath strings.Builder
	// Обработка входной строки
	for i, ch := range input {
		if nextState, exists := states[currentState]; exists {
			if currentState == "Start" {
				start = i
			}
			statepath.WriteString(fmt.Sprintf("%s-", currentState))
			currentState = nextState(ch)
			// Обновляем конечный индекс, если достигнуто конечное состояние или конец строки
			if currentState == "End" {
				statepath.WriteString(fmt.Sprintf("%s", currentState))
				end = i + 1
				phoneNumber = PhoneNumberInfoDKA{
					Number:    input[start:end],
					Start:     start + offset,
					End:       end + offset,
					Statepath: statepath.String(),
				}
				return phoneNumber
			}
		}
	}
	return phoneNumber

}

func terminate(text string) []PhoneNumberInfoDKA {
	delimiters := ",.;\n"
	offset := 0
	inputs := strings.FieldsFunc(text, func(r rune) bool {
		return strings.ContainsRune(delimiters, r)
	})

	var results []PhoneNumberInfoDKA
	for _, input := range inputs {
		phoneNumber := Process(input, offset)
		if phoneNumber.Number != "" {
			results = append(results, phoneNumber)
		}
		offset += len(input)
	}

	return results
}

func phoneNumberInfoToString(phoneNumbers []PhoneNumberInfoDKA) string {
	var builder strings.Builder
	for _, phoneNumber := range phoneNumbers {
		builder.WriteString(fmt.Sprintf("Номер: %s, Начальная позиция: %d\nСостояния:%s\n ", phoneNumber.Number, phoneNumber.Start, phoneNumber.Statepath))
	}
	return builder.String()
}
