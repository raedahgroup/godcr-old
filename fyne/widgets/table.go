package widgets

import (
	"fyne.io/fyne/widget"
)

type Table struct {
	tableData []*widget.Box
	heading   *widget.Box
	Result    *widget.Box
	Container *widget.ScrollContainer
}

func (table *Table) NewTable(heading *widget.Box, data ...*widget.Box) {
	table.heading = heading
	table.Container = widget.NewScrollContainer(nil)
	table.tableData = []*widget.Box{heading}
	table.tableData = append(table.tableData, data...)
	table.Result = widget.NewHBox()
	table.Refresh()
}

func (table *Table) Append(data ...*widget.Box) {
	table.tableData = append(table.tableData, data...)
	table.Refresh()
}

// Prepend is used to add to a stack
func (table *Table) Prepend(data ...*widget.Box) {
	table.tableData = append(data, table.tableData[1:]...)
	table.tableData = append([]*widget.Box{table.heading}, table.tableData...)
	table.Refresh()
}

// Delete method is used to delete object from stack. if tx notifier is created this remove the table from the stack thereby allowing call for for now we should just track transactions by comparing old with new
// Note: while using delete, consider heading
func (table *Table) Delete(tableNo int) {
	if len(table.tableData) < tableNo || tableNo >= len(table.tableData) {
		return
	}

	// cannot delete heading
	if tableNo == 0 {
		return
	}

	table.tableData = append(table.tableData[:tableNo], table.tableData[tableNo+1:]...)
	table.Refresh()
}

// Pop remove an object from the stack, Note it cant remove header
func (table *Table) Pop() {
	// not allowed to remove heading
	if len(table.tableData) <= 1 {
		return
	}
	table.tableData = table.tableData[:len(table.tableData)-1]
	table.Refresh()
}

func (table *Table) DeleteAll() {
	table.tableData = table.tableData[:1]
	table.Refresh()
}

func (table *Table) Refresh() {
	var container = widget.NewHBox()
	// get horizontals apart from heading
	for i := 0; i < len(table.heading.Children); i++ {
		// get vertical
		var getVerticals = widget.NewVBox()
		for _, data := range table.tableData {
			getVerticals.Append(data.Children[i])
		}
		container.Append(getVerticals)
	}

	table.Result.Children = container.Children
	table.Container.Content = widget.NewScrollContainer(table.Result).Content
	widget.Refresh(table.Result)
}