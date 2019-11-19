package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type Table struct {
	tableData []*widget.Box
	heading   *widget.Box
	Result    *widget.Box
	Container *widget.ScrollContainer
}

// NewTable creates a new table widget
func (table *Table) NewTable(heading *widget.Box, data ...*widget.Box) {
	table.heading = heading
	table.Result = widget.NewHBox()
	table.Container = widget.NewScrollContainer(table.Result)
	table.tableData = []*widget.Box{heading}
	table.tableData = append(table.tableData, data...)
	table.set()
}

// Append widget adds to the bottom row of a table.
func (table *Table) Append(data ...*widget.Box) {
	if len(table.tableData) == 0 {
		return
	}
	iTable := Table{
		heading:   table.heading,
		tableData: data,
		Result:    widget.NewHBox(),
		Container: widget.NewScrollContainer(nil),
	}
	iTable.set()

	for i := 0; i < len(table.heading.Children); i++ {
		tableBox, ok := interface{}(table.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		iTableBox, ok := interface{}(iTable.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		tableBox.Children = append(tableBox.Children, iTableBox.Children...)
		widget.Refresh(tableBox)
	}
}

// Prepend adds to the top row of a table.
func (table *Table) Prepend(data ...*widget.Box) {
	// Makes sure an heading is placed
	if len(table.Result.Children) == 0 || len(table.heading.Children) == 0 || len(table.tableData) == 0 {
		return
	}

	iTable := Table{
		heading:   table.heading,
		tableData: data,
		Result:    widget.NewHBox(),
		Container: widget.NewScrollContainer(nil),
	}
	iTable.set()

	for i := 0; i < len(table.heading.Children); i++ {
		tableBox, ok := interface{}(table.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		iTableBox, ok := interface{}(iTable.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		iTableBox.Children = append([]fyne.CanvasObject{tableBox.Children[0]}, iTableBox.Children...)
		iTableBox.Children = append(iTableBox.Children, tableBox.Children[1:]...)
		tableBox.Children = iTableBox.Children
		widget.Refresh(tableBox)
	}
	table.Container.Offset.Y = 1110
	widget.Refresh(table.Container)
}

// DeleteContent deletes all table data aside heading.
func (table *Table) DeleteContent() {
	if len(table.Result.Children) == 0 || len(table.heading.Children) == 0 || len(table.tableData) == 0 {
		return
	}

	for i := 0; i < len(table.heading.Children); i++ {
		a, ok := interface{}(table.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		a.Children = []fyne.CanvasObject{a.Children[0]}
		widget.Refresh(a)
	}

}

// NumberOfColumns returns the number of columns in a table neglecting header count.
func (table *Table) NumberOfColumns() (count int) {
	for i := 0; i < len(table.heading.Children); i++ {
		a, ok := interface{}(table.Result.Children[i]).(*widget.Box)
		if !ok {
			return 0
		}
		if count < len(a.Children) {
			count = len(a.Children)
		}

	}
	return count - 1
}

// Delete method deletes contents from the table NEGLECTING header.
// If deleting only one content then specify to as 0.
func (table *Table) Delete(from, to int) {
	if len(table.Result.Children) == 0 || len(table.heading.Children) == 0 || len(table.tableData) == 0 {
		return
	}
	if from < 0 || to < 0 {
		return
	}
	if to != 0 {
		for i := 0; i < len(table.heading.Children); i++ {
			a, ok := interface{}(table.Result.Children[i]).(*widget.Box)
			if !ok {
				return
			}
			a.Children = append(a.Children[:from+1], a.Children[to+1:]...)
			widget.Refresh(a)
		}
	} else {
		for i := 0; i < len(table.heading.Children); i++ {
			a, ok := interface{}(table.Result.Children[i]).(*widget.Box)
			if !ok {
				return
			}
			a.Children = append(a.Children[:from+1], a.Children[from+2:]...)
			widget.Refresh(a)
		}
	}

	table.Container.Offset.Y = 1110
	widget.Refresh(table.Container)
}

func (table *Table) set() {
	var container = widget.NewHBox()

	// Get horizontals apart from heading.
	for i := 0; i < len(table.heading.Children); i++ {
		// Get vertical.
		var getVerticals = widget.NewVBox()
		for _, data := range table.tableData {
			if len(table.heading.Children) > len(data.Children) && i > len(data.Children)-1 {
				continue
			}
			getVerticals.Append(data.Children[i])
		}
		container.Append(getVerticals)
	}

	table.Result.Children = container.Children
	widget.Refresh(table.Result)
}

