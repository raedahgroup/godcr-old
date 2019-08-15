package widgets

import (
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
	table.Container = widget.NewScrollContainer(nil)
	table.tableData = []*widget.Box{heading}
	table.tableData = append(table.tableData, data...)
	table.Result = widget.NewHBox()
	table.Refresh()
}

func (table *TableStruct) Refresh() {
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
