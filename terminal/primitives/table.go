package primitives

import "github.com/rivo/tview"

type Table struct {
	*tview.Table
}

func NewTable() *Table {
	return &Table{
		tview.NewTable(),
	}
}

func (table *Table) SetCellCenterAlign(row, column int, text string) *Table {
	cell := tview.NewTableCell(text).SetAlign(tview.AlignCenter)
	table.SetCell(row, column, cell)
	return table
}

func (table *Table) SetCellRightAlign(row, column int, text string) *Table {
	cell := tview.NewTableCell(text).SetAlign(tview.AlignRight)
	table.SetCell(row, column, cell)
	return table
}

func (table *Table) SetHeaderCell(row, column int, text string) *Table {
	headerCell := tview.NewTableCell(text).
		SetAlign(tview.AlignCenter).
		SetSelectable(false)
	table.SetCell(row, column, headerCell)
	return table
}
