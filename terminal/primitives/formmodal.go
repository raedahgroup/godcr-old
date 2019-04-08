package primitives

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type FormModal struct {
	*tview.Box

	// The framed embedded in the modal.
	frame *tview.Frame

	// The form embedded in the modal's frame.
	form *Form
}

func NewFormModal(modalTitle string) *FormModal {
	m := &FormModal{
		Box:       tview.NewBox(),
	}

	m.form = NewForm()
	m.form.SetBackgroundColor(tview.Styles.ContrastBackgroundColor).SetBorderPadding(0, 0, 0, 0)
	m.form.SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonsAlign(tview.AlignCenter).
		SetItemPadding(0)

	m.frame = tview.NewFrame(m.form).SetBorders(0, 0, 0, 0, 0, 0)
	m.frame.SetBorder(true).
		SetTitle(modalTitle).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tview.Styles.PrimaryTextColor).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)

	return m
}

func (m *FormModal) AddFormItem(item tview.FormItem) *FormModal {
	m.form.AddFormItem(item)
	m.form.boxItems[len(m.form.boxItems) - 1].SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}
		return event
	})
	return m
}

func (m *FormModal) AddButton(label string, selected func()) *FormModal {
	m.form.AddButton(label, selected)
	button := m.form.GetButton(m.form.GetButtonCount() - 1)
	button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyRight:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp, tcell.KeyLeft:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}
		return event
	})
	return m
}

// Focus is called when this primitive receives focus.
func (m *FormModal) Focus(delegate func(p tview.Primitive)) {
	delegate(m.form)
}

// HasFocus returns whether or not this primitive has focus.
func (m *FormModal) HasFocus() bool {
	return m.form.HasFocus()
}

// Draw draws this primitive onto the screen.
func (m *FormModal) Draw(screen tcell.Screen) {
	// Calculate the width of this modal,
	// which is the greater of the widest form field or the total buttons width
	formWidth := 0

	// first get buttons width
	for i := 0; i < m.form.GetButtonCount(); i++ {
		buttonLabel := m.form.GetButton(i).GetLabel()
		formWidth += tview.StringWidth(buttonLabel) + 4 + 2
	}
	formWidth -= 2

	// get longest form item width
	for i := 0; i < m.form.GetFormItemsCount(); i++ {
		itemWidth := m.form.GetFormItem(i).GetFieldWidth()
		if itemWidth > formWidth {
			formWidth = itemWidth
		}
	}

	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3
	if width < formWidth {
		width = formWidth
	}

	// Set the modal's position and size.
	height := (m.form.GetFormItemsCount() * 2) + 5
	width += 5
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}

func (m *FormModal) SetCancelFunc(fn func()) *FormModal {
	m.form.SetCancelFunc(fn)
	return m
}
