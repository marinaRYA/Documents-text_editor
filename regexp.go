package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"fyne.io/fyne/v2/widget"
)

// MobileNumberStorage хранит найденные мобильные номера
type PhoneNumberInfo struct {
	Number string
	Start  int
	End    int
}

type MobileNumberStorage struct {
	Numbers []PhoneNumberInfo
}

func find(editor *widget.Entry, result *widget.Label) *MobileNumberStorage {
	text := editor.Text

	delimiters := ",.;\n" // Разделители
	// Регулярное выражение для поиска российских мобильных номеров
	re := regexp.MustCompile(`(?:\+7|8)\s?\(?(?:900|9[1-6]\d|97[7-9]|98[0-9]|99[1-9])\)?(?:\s|\-)?\d{3}(?:\s|\-)?\d{2}(?:\s|\-)?\d{2}`)
	// Разделение текста по разделителям
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return strings.ContainsRune(delimiters, r)
	})
	storage := &MobileNumberStorage{}
	offset := 0
	for _, part := range parts {
		matches := re.FindAllStringIndex(part, -1)
		for _, match := range matches {
			start := match[0]
			end := match[1]
			foundExpression := part[start:end]
			storage.Numbers = append(storage.Numbers, PhoneNumberInfo{Number: foundExpression, Start: start + offset, End: end + offset})
		}
		offset += len(part) + len(delimiters)
	}

	// Формируем текст для вывода в Label
	var labelText strings.Builder
	for _, number := range storage.Numbers {
		labelText.WriteString(fmt.Sprintf("Номер: %s, Начало: %d\n", number.Number, number.Start))
	}
	result.SetText(labelText.String())

	return storage
}
func createLog(mobile *MobileNumberStorage, filename string) error {
	// Преобразуем данные о мобильных номерах в формат JSON
	jsonData, err := json.MarshalIndent(mobile, "", "    ")
	if err != nil {
		return err
	}

	// Записываем данные в файл
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
