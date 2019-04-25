package widgets

type Breadcrumb struct {
	Text   string
	Action func(string, *Window)
}

func (window *Window) AddBreadcrumb(breadcrumb []*Breadcrumb) {
	tableCells := []TableCell{}

	for index := range breadcrumb {
		var cell TableCell
		isLastItem := index == len(breadcrumb)-1
		if isLastItem {
			cell = NewLabelTableCell(breadcrumb[index].Text, "LC")
		} else {
			cell = NewLinkTableCell(breadcrumb[index].Text, "", breadcrumb[index].Action)
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
