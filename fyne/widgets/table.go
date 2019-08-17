package widgets

import (
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
		a.Append(T)
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
		T.Children = []fyne.CanvasObject{a.Children[0], T}
		T.Children = append(T.Children, a.Children[1:]...)
		a.Children = T.Children
		widget.Refresh(a)
	}
}

//Delete method is used to delete object from stack. if tx notifier is created this remove the table from the stack thereby allowing call for for now we should just track transactions by comparing old with new
//Note: while using delete, consider heading WIP, not needed in fyne.
func (table *TableStruct) Delete(tableNo int) {
	if len(table.tableData) < tableNo || tableNo >= len(table.tableData) {
		return
	}

	//cannot delete heading
	if tableNo == 0 {
		return
	}

	// table.tableData = append(table.tableData[:tableNo], table.tableData[tableNo+1:]...)
	// table.set()
}

func (table *TableStruct) Pop() {
	//not allowed to remove heading
	if len(table.tableData) <= 1 {
		return
	}
	// table.tableData = table.tableData[:len(table.tableData)-1]
	// table.set()
}

func (table *TableStruct) set() {
	var container = widget.NewHBox()

	//get horizontals apart from heading
	for i := 0; i < len(table.heading.Children); i++ {
		//get vertical
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
