package types

import "testing"

func TestBlockTradeDateOptionsMatchesTSShape(t *testing.T) {
	options := BlockTradeDateOptions{StartDate: "20241201", EndDate: "2024-12-31"}
	if options.StartDate != "20241201" || options.EndDate != "2024-12-31" {
		t.Fatalf("options = %+v", options)
	}
}
