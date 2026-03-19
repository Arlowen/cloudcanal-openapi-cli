package util

import "strings"

func FormatTable(headers []string, rows [][]string) string {
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i := range headers {
			value := safeCell(row, i)
			if len(value) > widths[i] {
				widths[i] = len(value)
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
		cells[i] = padRight(safeCell(values, i), widths[i])
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

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
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
