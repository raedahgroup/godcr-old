package primitives

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Form struct {
	*tview.Form

	// Keep record of all form items added to this form as *tview.Box objects
	formItems []*tview.Box

	// An optional capture function which receives a key event and returns the
	// event to be forwarded to the primitive's default input handler (nil if
	// nothing should be forwarded).
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
}

// NewForm returns a new form.
func NewForm() *Form {
	return &Form{
		Form: tview.NewForm(),
	}
}

// AddInputField adds an input field to the form. It has a label, an optional
// initial value, a field width (a value of 0 extends it as far as possible),
// an optional accept function to validate the item's value (set to nil to
// accept any text), and an (optional) callback function which is invoked when
// the input field's text has changed.
func (f *Form) AddInputField(label, value string, fieldWidth int, accept func(textToCheck string, lastChar rune) bool, changed func(text string)) *Form {
	f.AddFormItem(tview.NewInputField().
		SetLabel(label).
		SetText(value).
		SetFieldWidth(fieldWidth).
		SetAcceptanceFunc(accept).
		SetChangedFunc(changed))
	return f
}

// AddPasswordField adds a password field to the form. This is similar to an
// input field except that the user's input not shown. Instead, a "mask"
// character is displayed. The password field has a label, an optional initial
// value, a field width (a value of 0 extends it as far as possible), and an
// (optional) callback function which is invoked when the input field's text has
// changed.
func (f *Form) AddPasswordField(label, value string, fieldWidth int, mask rune, changed func(text string)) *Form {
	if mask == 0 {
		mask = '*'
	}
	f.AddFormItem(tview.NewInputField().
		SetLabel(label).
		SetText(value).
		SetFieldWidth(fieldWidth).
		SetMaskCharacter(mask).
		SetChangedFunc(changed))
	return f
}

func (f *Form) handlePasteEvents(inputField *tview.InputField) {
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlV {
			textToPaste, err := clipboard.ReadAll()
			if err == nil {
				inputField.SetText(textToPaste)
				return nil
			}
		}

		if f.inputCapture != nil {
			return f.inputCapture(event)
		}

		return event // allow this input field to handle key press event
	})
}

// AddDropDown adds a drop-down element to the form. It has a label, options,
// and an (optional) callback function which is invoked when an option was
// selected. The initial option may be a negative value to indicate that no
// option is currently selected.
func (f *Form) AddDropDown(label string, options []string, initialOption int, selected func(option string, optionIndex int)) *Form {
	f.AddFormItem(tview.NewDropDown().
		SetLabel(label).
		SetCurrentOption(initialOption).
		SetOptions(options, selected))
	return f
}

// AddCheckbox adds a checkbox to the form. It has a label, an initial state,
// and an (optional) callback function which is invoked when the state of the
// checkbox was changed by the user.
func (f *Form) AddCheckbox(label string, checked bool, changed func(checked bool)) *Form {
	f.AddFormItem(tview.NewCheckbox().
		SetLabel(label).
		SetChecked(checked).
		SetChangedFunc(changed))
	return f
}

// AddFormItem adds a new item to the form. This can be used to add your own
// objects to the form. Note, however, that the Form class will override some
// of its attributes to make it work in the form context. Specifically, these
// are:
//
//   - The label width
//   - The label color
//   - The background color
//   - The field text color
//   - The field background color
func (f *Form) AddFormItem(item tview.FormItem) *Form {
	var box *tview.Box

	if inputField, ok := item.(*tview.InputField); ok {
		box = inputField.Box
	} else if dropDown, ok := item.(*tview.DropDown); ok {
		box = dropDown.Box
	} else if checkBox, ok := item.(*tview.Checkbox); ok {
		box = checkBox.Box
	} else {
		// don't add this form item until we're sure what box element it is,
		// to prevent disparity between f.items and f.formItems
		// instead, add an input field with unknown element text to give some visual feedback
		return f.AddInputField(item.GetLabel(), "Unknown form item", item.GetFieldWidth(), nil, nil)
	}

	f.formItems = append(f.formItems, box)

	if inputField, ok := item.(*tview.InputField); ok {
		f.handlePasteEvents(inputField)
	} else {
		box.SetInputCapture(f.inputCapture)
	}

	f.Form.AddFormItem(item)
	return f
}

// AddButton adds a new button to the form. The "selected" function is called
// when the user selects this button. It may be nil.
func (f *Form) AddButton(label string, selected func()) *Form {
	f.Form.AddButton(label, selected)
	f.Form.GetButton(f.Form.GetButtonCount() - 1).SetInputCapture(f.inputCapture)
	return f
}

// Clear removes all input elements from the form, including the buttons if
// specified.
func (f *Form) Clear(includeButtons bool) *Form {
	f.formItems = nil
	f.Form.Clear(includeButtons)
	return f
}

// ClearFields clears texts from InputFields, unchecks Checkboxes and sets the selected index of DropDowns to 0
func (f *Form) ClearFields(includeButtons bool) *Form {
 for i := 0; i < f.GetFormItemsCount(); i++ {
  field := f.GetFormItem(i)
  if inputField, ok := field.(*tview.InputField); ok {
   inputField.SetText("")
  } else if dropDown, ok := field.(*tview.DropDown); ok {
   if selected, _ := dropDown.GetCurrentOption(); selected > 0 {
    dropDown.SetCurrentOption(0)
   }
  } else if checkBox, ok := field.(*tview.Checkbox); ok {
   checkBox.SetChecked(false)
  }
 }
 return f
}
// GetFormItemBox returns the form element at the given position as a *tview.Box object,
// starting with index 0. Elements are referenced in the order they were added.
// Buttons are not included.
func (f *Form) GetFormItemBox(index int) *tview.Box {
	return f.formItems[index]
}

func (f *Form) GetFormItemsCount() int {
	return len(f.formItems)
}

// SetInputCapture installs a function which captures key events on each form item and button
// before they are forwarded to the primitive's default key event handler. This function can
// then choose to forward that key event (or a different one) to the default
// handler by returning it. If nil is returned, the default handler will not
// be called.
//
// Providing a nil handler will remove a previously existing handler.
func (f *Form) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *Form {
	f.inputCapture = capture
	for _, box := range f.formItems {
		box.SetInputCapture(f.inputCapture)
	}
	return f
}

// RemoveFormItem removes the form element at the given position, starting with
// index 0. Elements are referenced in the order they were added. Buttons are
// not included.
func (f *Form) RemoveFormItem(index int) *Form {
	f.formItems = append(f.formItems[:index], f.formItems[index+1:]...)
	f.Form.RemoveFormItem(index)
	return f
}
