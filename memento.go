package main

// TextEditor структура для хранения состояния текста для undo/redo
type TextEditor struct {
	history   []string // История состояний текста для Undo
	redo      []string // История состояний текста для Redo
	index     int      // Текущий индекс состояния в истории Undo
	redoIndex int      // Текущий индекс состояния в истории Redo
}

// Undo откатывает последнее изменение текста
func (t *TextEditor) Undo() {
	if t.index > 0 {
		t.redo = append(t.redo, t.history[t.index]) // Сохраняем текущее состояние в историю Redo
		t.index--
	}
}

// Redo возвращает последнее отмененное изменение текста
func (t *TextEditor) Redo() {
	if len(t.redo) > 0 && t.redoIndex < len(t.redo) {
		t.index++
		t.history = append(t.history, t.redo[t.redoIndex])
		t.redoIndex++
	}
}

// SetText устанавливает текст и добавляет его в историю Undo
func (t *TextEditor) SetChange(text string) {
	if t.index < len(t.history)-1 {
		t.history = t.history[:t.index+1]
	}
	t.history = append(t.history, text)
	//t.redo = nil // Очищаем историю Redo при новом изменении текста
	t.index = len(t.history) - 1
	//t.redoIndex = 0
}

// Text возвращает текущий текст
func (t *TextEditor) Text() string {
	return t.history[t.index]
}

// Clear очищает историю изменений текста
func (t *TextEditor) Clear() {
	t.history = nil
	t.redo = nil
	t.index = 0
	t.redoIndex = 0
}

/*func updateBuffer(buffer []string, newItems []string) []string {
	// Определяем, сколько элементов нужно удалить из буфера
	removeCount := len(newItems)

	// Если количество элементов в буфере меньше, чем количество новых элементов,
	// мы удаляем только старые элементы, равное разнице между длинами
	if len(buffer) < removeCount {
		removeCount = len(buffer)
	}

	// Сдвигаем элементы влево для удаления старых элементов
	copy(buffer[:len(buffer)-removeCount], buffer[removeCount:])

	// Заменяем старые элементы новыми элементами в конце буфера
	copy(buffer[len(buffer)-removeCount:], newItems)

	return buffer
}*/
