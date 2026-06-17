package tencent

import (
	"context"
	"net/url"
	"testing"

	"github.com/ceheng-io/stock-go/types"
)

type fakeSearchClient struct {
	requestURL string
	text       string
}

func (f *fakeSearchClient) GetText(_ context.Context, requestURL string) (string, error) {
	f.requestURL = requestURL
	return f.text, nil
}

func (f *fakeSearchClient) TencentSearchURL(keyword string) string {
	return "https://smartbox.test/s3/?v=2&t=all&q=" + url.QueryEscape(keyword)
}

func TestSearchReturnsEmptyForBlankKeyword(t *testing.T) {
	client := &fakeSearchClient{text: `v_hint="sh~600519~name~~GP-A"`}

	got, err := Search(context.Background(), client, "  \t ")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("got %#v, want empty result", got)
	}
	if client.requestURL != "" {
		t.Fatalf("requestURL = %q, want no request", client.requestURL)
	}
}

func TestSearchReturnsEmptyForMissingOrNHint(t *testing.T) {
	for _, text := range []string{`console.log("no hint")`, `v_hint="N"`} {
		client := &fakeSearchClient{text: text}
		got, err := Search(context.Background(), client, "平安")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("text %q: got %#v, want empty result", text, got)
		}
	}
}

func TestSearchParsesSmartboxHintLikeTypeScript(t *testing.T) {
	client := &fakeSearchClient{text: `v_hint="sh~600519~\u8d35\u5dde\u8305\u53f0~GZMT~GP-A^jj~110011~\u6613\u65b9\u8fbe~YFD~QDII-ETF^hk~00700~Tencent~TENCENT~STOCK^sh~000001~\u4e0a\u8bc1\u6307\u6570~SZZS~ZS^xx~1~Other~O~UNKNOWN"`}

	got, err := Search(context.Background(), client, "茅台")
	if err != nil {
		t.Fatal(err)
	}

	want := []types.SearchResult{
		{Code: "sh600519", Name: "贵州茅台", Market: "sh", Type: "GP-A", Category: types.SearchStock},
		{Code: "jj110011", Name: "易方达", Market: "jj", Type: "QDII-ETF", Category: types.SearchFund},
		{Code: "hk00700", Name: "Tencent", Market: "hk", Type: "STOCK", Category: types.SearchStock},
		{Code: "sh000001", Name: "上证指数", Market: "sh", Type: "ZS", Category: types.SearchIndex},
		{Code: "xx1", Name: "Other", Market: "xx", Type: "UNKNOWN", Category: types.SearchOther},
	}

	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
	if client.requestURL != "https://smartbox.test/s3/?v=2&t=all&q=%E8%8C%85%E5%8F%B0" {
		t.Fatalf("requestURL = %q", client.requestURL)
	}
}

func TestNormalizeSearchTypeMatchesTypeScriptMapping(t *testing.T) {
	cases := map[string]types.SearchResultType{
		"QDII-LOF":     types.SearchFund,
		"ETF":          types.SearchFund,
		"LOF":          types.SearchFund,
		"KJ-HB":        types.SearchFund,
		"JJ":           types.SearchFund,
		"money-fund":   types.SearchFund,
		"GP-A":         types.SearchStock,
		"stock-us":     types.SearchStock,
		"ZS":           types.SearchIndex,
		"global-index": types.SearchIndex,
		"ZQ":           types.SearchBond,
		"corp-bond":    types.SearchBond,
		"QH":           types.SearchFutures,
		"future-main":  types.SearchFutures,
		"QZ":           types.SearchOption,
		"OPTION":       types.SearchOption,
		"UNKNOWN":      types.SearchOther,
	}

	for raw, want := range cases {
		if got := NormalizeSearchType(raw); got != want {
			t.Fatalf("NormalizeSearchType(%q) = %q, want %q", raw, got, want)
		}
	}
}
