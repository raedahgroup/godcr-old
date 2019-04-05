package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type Table struct {
	container *fyne.Container
}

func NewTable() *Table {
	return &Table{
		fyne.NewContainerWithLayout(layout.NewGridLayout(1)),
	}
}

func (table *Table) AddRow(objects ...fyne.CanvasObject) {
	row := fyne.NewContainerWithLayout(layout.NewGridLayout(len(objects)), objects...)
	table.container.AddObject(row)
}

func (table *Table) AddRowSimple(texts ...string) {
	tableCells := make([]fyne.CanvasObject, len(texts))
	for i, text := range texts {
		tableCells[i] = widget.NewLabel(text)
	}
	table.AddRow(tableCells...)
}

func (table *Table) Clear() {
	table.container.Objects = []fyne.CanvasObject{}
}

// DefaultTable returns a table that grows beyond the minimum size to cover all available space
func (table *Table) DefaultTable() *fyne.Container {
	return table.container
}

// CondensedTable returns a table that does not grow beyond the minimum size required to display the longest row
func (table *Table) CondensedTable() *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(table.container.MinSize()), table.container)
}
