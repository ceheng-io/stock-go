package eastmoney

import (
	"context"
	"net/url"
	"testing"
)

type fakeFundClient struct {
	lastURL string
	urls    []string
	text    string
	texts   []string
}

func (f *fakeFundClient) GetText(_ context.Context, requestURL string) (string, error) {
	f.lastURL = requestURL
	f.urls = append(f.urls, requestURL)
	if len(f.texts) > 0 {
		text := f.texts[0]
		f.texts = f.texts[1:]
		return text, nil
	}
	return f.text, nil
}

func TestGetFundEstimateBuildsRequestAndParsesJSONP(t *testing.T) {
	client := &fakeFundClient{text: `jsonpgz({"fundcode":"110011","name":"易方达中小盘","jzrq":"2024-12-16","dwjz":"3.5000","gsz":"3.5600","gszzl":"1.71","gztime":"2024-12-17 15:00"});`}

	row, err := GetFundEstimate(context.Background(), client, "https://fundgz.test/js", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if row.Code != "110011" || row.Name == nil || *row.Name != "易方达中小盘" || row.NavDate == nil || *row.NavDate != "2024-12-16" {
		t.Fatalf("row = %+v", row)
	}
	if row.Nav == nil || *row.Nav != 3.5 || row.EstimatedNav == nil || *row.EstimatedNav != 3.56 || row.EstimatedChangePercent == nil || *row.EstimatedChangePercent != 1.71 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	if parsed.Path != "/js/110011.js" {
		t.Fatalf("path = %q", parsed.Path)
	}
	if parsed.Query().Get("rt") == "" {
		t.Fatalf("missing rt in %s", client.lastURL)
	}
}

func TestGetFundEstimateUsesInputCodeAndNilNumbersForBlankPayload(t *testing.T) {
	client := &fakeFundClient{text: `jsonpgz({"name":"","dwjz":"--","gsz":"","gszzl":"","gztime":""});`}

	row, err := GetFundEstimate(context.Background(), client, "https://fundgz.test/js", "005827")
	if err != nil {
		t.Fatal(err)
	}
	if row.Code != "005827" || row.Name != nil || row.NavDate != nil || row.Nav != nil || row.EstimatedNav != nil || row.EstimatedChangePercent != nil || row.EstimateTime != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetFundEstimateUsesFallbackWhenJSONPIsBlank(t *testing.T) {
	client := &fakeFundClient{text: "  "}

	row, err := GetFundEstimate(context.Background(), client, "https://fundgz.test/js", "005827")
	if err != nil {
		t.Fatal(err)
	}
	if row.Code != "005827" || row.Name != nil || row.NavDate != nil || row.Nav != nil || row.EstimatedNav != nil || row.EstimatedChangePercent != nil || row.EstimateTime != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetFundEstimateUsesFallbackWhenJSONPIsInvalid(t *testing.T) {
	client := &fakeFundClient{text: "not jsonp"}

	row, err := GetFundEstimate(context.Background(), client, "https://fundgz.test/js", "005827")
	if err != nil {
		t.Fatal(err)
	}
	if row.Code != "005827" || row.Name != nil || row.NavDate != nil || row.Nav != nil || row.EstimatedNav != nil || row.EstimatedChangePercent != nil || row.EstimateTime != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetFundNavHistoryBuildsRequestAndMergesAccumulatedNav(t *testing.T) {
	client := &fakeFundClient{text: `
var fS_code = "110011";
var fS_name = "易方达中小盘";
var Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"equityReturn":"1.2","unitMoney":""},{"x":1702944000000,"y":3.6,"equityReturn":"","unitMoney":"0.01"}];
var Data_ACWorthTrend = [[1702857600000,5.6]];
`}

	history, err := GetFundNavHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if history.Code != "110011" || history.Name == nil || *history.Name != "易方达中小盘" || len(history.Items) != 2 {
		t.Fatalf("history = %+v", history)
	}
	first := history.Items[0]
	if first.Date != "2023-12-18" || first.Timestamp == nil || *first.Timestamp != 1702857600000 || first.Nav != 3.5 || first.AccNav == nil || *first.AccNav != 5.6 || first.DailyReturn == nil || *first.DailyReturn != 1.2 {
		t.Fatalf("first = %+v", first)
	}
	second := history.Items[1]
	if second.AccNav != nil || second.DailyReturn != nil || second.UnitMoney != "0.01" {
		t.Fatalf("second = %+v", second)
	}
	parsed, _ := url.Parse(client.lastURL)
	if parsed.Path != "/pingzhongdata/110011.js" {
		t.Fatalf("path = %q", parsed.Path)
	}
}

func TestGetFundNavHistoryUsesFallbackCodeAndNilName(t *testing.T) {
	client := &fakeFundClient{text: `
var Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"equityReturn":"","unitMoney":""}];
var Data_ACWorthTrend = [];
`}

	history, err := GetFundNavHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if history.Code != "110011" || history.Name != nil || len(history.Items) != 1 {
		t.Fatalf("history = %+v", history)
	}
}

func TestGetFundNavHistoryKeepsNilTimestampForMissingPointTime(t *testing.T) {
	client := &fakeFundClient{text: `
var Data_netWorthTrend = [{"y":3.5,"equityReturn":"","unitMoney":""}];
var Data_ACWorthTrend = [];
`}

	history, err := GetFundNavHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(history.Items) != 1 {
		t.Fatalf("history = %+v", history)
	}
	point := history.Items[0]
	if point.Timestamp != nil || point.Date != "" || point.AccNav != nil {
		t.Fatalf("point = %+v", point)
	}
}

func TestGetFundNavHistorySkipsInvalidAccumulatedRowsLikeTypeScript(t *testing.T) {
	client := &fakeFundClient{text: `
var Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"equityReturn":"","unitMoney":""},{"x":1702944000000,"y":3.6,"equityReturn":"","unitMoney":""}];
var Data_ACWorthTrend = [{"bad":true},[1702857600000,5.6],[1702944000000]];
`}

	history, err := GetFundNavHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(history.Items) != 2 {
		t.Fatalf("history = %+v", history)
	}
	if history.Items[0].AccNav == nil || *history.Items[0].AccNav != 5.6 {
		t.Fatalf("first = %+v, want accumulated nav merged", history.Items[0])
	}
	if history.Items[1].AccNav != nil {
		t.Fatalf("second = %+v, want invalid short row skipped", history.Items[1])
	}
}

func TestGetFundRankHistoryBuildsRequestAndMergesPercentile(t *testing.T) {
	client := &fakeFundClient{text: `
var fS_code = "110011";
var fS_name = "易方达中小盘";
var Data_rateInSimilarType = [{"x":1702857600000,"y":"12","sc":"300"},{"x":1702944000000,"y":"","sc":""}];
var Data_rateInSimilarPersent = [[1702857600000,4.0]];
`}

	history, err := GetFundRankHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if history.Code != "110011" || history.Name == nil || *history.Name != "易方达中小盘" || len(history.Items) != 2 {
		t.Fatalf("history = %+v", history)
	}
	first := history.Items[0]
	if first.Date != "2023-12-18" || first.Timestamp == nil || *first.Timestamp != 1702857600000 || first.Rank == nil || *first.Rank != 12 || first.Total == nil || *first.Total != 300 || first.Percentile == nil || *first.Percentile != 4 {
		t.Fatalf("first = %+v", first)
	}
	second := history.Items[1]
	if second.Rank != nil || second.Total != nil || second.Percentile != nil {
		t.Fatalf("second = %+v", second)
	}
}

func TestGetFundRankHistoryUsesFallbackCodeAndNilName(t *testing.T) {
	client := &fakeFundClient{text: `
var Data_rateInSimilarType = [{"x":1702857600000,"y":"","sc":""}];
var Data_rateInSimilarPersent = [];
`}

	history, err := GetFundRankHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if history.Code != "110011" || history.Name != nil || len(history.Items) != 1 {
		t.Fatalf("history = %+v", history)
	}
}

func TestGetFundRankHistoryKeepsNilTimestampForMissingPointTime(t *testing.T) {
	client := &fakeFundClient{text: `
var Data_rateInSimilarType = [{"y":"12","sc":"300"}];
var Data_rateInSimilarPersent = [];
`}

	history, err := GetFundRankHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(history.Items) != 1 {
		t.Fatalf("history = %+v", history)
	}
	point := history.Items[0]
	if point.Timestamp != nil || point.Date != "" || point.Percentile != nil {
		t.Fatalf("point = %+v", point)
	}
}

func TestGetFundRankHistorySkipsInvalidPercentileRowsLikeTypeScript(t *testing.T) {
	client := &fakeFundClient{text: `
var Data_rateInSimilarType = [{"x":1702857600000,"y":"12","sc":"300"},{"x":1702944000000,"y":"13","sc":"300"}];
var Data_rateInSimilarPersent = [{"bad":true},[1702857600000,4.0],[1702944000000]];
`}

	history, err := GetFundRankHistory(context.Background(), client, "https://fund.test/pingzhongdata", "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(history.Items) != 2 {
		t.Fatalf("history = %+v", history)
	}
	if history.Items[0].Percentile == nil || *history.Items[0].Percentile != 4 {
		t.Fatalf("first = %+v, want percentile merged", history.Items[0])
	}
	if history.Items[1].Percentile != nil {
		t.Fatalf("second = %+v, want invalid short row skipped", history.Items[1])
	}
}

func TestGetFundDividendListBuildsRequestAndParsesPage(t *testing.T) {
	client := &fakeFundClient{text: `
var pageinfo = [3,20,2];
var jjfh_data = [["110011","易方达中小盘","2024-12-16","2024-12-17","0.12","2024-12-18","混合型"],["000001","华夏成长","","","--","",""]];
`}

	result, err := GetFundDividendList(context.Background(), client, "https://fund.test/Data/funddataIndex_Interface.aspx", FundDividendListOptions{
		Year:     "2024",
		Page:     2,
		FundType: "混合型",
		Rank:     FundDividendRankCode,
		Sort:     FundSortAsc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalPages != 3 || result.PageSize != 20 || result.CurrentPage != 2 || len(result.Items) != 2 {
		t.Fatalf("result = %+v", result)
	}
	first := result.Items[0]
	if first.Code != "110011" || first.Name != "易方达中小盘" || first.EquityRecordDate == nil || *first.EquityRecordDate != "2024-12-16" || first.DividendPerShare == nil || *first.DividendPerShare != 0.12 {
		t.Fatalf("first = %+v", first)
	}
	second := result.Items[1]
	if second.EquityRecordDate != nil || second.ExDividendDate != nil || second.DividendPerShare != nil || second.PayDate != nil {
		t.Fatalf("second = %+v", second)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	if query.Get("dt") != "8" || query.Get("page") != "2" || query.Get("year") != "2024" || query.Get("ftype") != "混合型" || query.Get("rank") != "BZDM" || query.Get("sort") != "asc" {
		t.Fatalf("query = %s", parsed.RawQuery)
	}
}

func TestGetFundDividendListAllPagesFiltersByCode(t *testing.T) {
	client := &fakeFundClient{texts: []string{
		`var pageinfo = [2,20,1]; var jjfh_data = [["110011","易方达中小盘","2024-12-16","2024-12-17","0.12","2024-12-18","混合型"],["000001","华夏成长","2024-01-01","2024-01-02","0.03","2024-01-03","混合型"]];`,
		`var pageinfo = [2,20,2]; var jjfh_data = [["110011","易方达中小盘","2023-12-16","2023-12-17","0.10","2023-12-18","混合型"]];`,
	}}

	result, err := GetFundDividendList(context.Background(), client, "https://fund.test/Data/funddataIndex_Interface.aspx", FundDividendListOptions{
		Year: "2024",
		Page: "all",
		Code: "110011",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalPages != 2 || result.CurrentPage != -1 || len(result.Items) != 2 {
		t.Fatalf("result = %+v", result)
	}
	if len(client.urls) != 2 {
		t.Fatalf("requests = %d, want 2", len(client.urls))
	}
	firstURL, _ := url.Parse(client.urls[0])
	secondURL, _ := url.Parse(client.urls[1])
	if firstURL.Query().Get("page") != "1" || secondURL.Query().Get("page") != "2" {
		t.Fatalf("urls = %#v", client.urls)
	}
}

func TestGetFundDividendListAllPagesCompatibilityFlag(t *testing.T) {
	client := &fakeFundClient{texts: []string{
		`var pageinfo = [1,20,1]; var jjfh_data = [["110011","易方达中小盘","2024-12-16","2024-12-17","0.12","2024-12-18","混合型"]];`,
	}}

	result, err := GetFundDividendList(context.Background(), client, "https://fund.test/Data/funddataIndex_Interface.aspx", FundDividendListOptions{
		Year:     "2024",
		AllPages: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.CurrentPage != -1 || len(client.urls) != 1 {
		t.Fatalf("result=%+v urls=%#v", result, client.urls)
	}
}
