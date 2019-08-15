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
	table.Refresh()
}

func (table *TableStruct) Append(data ...*widget.Box) {
	if len(table.tableData) == 0 {
		return
	}
	//fmt.Println(table.Result.Children[0])
	_, ok := interface{}(table.Result.Children[0]).(*widget.Box)
	fmt.Println("IS this correct", ok)

	//for i := range data {
	//data = append(, fyne.CanvasObject(data[i].Children))
	//}
	//var newData TableStruct
	//newData.NewTable(table.heading, data...)
	//newData.Result.Children.(fyne.CanvasObject) = table.Result.Children
	//table.Result.Children = append(table.Result.Children, newData.Result.Children...)
	//table.Result.Children = append(table.Result.Children, newData.Result.Children[1:]...) //Append(newData.Container.Content)
	//widget.Refresh(table.Result)
	//canvas.Refresh(table.Result)
	//widget.Refresh(table.Container)
	//widget.Refresh(table.Result)
}

//Prepend is used to add to a stack
func (table *TableStruct) Prepend(data ...*widget.Box) {
	// this makes sure an heading is placed
	if len(table.Result.Children) == 0 {
		return
	}
	newData := &TableStruct{
		tableData: data,
		heading:   table.heading,
	}
	newData.Refresh()

	newData.Result.Children = append([]fyne.CanvasObject{table.Result.Children[0]}, newData.Result.Children...)
	newData.Result.Children = append(newData.Result.Children, table.Result.Children[0:]...)

	// table.Result.Children = []fyne.CanvasObject{table.Result.Children[0]}
	// table.Result.Children = append(table.Result.Children, newData.Result)

	// newData.tableData = append(data, newData.tableData[1:]...)
	// newData.tableData = append([]*widget.Box{table.heading}, table.tableData...)
	// newData.Refresh()

	table.Result.Children = newData.Result.Children
	widget.Refresh(table.Result)
}

//Delete method is used to delete object from stack. if tx notifier is created this remove the table from the stack thereby allowing call for for now we should just track transactions by comparing old with new
//Note: while using delete, consider heading
func (table *TableStruct) Delete(tableNo int) {
	if len(table.tableData) < tableNo || tableNo >= len(table.tableData) {
		return
	}

	//cannot delete heading
	if tableNo == 0 {
		return
	}

	table.tableData = append(table.tableData[:tableNo], table.tableData[tableNo+1:]...)
	table.Refresh()
}

//Pop remove an object from the stack, Note it cant remove header
func (table *TableStruct) Pop() {
	//not allowed to remove heading
	if len(table.tableData) <= 1 {
		return
	}
	table.tableData = table.tableData[:len(table.tableData)-1]
	table.Refresh()
}

func (table *TableStruct) Refresh() {
	var container = widget.NewHBox()

	//get horizontals apart from heading
	for i := 0; i < len(table.heading.Children); i++ {
		//get vertical
		var getVerticals = widget.NewVBox()
		for _, data := range table.tableData {
			getVerticals.Append(data.Children[i])
		}
		container.Append(getVerticals)
	}

	table.Result.Children = container.Children
	widget.Refresh(table.Result)
}
