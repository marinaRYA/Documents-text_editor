package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

func ShowFileSave(callback func(writer fyne.URIWriteCloser, err error), w fyne.Window, text string) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			callback(nil, err)
			return
		}

		if writer != nil {
			defer writer.Close()
			if !isTextFile(writer.URI().Name()) {
				dialog.ShowInformation("Ошибка", "Сохраненный файл не является текстовым файлом", w)
				return
			}

			_, err := writer.Write([]byte(text))
			if err != nil {
				dialog.ShowError(err, w)
				callback(nil, err)
				return
			}

			callback(writer, nil)
			infoDialog := dialog.NewInformation("Успех", "Файл успешно сохранен", w)
			infoDialog.Show()
		}
	}, w)
}
func ShowFileOpen(callback func(data []byte, err error), w fyne.Window, editor *widget.Entry, name *string, t *widget.Label) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {

		if err != nil {
			dialog.ShowError(err, w)
			callback(nil, err)
			return
		}

		if reader == nil {
			callback(nil, nil)
			return
		}
		defer reader.Close()
		data, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, w)
			callback(nil, err)
			return
		}
		if !isTextFile(reader.URI().String()) {
			dialog.ShowInformation("Ошибка", "Выбранный файл не является текстовым файлом", w)
			callback(nil, errors.New("выбранный файл не является текстовым файлом"))
			return
		}

		*name = reader.URI().Path() //название файла
		t.SetText(*name)
		callback(data, nil)
		// Устанавливаем содержимое файла в текстовое поле
		editor.SetText(string(data))
		editor.SetPlaceHolder("Введите ваш текст здесь...")
		editor.Enable()
	}, w)
}
func createFile(w fyne.Window, name *string) {
	// Исходное имя файла
	baseName := "new.txt"

	// Проверяем, существует ли файл с исходным именем
	if _, err := os.Stat(baseName); err == nil {
		// Файл существует
		// Извлекаем расширение файла
		ext := ""
		lastDot := strings.LastIndex(baseName, ".")
		if lastDot != -1 {
			ext = baseName[lastDot:]
		}

		// Генерируем новое имя файла с добавлением числового суффикса
		for i := 1; ; i++ {
			newName := fmt.Sprintf("new%d%s", i, ext)
			if _, err := os.Stat(newName); os.IsNotExist(err) {
				// Файл с таким именем не существует, используем его
				baseName = newName
				break
			}
		}
	}

	// Создаем файл
	file, err := os.Create(baseName)
	if err != nil {
		dialog.ShowInformation("Ошибка", "Ошибка при создании файла: "+err.Error(), w)
		return
	}
	defer file.Close()

	// Устанавливаем имя файла
	*name = baseName
}
func saveToFile(fileName string, text string, w fyne.Window) {
	// Открытие файла для записи
	file, err := os.Create(fileName)
	if err != nil {
		dialog.ShowInformation("Ошибка", "Ошибка при сохранении в файл: "+err.Error(), w)
		return
	}
	defer file.Close()

	// Запись текста в файл
	_, err = file.WriteString(text)
	if err != nil {
		dialog.ShowInformation("Ошибка", "Ошибка при записи в файл: "+err.Error(), w)
		return
	}

}
func isTextFile(filename string) bool {
	textFileExtensions := []string{".txt", ".csv", ".doc", ".rtf"}
	for _, ext := range textFileExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}
