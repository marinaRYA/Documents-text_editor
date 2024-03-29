package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type State int

const (
	Start State = iota
	CountryCode
	AreaCode
	ExchangeCode
	SubscriberNumber
	End
)

type PhoneNumberFSM struct {
	currentState State
}

func NewPhoneNumberFSM() *PhoneNumberFSM {
	return &PhoneNumberFSM{
		currentState: Start,
	}
}

func (fsm *PhoneNumberFSM) Process(input string, offset int) PhoneNumberInfo {
	var phoneNumber PhoneNumberInfo
	start := 0
	end := 0
	for i := 0; i < len(input); i++ {
		char := rune(input[i])
		switch fsm.currentState {
		case Start:
			if char == '8' {
				start = i
				fsm.currentState = CountryCode
			} else if char == '+' {
				if input[i+1] == '7' {
					start = i
					fsm.currentState = CountryCode
					i++
				} else {
					fsm.currentState = End
				}
			}

		case CountryCode:
			if (char == ' ' || char == '-' || char == '(') && input[i+1] == '9' {
				if char == '(' && input[i+4] != ')' {
					fsm.currentState = End
					break
				}
				continue
			} else {
				codeChar := string(char) + string(input[i+1]) + string(input[i+2])
				code, err := strconv.Atoi(codeChar)
				if err == nil && ((code >= 900 && code <= 969) || (code >= 977 && code <= 989) || (code >= 991 && code <= 999)) {
					i += 2
					fsm.currentState = AreaCode
				} else {
					fsm.currentState = End
				}
			}

		case AreaCode:
			if char == ' ' || char == '-' || char == ')' {
				continue
			} else if unicode.IsDigit(char) && unicode.IsDigit(rune(input[i+1])) && unicode.IsDigit(rune(input[i+2])) {
				i += 2
				fsm.currentState = ExchangeCode
			} else {
				fsm.currentState = End
			}

		case ExchangeCode:
			if char == ' ' || char == '-' {
				continue
			} else if unicode.IsDigit(char) && unicode.IsDigit(rune(input[i+1])) {
				i++
				fsm.currentState = SubscriberNumber
			} else {
				fsm.currentState = End
			}

		case SubscriberNumber:
			if char == ' ' || char == '-' {
				continue
			} else if unicode.IsDigit(char) && unicode.IsDigit(rune(input[i+1])) {
				end = i + 2
				i = len(input)
			} else {
				fsm.currentState = End
			}

		case End:
			fsm.currentState = Start
		}
	}

	if fsm.currentState == SubscriberNumber {
		phoneNumber = PhoneNumberInfo{
			Number: input[start:end],
			Start:  start + offset,
			End:    end + offset,
		}
	}

	return phoneNumber
}

func terminate(text string) []PhoneNumberInfo {
	fsm := NewPhoneNumberFSM()
	delimiters := ",.;\n"
	offset := 0
	inputs := strings.FieldsFunc(text, func(r rune) bool {
		return strings.ContainsRune(delimiters, r)
	})

	var results []PhoneNumberInfo
	for _, input := range inputs {
		phoneNumber := fsm.Process(input, offset)
		if phoneNumber.Number != "" {
			results = append(results, phoneNumber)
		}
		fsm.currentState = Start
		offset += len(input) + len(delimiters)
	}

	return results
}
func phoneNumberInfoToString(phoneNumbers []PhoneNumberInfo) string {
	var builder strings.Builder

	for _, phoneNumber := range phoneNumbers {
		builder.WriteString(fmt.Sprintf("Номер: %s, Начальная позиция: %d\nСостояния:\n Start-CountryCode-AreaCode-ExchangeCode-SubscriberNumber\n", phoneNumber.Number, phoneNumber.Start))
	}

	return builder.String()
}
