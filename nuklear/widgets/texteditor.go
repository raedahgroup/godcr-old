package widgets

import "github.com/aarzilli/nucular"

const EditorHeight = 25
const defaultEditorWidth = 200

func (window *Window) AddEditors(textEditors ...*nucular.TextEditor) {
	widths := make([]int, len(textEditors))
	for i := range textEditors {
		widths[i] = defaultEditorWidth
	}
	window.AddEditorsWithWidths(widths, textEditors...)
}

func (window *Window) AddEditorsWithWidths(widths []int, textEditors ...*nucular.TextEditor) {
	if len(textEditors) == 0 || len(widths) == 0 {
		// don't add row that will never be populated
		return
	}
	window.Row(EditorHeight).Static(widths...)
	window.AddEditorsToCurrentRow(textEditors...)
}

func (window *Window) AddEditorsToCurrentRow(textEditors ...*nucular.TextEditor) {
	for _, editor := range textEditors {
		window.AddEditorToCurrentRow(editor)
	}
}

func (window *Window) AddEditorToCurrentRow(textEditor *nucular.TextEditor) {
	if textEditor == nil {
		// need to fill this column with empty space so the next cell is added to the next column instead of this column
		window.Spacing(1)
	} else {
		textEditor.Edit(window.Window)
	}
}
