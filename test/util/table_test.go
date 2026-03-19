package util_test

import (
	"cloudcanal-openapi-cli/internal/util"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestFormatTableAlignsMixedWidthContent(t *testing.T) {
	output := util.FormatTable(
		[]string{"ID", "名称", "健康度"},
		[][]string{
			{"1", "worker71kah6ac07o", "Health"},
			{"2", "同步任务A", "健康"},
		},
	)

	assertPipeAlignment(t, output)

	for _, want := range []string{"名称", "同步任务A", "worker71kah6ac07o"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q in %q", want, output)
		}
	}
}

func assertPipeAlignment(t *testing.T, output string) {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var baseline []int
	for _, line := range lines {
		positions := pipeDisplayPositions(line)
		if len(positions) == 0 {
			continue
		}
		if baseline == nil {
			baseline = positions
			continue
		}
		if len(positions) != len(baseline) {
			t.Fatalf("pipe count mismatch for %q: got %v want %v", line, positions, baseline)
		}
		for i := range positions {
			if positions[i] != baseline[i] {
				t.Fatalf("pipe position mismatch for %q: got %v want %v", line, positions, baseline)
			}
		}
	}
}

func pipeDisplayPositions(line string) []int {
	positions := make([]int, 0, strings.Count(line, "|"))
	width := 0
	for _, r := range line {
		if r == '|' {
			positions = append(positions, width)
		}
		width += runewidth.RuneWidth(r)
	}
	return positions
}
