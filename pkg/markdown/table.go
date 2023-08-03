package markdown

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Table struct {
	header []any
	body   [][]any
	size   []int
}

func NewTable(header ...string) *Table {
	row := make([]any, 0, len(header))
	size := make([]int, 0, len(header))
	for _, v := range header {
		row = append(row, v)
		size = append(size, len(v))
	}
	return &Table{
		header: row,
		size:   size,
	}
}

func (t *Table) AddRow(row ...any) {
	body := make([]any, len(row))
	for i, row := range row {
		var item string
		switch row := row.(type) {
		case string:
			item = row
		case int:
			item = strconv.Itoa(row)
		default:
			item = fmt.Sprintf("%v", row)
		}
		body[i] = item
		if len(item) > t.size[i] {
			t.size[i] = len(item)
		}
	}
	t.body = append(t.body, body)
}

func (t *Table) Print(w io.Writer) error {
	// print header
	formatBuilder := new(strings.Builder)
	formatBuilder.WriteByte('|')
	for _, size := range t.size {
		fmt.Fprintf(formatBuilder, " %%-%ds |", size)
	}
	formatBuilder.WriteByte('\n')
	format := formatBuilder.String()
	if _, err := fmt.Fprintf(w, format, t.header...); err != nil {
		return err
	}

	// print separator
	if _, err := w.Write([]byte("|")); err != nil {
		return err
	}
	for _, size := range t.size {
		if _, err := fmt.Fprintf(w, "%s|", strings.Repeat("-", size+2)); err != nil {
			return err
		}
	}

	// print body
	for _, row := range t.body {
		if _, err := fmt.Fprintf(w, format, row...); err != nil {
			return err
		}
	}

	return nil
}
