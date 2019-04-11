package widgets

const tableRowHeight = 25

type TableRow struct {
	cells []TableCell
}

type Table struct {
	rows []*TableRow
}

func NewTable() *Table {
	return &Table{}
}

func (table *Table) AddRow(cells ...TableCell) {
	row := &TableRow{
		cells: cells,
	}
	table.rows = append(table.rows, row)
}

func (table *Table) Render(window *Window) {
	// first calculate max column widths
	maxColumnWidths := make([]int, 0)
	for _, row := range table.rows {
		for i, cell := range row.cells {
			if cell == nil {
				continue
			}

			cellWidth := cell.MinWidth(window)
			if i >= len(maxColumnWidths) {
				maxColumnWidths = append(maxColumnWidths, cellWidth)
			} else if cellWidth > maxColumnWidths[i] {
				maxColumnWidths[i] = cellWidth
			}
		}
	}

	// create row constructor for each row of items and call draw on the items
	for _, row := range table.rows {
		window.Row(tableRowHeight).Static(maxColumnWidths...)
		for _, cell := range row.cells {
			if cell == nil {
				continue
			}
			cell.Render(window)
		}
	}
}
