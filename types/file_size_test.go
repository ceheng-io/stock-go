package types

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const maxTypesFileLines = 1000

func TestTypesFilesStaySmall(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Count(string(body), "\n")
		if lines > maxTypesFileLines {
			t.Fatalf("%s has %d lines, want <= %d; split types by domain", file, lines, maxTypesFileLines)
		}
	}
}

func TestDomainTypesStayInDomainFiles(t *testing.T) {
	tests := []struct {
		file     string
		required []string
	}{
		{
			file: "block_trade.go",
			required: []string{
				"type BlockTradeMarketStatItem",
				"type BlockTradeDetailItem",
				"type BlockTradeDailyStatItem",
			},
		},
		{
			file: "margin.go",
			required: []string{
				"type MarginAccountItem",
				"type MarginTargetItem",
			},
		},
		{
			file: "dividend.go",
			required: []string{
				"type DividendDetail",
			},
		},
	}

	for _, tt := range tests {
		body, err := os.ReadFile(tt.file)
		if err != nil {
			t.Fatal(err)
		}
		text := string(body)
		for _, required := range tt.required {
			if !strings.Contains(text, required) {
				t.Fatalf("%s does not contain %q; keep domain types in their own files", tt.file, required)
			}
		}
	}
}
