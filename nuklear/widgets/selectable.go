package widgets

func (window *Window) AddCurrentNavButton(text string, clickFunc func()) {
	selected := true

	window.Row(bigButtonHeight).Dynamic(1)
	if window.SelectableLabel(text, CenterAlign, &selected) {
		clickFunc()
	}
}
