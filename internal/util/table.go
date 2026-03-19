package util

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func FormatTable(headers []string, rows [][]string) string {
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = DisplayWidth(header)
	}
	for _, row := range rows {
		for i := range headers {
			value := safeCell(row, i)
			if DisplayWidth(value) > widths[i] {
				widths[i] = DisplayWidth(value)
			}
		}
	}

	var lines []string
	lines = append(lines, renderRow(headers, widths))
	lines = append(lines, renderSeparator(widths))
	for _, row := range rows {
		lines = append(lines, renderRow(row, widths))
	}
	return strings.Join(lines, "\n")
}

func renderRow(values []string, widths []int) string {
	cells := make([]string, len(widths))
	for i := range widths {
		cells[i] = PadDisplayRight(safeCell(values, i), widths[i])
	}
	return strings.Join(cells, " | ")
}

func renderSeparator(widths []int) string {
	cells := make([]string, len(widths))
	for i, width := range widths {
		cells[i] = strings.Repeat("-", width)
	}
	return strings.Join(cells, "-+-")
}

func DisplayWidth(value string) int {
	return runewidth.StringWidth(value)
}

func PadDisplayRight(value string, width int) string {
	if DisplayWidth(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-DisplayWidth(value))
}

func safeCell(values []string, index int) string {
	if index >= len(values) {
		return "-"
	}
	value := strings.TrimSpace(values[index])
	if value == "" {
		return "-"
	}
	return value
}
