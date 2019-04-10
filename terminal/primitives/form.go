package primitives

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Form struct {
	*tview.Form

	// Keep record of all form items added to this form as *tview.Box objects
	boxItems []*tview.Box

	// The box's background color.
	// Shadows tview.Box private property.
	backgroundColor tcell.Color

	// If set to true, instead of position items and buttons from top to bottom,
	// they are positioned from left to right.
	// Shadows tview.Form private property.
	horizontal bool

	// Whether labels should be right aligned.
	rightAlignLabels bool

	// The alignment of the buttons.
	// Shadows tview.Form private property.
	buttonsAlign int

	// The number of empty rows between items.
	// Shadows tview.Form private property.
	itemPadding int

	// The label color.
	// Shadows tview.Form private property.
	labelColor tcell.Color

	// The background color of the input area.
	// Shadows tview.Form private property.
	fieldBackgroundColor tcell.Color

	// The text color of the input area.
	// Shadows tview.Form private property.
	fieldTextColor tcell.Color

	// The background color of the buttons.
	// Shadows tview.Form private property.
	buttonBackgroundColor tcell.Color

	// The color of the button text.
	// Shadows tview.Form private property.
	buttonTextColor tcell.Color

	// An optional capture function which receives a key event and returns the
	// event to be forwarded to the primitive's default input handler (nil if
	// nothing should be forwarded).
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
}

type FormItem interface {
	tview.FormItem

	// CalculateFieldSize is used to calculate and set the height and width needed to display the form item completely
	// given the max width that the item can occupy
	CalculateFieldSize(maxWidth int)

	// GetFieldHeight returns the height to use in drawing the form item.
	// Default height is 1.
	GetFieldHeight() int
}

// NewForm returns a new form.
func NewForm(rightAlignLabels bool) *Form {
	return &Form{
		Form:                  tview.NewForm(),
		backgroundColor:       tview.Styles.PrimitiveBackgroundColor,
		rightAlignLabels:      rightAlignLabels,
		itemPadding:           1,
		labelColor:            tview.Styles.SecondaryTextColor,
		fieldBackgroundColor:  tview.Styles.ContrastBackgroundColor,
		fieldTextColor:        tview.Styles.PrimaryTextColor,
		buttonBackgroundColor: tview.Styles.ContrastBackgroundColor,
		buttonTextColor:       tview.Styles.PrimaryTextColor,
	}
}

// SetBackgroundColor sets the box's background color.
func (f *Form) SetBackgroundColor(color tcell.Color) *Form {
	f.backgroundColor = color
	f.Form.SetBackgroundColor(color)
	return f
}

// SetLabelColor sets the color of the labels.
func (f *Form) SetLabelColor(color tcell.Color) *Form {
	f.labelColor = color
	f.Form.SetLabelColor(color)
	return f
}

// SetFieldBackgroundColor sets the background color of the input areas.
func (f *Form) SetFieldBackgroundColor(color tcell.Color) *Form {
	f.fieldBackgroundColor = color
	f.Form.SetFieldBackgroundColor(color)
	return f
}

// SetFieldTextColor sets the text color of the input areas.
func (f *Form) SetFieldTextColor(color tcell.Color) *Form {
	f.fieldTextColor = color
	f.Form.SetFieldTextColor(color)
	return f
}

// SetButtonsAlign sets how the buttons align horizontally, one of AlignLeft
// (the default), AlignCenter, and AlignRight. This is only
func (f *Form) SetButtonsAlign(align int) *Form {
	f.buttonsAlign = align
	f.Form.SetButtonsAlign(align)
	return f
}

// SetButtonBackgroundColor sets the background color of the buttons.
func (f *Form) SetButtonBackgroundColor(color tcell.Color) *Form {
	f.buttonBackgroundColor = color
	f.Form.SetButtonBackgroundColor(color)
	return f
}

// SetButtonTextColor sets the color of the button texts.
func (f *Form) SetButtonTextColor(color tcell.Color) *Form {
	f.buttonTextColor = color
	f.Form.SetButtonTextColor(color)
	return f
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
	} else if textView, ok := item.(*TextViewFormItem); ok {
		box = textView.GetTextView().Box
	} else {
		// don't add this form item until we're sure what box element it is,
		// to prevent disparity between f.items and f.boxItems
		// instead, add an input field with unknown element text to give some visual feedback
		return f.AddInputField(item.GetLabel(), "Unknown form item", item.GetFieldWidth(), nil, nil)
	}

	f.boxItems = append(f.boxItems, box)

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
	f.boxItems = nil
	f.Form.Clear(includeButtons)
	return f
}

// ClearFields clears texts from InputFields, unchecks Checkboxes and sets the selected index of DropDowns to 0
func (f *Form) ClearFields() *Form {
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
	return f.boxItems[index]
}

func (f *Form) GetFormItemsCount() int {
	return len(f.boxItems)
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
	for _, box := range f.boxItems {
		box.SetInputCapture(f.inputCapture)
	}
	return f
}

// RemoveFormItem removes the form element at the given position, starting with
// index 0. Elements are referenced in the order they were added. Buttons are
// not included.
func (f *Form) RemoveFormItem(index int) *Form {
	f.boxItems = append(f.boxItems[:index], f.boxItems[index+1:]...)
	f.Form.RemoveFormItem(index)
	return f
}

// SetHorizontal sets the direction the form elements are laid out. If set to
// true, instead of positioning them from top to bottom (the default), they are
// positioned from left to right, moving into the next row if there is not
// enough space.
func (f *Form) SetHorizontal(horizontal bool) *Form {
	f.horizontal = horizontal
	f.Form.SetHorizontal(horizontal)
	return f
}

// Draw draws this primitive onto the screen.
func (f *Form) Draw(screen tcell.Screen) {
	f.Box.Draw(screen)

	// Determine the dimensions.
	x, y, width, height := f.GetInnerRect()
	topLimit := y
	bottomLimit := y + height
	rightLimit := x + width
	startX := x

	formItems := make([]tview.FormItem, f.GetFormItemsCount())
	for i := range formItems {
		formItems[i] = f.GetFormItem(i)
	}

	formButtons := make([]*tview.Button, f.GetButtonCount())
	for i := range formButtons {
		formButtons[i] = f.GetButton(i)
	}

	// Find the longest label.
	var maxLabelWidth int
	for _, item := range formItems {
		labelWidth := tview.StringWidth(item.GetLabel())
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}

	if maxLabelWidth > 0 {
		maxLabelWidth++ // Add one space to separate label from field.
	}

	// Calculate positions of form items.
	positions := make([]struct{ x, y, width, height int }, len(formItems)+len(formButtons))
	var focusedPosition struct{ x, y, width, height int }
	for index, item := range formItems {
		// Calculate the space needed.
		itemHeight := 1
		if formItem, ok := item.(FormItem); ok {
			formItem.CalculateFieldSize(width)

			itemHeight = formItem.GetFieldHeight()
			if itemHeight <= 0 {
				itemHeight = 1
			}
		}

		labelWidth := tview.StringWidth(item.GetLabel())
		var itemWidth int
		if f.horizontal {
			fieldWidth := item.GetFieldWidth()
			if fieldWidth == 0 {
				fieldWidth = tview.DefaultFormFieldWidth
			}
			labelWidth++
			itemWidth = labelWidth + fieldWidth
		} else {
			if !f.rightAlignLabels {
				// We want all fields to align vertically on the left and have equal spacing to the right.
				labelWidth = maxLabelWidth
			} else if labelWidth > 0 {
				labelWidth++ // this item has a label and so should have a single space between the label and the field
			}
			itemWidth = width
		}

		// Advance to next line if there is no space.
		if f.horizontal && x+labelWidth+1 >= rightLimit {
			x = startX
			y += 2
		}

		// Adjust the item's attributes.
		if x+itemWidth >= rightLimit {
			itemWidth = rightLimit - x
		}
		item.SetFormAttributes(
			labelWidth,
			f.labelColor,
			f.backgroundColor,
			f.fieldTextColor,
			f.fieldBackgroundColor,
		)

		// Save position.
		positions[index].x = x
		positions[index].y = y
		positions[index].width = itemWidth
		positions[index].height = itemHeight
		if item.GetFocusable().HasFocus() {
			focusedPosition = positions[index]
		}

		// Advance to next item.
		if f.horizontal {
			x += itemWidth + f.itemPadding
		} else {
			y += itemHeight + f.itemPadding
		}
	}

	// How wide are the buttons?
	buttonWidths := make([]int, len(formButtons))
	buttonsWidth := 0
	for index, button := range formButtons {
		w := tview.StringWidth(button.GetLabel()) + 4
		buttonWidths[index] = w
		buttonsWidth += w + 1
	}
	buttonsWidth--

	// Where do we place them?
	if !f.horizontal && x+buttonsWidth < rightLimit {
		if f.buttonsAlign == tview.AlignRight {
			x = rightLimit - buttonsWidth
		} else if f.buttonsAlign == tview.AlignCenter {
			x = (x + rightLimit - buttonsWidth) / 2
		}

		// In vertical layouts, buttons always appear after an empty line.
		if f.itemPadding == 0 {
			y++
		}
	}

	// Calculate positions of buttons.
	for index, button := range formButtons {
		space := rightLimit - x
		buttonWidth := buttonWidths[index]
		if f.horizontal {
			if space < buttonWidth-4 {
				x = startX
				y += 2
				space = width
			}
		} else {
			if space < 1 {
				break // No space for this button anymore.
			}
		}
		if buttonWidth > space {
			buttonWidth = space
		}
		button.SetLabelColor(f.buttonTextColor).
			SetLabelColorActivated(f.buttonBackgroundColor).
			SetBackgroundColorActivated(f.buttonTextColor).
			SetBackgroundColor(f.buttonBackgroundColor)

		buttonIndex := index + len(formItems)
		positions[buttonIndex].x = x
		positions[buttonIndex].y = y
		positions[buttonIndex].width = buttonWidth
		positions[buttonIndex].height = 1

		if button.HasFocus() {
			focusedPosition = positions[buttonIndex]
		}

		x += buttonWidth + 1
	}

	// Determine vertical offset based on the position of the focused item.
	var offset int
	if focusedPosition.y+focusedPosition.height > bottomLimit {
		offset = focusedPosition.y + focusedPosition.height - bottomLimit
		if focusedPosition.y-offset < topLimit {
			offset = focusedPosition.y - topLimit
		}
	}

	// Draw items.
	for index, item := range formItems {
		// make the item start from somewhere closer to the right if labels are set to right align
		// only makes sense if form items are arranged vertically
		labelWidth := tview.StringWidth(item.GetLabel())
		if !f.horizontal && f.rightAlignLabels && labelWidth < maxLabelWidth {
			diff := maxLabelWidth - labelWidth
			positions[index].x = positions[index].x + diff
			positions[index].width = positions[index].width + diff
		}

		// Set position.
		y := positions[index].y - offset
		height := positions[index].height
		item.SetRect(positions[index].x, y, positions[index].width, height)

		// Is this item visible?
		if y+height <= topLimit || y >= bottomLimit {
			continue
		}

		// Draw items with focus last (in case of overlaps).
		if item.GetFocusable().HasFocus() {
			defer item.Draw(screen)
		} else {
			item.Draw(screen)
		}
	}

	// Draw buttons.
	for index, button := range formButtons {
		// Set position.
		buttonIndex := index + len(formItems)
		y := positions[buttonIndex].y - offset
		height := positions[buttonIndex].height
		button.SetRect(positions[buttonIndex].x, y, positions[buttonIndex].width, height)

		// Is this button visible?
		if y+height <= topLimit || y >= bottomLimit {
			continue
		}

		// Draw button.
		button.Draw(screen)
	}
}
