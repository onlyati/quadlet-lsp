package terminal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Table struct {
	Rows     [][][]string
	Style    TableStyle
	MaxWidth int

	HeadingSeparator       string
	RowSeparator           string
	TopDividerSeparator    string
	TopDivider             string
	MiddleDividerSeparator string
	MiddleDivider          string
	BottomDividerSeparator string
	BottomDivider          string
}

type TableStyle int

const (
	TopSeparatorTableStyle TableStyle = 1 << iota
	ColumnSeparatorTableStyle
	RowSeparatorTableStyle
	BottomSeparatorTableStyle
)

func NewTable(width int, headings ...string) *Table {
	if width == 0 {
		var err error
		if width, _, err = GetSize(); err != nil {
			width = -1
		}
	}

	self := Table{
		Style:    TopSeparatorTableStyle | ColumnSeparatorTableStyle | RowSeparatorTableStyle | BottomSeparatorTableStyle,
		MaxWidth: width,

		// https://en.wikipedia.org/wiki/Box-drawing_character#DOS
		HeadingSeparator:       " ",
		RowSeparator:           "│",
		TopDividerSeparator:    "╤",
		TopDivider:             "═",
		MiddleDividerSeparator: "┼",
		MiddleDivider:          "─",
		BottomDividerSeparator: "╧",
		BottomDivider:          "═",
	}

	self.Add(headings...)

	return &self
}

func (self *Table) ColumnWidths() (int, []int) {
	columns := len(self.Rows[0])
	columnWidths := make([]int, columns)

	for _, row := range self.Rows {
		for index, cell := range row {
			width := cellWidth(cell)
			if width > columnWidths[index] {
				columnWidths[index] = width
			}
		}
	}

	return columns, columnWidths
}

func (self *Table) Add(cells ...string) {
	if len(self.Rows) > 0 {
		columns := len(self.Rows[0])
		length := len(cells)
		if length != columns {
			panic(fmt.Sprintf("row has %d columns but must have %d", length, columns))
		}
	}

	self.Rows = append(self.Rows, splitCells(cells))
}

func splitCells(cells []string) [][]string {
	length := len(cells)
	splitCells := make([][]string, length)
	for index, cell := range cells {
		splitCells[index] = strings.Split(strings.Trim(cell, "\n"), "\n")
	}
	return splitCells
}

func (self *Table) Wrap() {
	if self.MaxWidth <= 0 {
		return
	}

	columns, columnWidths := self.ColumnWidths()

	cellsTotalWidth := 0
	for column := 0; column < columns; column++ {
		cellsTotalWidth += columnWidths[column]
	}
	totalWidth := cellsTotalWidth + columns - 1 // with dividers

	if totalWidth > self.MaxWidth {
		// Column shrink factor
		cellsMaxWidth := self.MaxWidth - (columns - 1)
		factor := float64(cellsMaxWidth) / float64(cellsTotalWidth)

		// Leftover
		realWidth := columns - 1 // just the dividers
		for _, columnWidth := range columnWidths {
			realWidth += int(float64(columnWidth) * factor)
		}
		leftover := self.MaxWidth - realWidth

		// New widths
		for column, columnWidth := range columnWidths {
			columnWidth = int(float64(columnWidth) * factor)

			// We'll apply the leftover from left to right
			if leftover > 0 {
				columnWidth += 1
				leftover -= 1
			}

			// Minimum column width
			if columnWidth < 1 {
				columnWidth = 1
			}

			columnWidths[column] = columnWidth
		}

		// Wrap rows
		for _, row := range self.Rows {
			for column, cell := range row {
				row[column] = wrap(cell, columnWidths[column])
			}
		}
	}
}

func (self *Table) Write(writer io.Writer, stylist *Stylist) {
	if len(self.Rows) <= 1 {
		// Empty table
		return
	}

	if stylist == nil {
		stylist = NewStylist(false)
	}

	self.Wrap()

	columns, columnWidths := self.ColumnWidths()
	rows := len(self.Rows)

	separator := func(divider string, dividerSeparator string) string {
		r := ""
		for column := 0; column < columns; column++ {
			r += strings.Repeat(divider, columnWidths[column])
			if column < columns-1 {
				r += dividerSeparator
			}
		}
		return r
	}

	rowSeparator := separator(self.MiddleDivider, self.MiddleDividerSeparator)
	bottomSeparator := separator(self.BottomDivider, self.BottomDividerSeparator)

	for r, row := range self.Rows {
		height := rowHeight(row)
		for line := 0; line < height; line++ {
			for column, cell := range row {
				if line < len(cell) {
					if r != 0 {
						fmt.Fprint(writer, stylist.Value(pad(cell[line], columnWidths[column])))
					} else {
						// Heading
						fmt.Fprint(writer, stylist.TypeName(pad(cell[line], columnWidths[column])))
					}
				} else {
					// Pad lines
					fmt.Fprint(writer, strings.Repeat(" ", columnWidths[column]))
				}

				if column < columns-1 {
					if (r != 0) && (self.Style&ColumnSeparatorTableStyle != 0) {
						fmt.Fprint(writer, self.RowSeparator)
					} else {
						fmt.Fprint(writer, " ")
					}
				}
			}

			fmt.Fprint(writer, "\n")
		}

		if r == 0 {
			// Heading
			if self.Style&TopSeparatorTableStyle != 0 {
				if self.Style&ColumnSeparatorTableStyle != 0 {
					fmt.Fprintln(writer, separator(self.TopDivider, self.TopDividerSeparator))
				} else {
					fmt.Fprintln(writer, separator(self.TopDivider, self.TopDivider))
				}
			}
		} else if r < rows-1 {
			// Middle row
			if self.Style&RowSeparatorTableStyle != 0 {
				fmt.Fprintln(writer, rowSeparator)
			}
		} else if self.Style&BottomSeparatorTableStyle != 0 {
			// Bottom row
			fmt.Fprintln(writer, bottomSeparator)
		}
	}
}

func (self *Table) Print() {
	self.Write(os.Stdout, StdoutStylist)
}

// Utils

func cellWidth(cell []string) int {
	width := 0
	for _, line := range cell {
		w := len(line)
		if w > width {
			width = w
		}
	}
	return width
}

func rowHeight(row [][]string) int {
	height := 0
	for _, cell := range row {
		h := len(cell)
		if h > height {
			height = h
		}
	}
	return height
}

func wrap(lines []string, width int) []string {
	var newLines []string
	for _, line := range lines {
		for {
			remainder := ""
			if len(line) > width {
				line, remainder = line[:width], line[width:]
			}
			newLines = append(newLines, line)
			line = remainder
			if line == "" {
				break
			}
		}
	}
	return newLines
}

func pad(s string, width int) string {
	return s + strings.Repeat(" ", width-len(s))
}
