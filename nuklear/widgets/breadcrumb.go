package widgets

type Breadcrumb struct {
	Text   string
	Action func(string, *Window)
}

func (window *Window) AddBreadcrumb(breadcrumbs []*Breadcrumb) {
	tableCells := []TableCell{}

	for index := range breadcrumbs {
		var cell TableCell
		isLastItem := breadcrumbs[index].Action == nil
		if isLastItem {
			cell = NewLabelTableCell(breadcrumbs[index].Text, "LC")
		} else {
			cell = NewLinkTableCell(breadcrumbs[index].Text, "", breadcrumbs[index].Action)
		}
		tableCells = append(tableCells, cell)
		if !isLastItem {
			tableCells = append(tableCells, NewLabelTableCell("/", "LC"))
		}
	}

	table := NewTable()
	table.AddRow(tableCells...)
	table.Render(window)
}
