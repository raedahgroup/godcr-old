package widgets

const checkboxLabelHeight = 20

func (window *Window) AddCheckbox(text string, checked *bool, checkChanged func()) {
	window.Row(checkboxLabelHeight).Dynamic(1)
	// add space before label text so the box doesn't glue to the text
	if window.CheckboxText(" "+text, checked) && checkChanged != nil {
		checkChanged()
	}
}
