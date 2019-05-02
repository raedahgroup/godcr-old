package widgets

type Breadcrumb struct {
	Text   string
	Action func()
}

func (b *Box) AddBreadcrumb(breadcrumb []*Breadcrumb) {
	tableCells := []TableCell{}

	for index := range breadcrumb {
		isLastItem := index == len(breadcrumb)-1

		tableCells = append(tableCells, NewLink(breadcrumb[index].Text, breadcrumb[index].Action))
		if !isLastItem {
			tableCells = append(tableCells, NewLink("/", nil))
		}
	}

	table := NewTable()
	table.AddRow(tableCells...)
	b.Add(table.CondensedTable())
}
