package widgets

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type TableStruct struct {
	tableData []*widget.Box
	heading   *widget.Box
	Result    *widget.Box
	Container *widget.ScrollContainer
}

func (table *TableStruct) NewTable(heading *widget.Box, data ...*widget.Box) {
	table.heading = heading
	table.Result = widget.NewHBox()
	table.Container = widget.NewScrollContainer(table.Result)
	table.tableData = []*widget.Box{heading}
	table.tableData = append(table.tableData, data...)
	table.set()
}

func (table *TableStruct) Append(data ...*widget.Box) {
	if len(table.tableData) == 0 {
		return
	}
	iTable := TableStruct{
		heading:   table.heading,
		tableData: data,
		Result:    widget.NewHBox(),
		Container: widget.NewScrollContainer(nil),
	}
	iTable.set()

	for i := 0; i < len(table.heading.Children); i++ {
		a, ok := interface{}(table.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		T, ok := interface{}(iTable.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		a.Children = append(a.Children, T.Children...)
		widget.Refresh(a)
	}
}

// Prepend adds to a table.
func (table *TableStruct) Prepend(data ...*widget.Box) {
	// this makes sure an heading is placed
	if len(table.Result.Children) == 0 || len(table.heading.Children) == 0 || len(table.tableData) == 0 {
		return
	}

	iTable := TableStruct{
		heading:   table.heading,
		tableData: data,
		Result:    widget.NewHBox(),
		Container: widget.NewScrollContainer(nil),
	}
	iTable.set()

	for i := 0; i < len(table.heading.Children); i++ {
		a, ok := interface{}(table.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		T, ok := interface{}(iTable.Result.Children[i]).(*widget.Box)
		if !ok {
			return
		}
		T.Children = append([]fyne.CanvasObject{a.Children[0]}, T.Children...)
		T.Children = append(T.Children, a.Children[1:]...)
		a.Children = T.Children
		widget.Refresh(a)
	}
	table.Container.Offset.Y = 1110
	widget.Refresh(table.Container)
}

// DeleteContent deletes all table data aside heading.
func (table *TableStruct) DeleteContent() {
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

func (table *TableStruct) NumberOfRows() (count int) {
	for i := 0; i < len(table.heading.Children); i++ {
		a, ok := interface{}(table.Result.Children[i]).(*widget.Box)
		if !ok {
			return 0
		}
		fmt.Println("Number of row is", len(a.Children))
		if count < len(a.Children) {
			count = len(a.Children)
		}

	}
	return count - 1
}

// Delete method deletes contents from the table.
// if deleting only one content then specify to as 0
func (table *TableStruct) Delete(from, to int) {
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
