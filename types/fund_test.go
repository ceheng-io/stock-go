package types

import "testing"

func TestFundDividendOptionTypesMatchTSUnions(t *testing.T) {
	ranks := []FundDividendRank{
		FundDividendRankCode,
		FundDividendRankName,
		FundDividendRankEquityRecordDate,
		FundDividendRankExDividendDate,
		FundDividendRankDividendPerShare,
		FundDividendRankPayDate,
	}
	wantRanks := []string{"BZDM", "ABBNAME", "DJR", "FSRQ", "FHFCZ", "FFR"}
	for index, value := range ranks {
		if string(value) != wantRanks[index] {
			t.Fatalf("ranks[%d] = %q, want %q", index, value, wantRanks[index])
		}
	}

	sorts := []FundSortDirection{FundSortAsc, FundSortDesc}
	wantSorts := []string{"asc", "desc"}
	for index, value := range sorts {
		if string(value) != wantSorts[index] {
			t.Fatalf("sorts[%d] = %q, want %q", index, value, wantSorts[index])
		}
	}

	options := FundDividendListOptions{
		Year:     "2024",
		Page:     "all",
		FundType: "股票型",
		Rank:     FundDividendRankExDividendDate,
		Sort:     FundSortDesc,
		Code:     "110011",
	}
	if options.Rank != FundDividendRankExDividendDate || options.Sort != FundSortDesc {
		t.Fatalf("options = %+v", options)
	}
}
