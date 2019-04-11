package widgets

const tableRowHeight = 25

type Table struct {
	rows []*TableRow
}

type TableRow struct {
	cells []TableCell
	font fontFace
	isFontSet bool
}

func (row *TableRow) Render(window *Window) {
	for _, cell := range row.cells {
		if cell == nil {
			// need to fill this column with empty space so the next cell is added to the next column instead of this column
			window.Spacing(1)
		} else {
			cell.Render(window)
		}
	}
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

func (table *Table) AddRowWithFont(font fontFace, cells ...TableCell) {
	row := &TableRow{
		cells: cells,
		font: font,
		isFontSet: true,
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
		if row.isFontSet {
			window.UseFontAndResetToPrevious(row.font, func() {
				row.Render(window)
			})
		} else {
			row.Render(window)
		}
	}
}
