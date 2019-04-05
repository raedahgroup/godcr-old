package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type Table struct {
	*fyne.Container
}

func NewTable() *Table {
	return &Table{
		fyne.NewContainerWithLayout(layout.NewGridLayout(1)),
	}
}

func (table *Table) AddRow(objects ...fyne.CanvasObject) *Table {
	row := fyne.NewContainerWithLayout(layout.NewGridLayout(len(objects)), objects...)
	table.AddObject(row)
	return table
}

func (table *Table) AddRowSimple(texts ...string) *Table {
	tableCells := make([]fyne.CanvasObject, len(texts))
	for i, text := range texts {
		tableCells[i] = widget.NewLabel(text)
	}
	return table.AddRow(tableCells...)
}
