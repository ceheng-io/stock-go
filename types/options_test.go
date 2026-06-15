package types

import "testing"

func TestOptionCategoryConstantsMatchTSUnions(t *testing.T) {
	products := []IndexOptionProduct{
		IndexOptionProductHO,
		IndexOptionProductIO,
		IndexOptionProductMO,
	}
	wantProducts := []string{"ho", "io", "mo"}
	for index, value := range products {
		if string(value) != wantProducts[index] {
			t.Fatalf("products[%d] = %q, want %q", index, value, wantProducts[index])
		}
	}

	categories := []ETFOptionCate{
		ETFOptionCate50ETF,
		ETFOptionCate300ETF,
		ETFOptionCate500ETF,
		ETFOptionCateKechuang50,
		ETFOptionCateKechuangBoard50,
	}
	wantCategories := []string{"50ETF", "300ETF", "500ETF", "科创50", "科创板50"}
	for index, value := range categories {
		if string(value) != wantCategories[index] {
			t.Fatalf("categories[%d] = %q, want %q", index, value, wantCategories[index])
		}
	}
}
