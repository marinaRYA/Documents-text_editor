package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CtrlPlusShortcut struct{}

func (s *CtrlPlusShortcut) ShortcutName() string {
	return "CtrlPlusShortcut"
}

func (s *CtrlPlusShortcut) Key() fyne.KeyName {
	return fyne.KeyPlus
}

func (s *CtrlPlusShortcut) Modifier() fyne.KeyModifier {
	return fyne.KeyModifierControl
}

// Custom shortcut for Ctrl-
type CtrlMinusShortcut struct{}

func (s *CtrlMinusShortcut) ShortcutName() string {
	return "CtrlMinusShortcut"
}

func (s *CtrlMinusShortcut) Key() fyne.KeyName {
	return fyne.KeyMinus
}

func (s *CtrlMinusShortcut) Modifier() fyne.KeyModifier {
	return fyne.KeyModifierControl
}

var (
	filename       string = ""
	isFileModified bool   = false
	editorFontSize int    = 14
	//resultFontSize int    = 14
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Text Editor")
	w.Resize(fyne.NewSize(800, 600))
	changedtext := TextEditor{}

	titleLabel := widget.NewLabel("")
	// Создание текстового редактора
	editor := widget.NewMultiLineEntry()
	editor.Wrapping = fyne.TextWrapBreak
	editor.SetPlaceHolder("Откройте или создайте файл")
	editor.Disable()
	editor.OnChanged = func(t string) {
		isFileModified = true
		if changedtext.index != 0 {
			if changedtext.history[changedtext.index] != t {
				changedtext.SetChange(t)
			}
			return
		}
		changedtext.SetChange(t)

	}

	// Создание области для отображения результатов
	result := widget.NewLabel("Результат будет здесь")

	// Создание кнопок панели инструментов
	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		changedtext.Undo()
		if changedtext.index > 0 {
			editor.SetText(changedtext.Text())
		}
		fmt.Println("История Undo:", changedtext.history)
	})

	forwardButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		changedtext.Redo()
		if changedtext.index > 0 {
			editor.SetText(changedtext.Text())
		}

		fmt.Println("История Redo:", changedtext.redo)
	})

	saveButton := widget.NewButtonWithIcon("Сохранить", theme.DocumentSaveIcon(), func() {
		isFileModified = false
		saveToFile(filename, editor.Text, w)

	})
	clearButton := widget.NewButtonWithIcon("Очистить", theme.CancelIcon(), func() {
		editor.SetText("")
	})

	runButton := widget.NewButtonWithIcon("Запуск", theme.MediaPlayIcon(), func() {
		result.SetText("")
		input := editor.Text
		offset := 0
		if editor.Text != "" {
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
					if len(errors) == 0 {
						result.SetText(result.Text + fmt.Sprintf("Выражение верное. Позиция %d %d\n", offset, len(part)))
					} else {
						for _, err := range errors {
							result.SetText(result.Text + fmt.Sprintf("Error: %s в позиции %d\n ", err.ExceptionDesc, err.Start))

						}
					}
				}
				offset += len(part)
			}
		}
	})

	// Создание главного меню
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("Файл",
			fyne.NewMenuItem("Создать", func() {
				text := editor.Text
				if editor.Text != "" && isFileModified {
					confirmDialog := dialog.NewConfirm("Предупреждение", "Содержимое редактора будет заменено содержимым файла. Хотите сохранить текущее содержимое?", func(confirmed bool) {
						if confirmed {
							ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
								if err != nil {
									fyne.LogError("Ошибка при сохранении файла", err)
									return
								}

							}, w, text)

						}
					}, w)
					confirmDialog.SetDismissText("Нет")
					confirmDialog.SetConfirmText("Да")
					confirmDialog.Show()
					editor.SetText("")
				}
				createFile(w, &filename)
				titleLabel.SetText(filename)
				if filename != "" {
					editor.SetPlaceHolder("Введите ваш текст здесь...")
					editor.Enable()
				}
				changedtext.Clear()
				editor.SetText("")
			}),
			fyne.NewMenuItem("Открыть", func() {

				ShowFileOpen(func(data []byte, err error) {
					if err != nil {
						fyne.LogError("Ошибка при открытии файла", err)
						return
					}
					changedtext.Clear()
					editor.SetText(string(data))
					isFileModified = false
				}, w, editor, &filename, titleLabel)
				if editor.Text != "" && isFileModified {
					confirmDialog := dialog.NewConfirm("Предупреждение", "Содержимое редактора будет заменено содержимым файла. Хотите сохранить текущее содержимое?", func(confirmed bool) {
						if confirmed {
							ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
								if err != nil {
									fyne.LogError("Ошибка при сохранении файла", err)
									return
								}

							}, w, editor.Text)
						}
					}, w)
					confirmDialog.SetDismissText("Нет")
					confirmDialog.SetConfirmText("Да")
					confirmDialog.Show()
				}

			}),
			fyne.NewMenuItem("Сохранить как", func() {

				if !editor.Disabled() {
					ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
						if err != nil {
							fyne.LogError("Ошибка при сохранении файла", err)
							return
						}

					}, w, editor.Text)
				} else {
					dialog.ShowInformation("Предупреждение", "Откройте или создайте файл", w)
				}
			}),
		),
		fyne.NewMenu("Текст",
			fyne.NewMenuItem("Регулярные выражения", func() {
				if filename != "" {
					mobiledata := *find(editor, result)
					createLog(&mobiledata, "log")
				}
			}),
			fyne.NewMenuItem("ДКА", func() {
				phoneNumbers := terminate(editor.Text)
				output := phoneNumberInfoToString(phoneNumbers)
				result.SetText(output)
			}),
			fyne.NewMenuItem("Сканер", func() {
				tokens := leksem(editor.Text)
				result.SetText(tokens)
			}),
			fyne.NewMenuItem("Анализатор", func() {
				result.SetText("")
				input := editor.Text
				offset := 0
				if editor.Text != "" {
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
							if len(errors) == 0 {
								result.SetText(result.Text + fmt.Sprintf("Выражение верное. Позиция %d %d\n", offset, len(part)))
							} else {
								for _, err := range errors {
									result.SetText(result.Text + fmt.Sprintf("Error: %s в позиции %d\n ", err.ExceptionDesc, err.Start))

								}
							}
						}
						offset += len(part)
					}
				}
			}),
		),
		fyne.NewMenu("Редактировать",
			fyne.NewMenuItem("Назад", func() {
				changedtext.Undo()
				editor.SetText(changedtext.Text())
			}),
			fyne.NewMenuItem("Вперед", func() {
				changedtext.Redo()
				editor.SetText(changedtext.Text())
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Копировать", func() { editor.TypedShortcut(&fyne.ShortcutCopy{Clipboard: w.Clipboard()}) }),
			fyne.NewMenuItem("Вставить", func() { editor.TypedShortcut(&fyne.ShortcutPaste{Clipboard: w.Clipboard()}) }),
			fyne.NewMenuItem("Вырезать", func() { editor.TypedShortcut(&fyne.ShortcutCut{Clipboard: w.Clipboard()}) }),
			fyne.NewMenuItem("Удалить", func() {
				if editor.SelectedText() != "" {
					newText := strings.Replace(editor.Text, editor.SelectedText(), "", 1)
					editor.SetText(newText)
				}
			}),
			fyne.NewMenuItem("Выделить всё", func() { editor.TypedShortcut(&fyne.ShortcutSelectAll{}) }),
		),
		fyne.NewMenu("Помощь",
			fyne.NewMenuItem("Справка", func() {
				file := "help1.html"
				if _, err := os.Stat(file); os.IsNotExist(err) {

					return
				}
				err := openBrowser(file)
				if err != nil {
					println("Failed to open browser:", err)
				}
			}),
		),
	)
	w.SetMainMenu(mainMenu)

	w.SetContent(container.NewBorder(
		nil,
		nil,
		nil,
		container.NewVBox(
			titleLabel,
			container.NewHBox(
				backButton,
				forwardButton,
				saveButton,
				clearButton,
				runButton,
			),
			result,
		),
		editor,
	))
	updateFontSize := func(delta int) {
		editorFontSize += delta
		if editorFontSize < 10 {
			editorFontSize = 10
		} else if editorFontSize > 30 {
			editorFontSize = 30
		}

		editor.TextStyle = fyne.TextStyle{
			Bold:      true,
			Italic:    true,
			Monospace: true,
			Symbol:    true,
			TabWidth:  editorFontSize}

		w.Canvas().Content().Refresh()
		w.Canvas().Refresh(result)
	}
	result.Resize(fyne.NewSize(600, 600))
	w.Canvas().AddShortcut(&CtrlPlusShortcut{}, func(shortcut fyne.Shortcut) {
		updateFontSize(2) // Увеличиваем размер на 2
	})

	// Обработчик для Ctrl-
	w.Canvas().AddShortcut(&CtrlMinusShortcut{}, func(shortcut fyne.Shortcut) {
		updateFontSize(-2) // Уменьшаем размер на 2
	})

	w.ShowAndRun()
}
