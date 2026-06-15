package screener

import "testing"

type quoteFixture struct {
	Code          string
	PE            *float64
	ChangePercent *float64
	Amount        *float64
}

func TestScreenFiltersSortsAndTops(t *testing.T) {
	items := []quoteFixture{
		{Code: "600000", PE: Float(12), ChangePercent: Float(4), Amount: Float(100)},
		{Code: "600001", PE: Float(30), ChangePercent: Float(5), Amount: Float(500)},
		{Code: "600002", PE: Float(18), ChangePercent: Float(2), Amount: Float(300)},
		{Code: "600003", PE: Float(16), ChangePercent: Float(6), Amount: nil},
	}

	picks, err := Screen(items).
		Where(func(item quoteFixture) bool { return item.PE != nil && *item.PE < 20 }).
		Where(func(item quoteFixture) bool { return item.ChangePercent != nil && *item.ChangePercent > 3 }).
		SortBy(func(item quoteFixture) *float64 { return item.Amount }, Desc).
		Top(2)
	if err != nil {
		t.Fatal(err)
	}

	if len(picks) != 2 {
		t.Fatalf("len(picks) = %d, want 2", len(picks))
	}
	if picks[0].Code != "600000" || picks[1].Code != "600003" {
		t.Fatalf("picks = %+v", picks)
	}
}

func TestScreenRejectsNegativeTop(t *testing.T) {
	_, err := Screen([]int{1, 2, 3}).Top(-1)
	if err == nil {
		t.Fatal("expected error")
	}
}
