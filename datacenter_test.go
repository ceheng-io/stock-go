package stock

import "testing"

func TestRootReExportsDatacenterUtilities(t *testing.T) {
	fetchAllPages := false
	query := DatacenterQuery{
		ReportName:    "RPT_TEST",
		PageSize:      100,
		FetchAllPages: &fetchAllPages,
	}
	if query.ReportName != "RPT_TEST" || query.PageSize != 100 || query.FetchAllPages == nil || *query.FetchAllPages {
		t.Fatalf("DatacenterQuery alias = %+v", query)
	}

	result := DatacenterResult[string]{
		Data:  []string{"a"},
		Total: 1,
		Pages: 1,
	}
	if result.Data[0] != "a" || result.Total != 1 || result.Pages != 1 {
		t.Fatalf("DatacenterResult alias = %+v", result)
	}

	if got := ParseDCDate("2024-01-15T00:00:00.000"); got != "2024-01-15" {
		t.Fatalf("ParseDCDate = %q", got)
	}
	if got := ParseDcDate("2024-01-15 00:00:00"); got != "2024-01-15" {
		t.Fatalf("ParseDcDate = %q", got)
	}
}
