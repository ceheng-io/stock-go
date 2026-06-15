package sina

import (
	"context"
	"net/url"
	"testing"
)

type fakeETFOptionClient struct {
	urls      []string
	responses []string
}

func (f *fakeETFOptionClient) GetText(_ context.Context, requestURL string) (string, error) {
	f.urls = append(f.urls, requestURL)
	if len(f.responses) == 0 {
		return `callback({})`, nil
	}
	text := f.responses[0]
	f.responses = f.responses[1:]
	return text, nil
}

func TestGetETFOptionMonthsBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":{"contractMonth":["全部","2024-06","2024-07"],"stockId":"510050","cateId":"50ETF","cateList":["50ETF","300ETF"]}}})`,
	}}

	months, err := GetETFOptionMonths(context.Background(), client, "https://sina.test/list", ETFOptionCate50ETF)
	if err != nil {
		t.Fatal(err)
	}
	if len(months.Months) != 2 || months.Months[0] != "2024-06" || months.Months[1] != "2024-07" {
		t.Fatalf("months = %+v", months)
	}
	if months.StockID != "510050" || months.CateID != "50ETF" || len(months.CateList) != 2 {
		t.Fatalf("months metadata = %+v", months)
	}

	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	if query.Get("exchange") != "null" || query.Get("cate") != "50ETF" {
		t.Fatalf("query = %v", query)
	}
}

func TestGetETFOptionMonthsReturnsEmptySlicesForMissingArrays(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":{"stockId":"510050","cateId":"50ETF"}}})`,
	}}

	months, err := GetETFOptionMonths(context.Background(), client, "https://sina.test/list", ETFOptionCate50ETF)
	if err != nil {
		t.Fatal(err)
	}
	if months.Months == nil {
		t.Fatalf("Months = nil, want empty slice")
	}
	if months.CateList == nil {
		t.Fatalf("CateList = nil, want empty slice")
	}
	if len(months.Months) != 0 || len(months.CateList) != 0 {
		t.Fatalf("months = %+v, want empty slices", months)
	}
}

func TestGetETFOptionExpireDayRetriesXDCategoryWhenRemainderNegative(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":{"expireDay":"2024-06-26","remainderDays":-1,"stockId":"510050","other":{"name":"上证50ETF期权"}}}})`,
		`callback({"result":{"data":{"expireDay":"2024-06-26","remainderDays":12,"stockId":"510050","other":{"name":"上证50ETF期权"}}}})`,
	}}

	expire, err := GetETFOptionExpireDay(context.Background(), client, "https://sina.test/expire", ETFOptionCate50ETF, "2024-06")
	if err != nil {
		t.Fatal(err)
	}
	if expire.ExpireDay != "2024-06-26" || expire.RemainderDays != 12 || expire.StockID != "510050" || expire.Name != "上证50ETF期权" {
		t.Fatalf("expire = %+v", expire)
	}
	if len(client.urls) != 2 {
		t.Fatalf("urls = %#v", client.urls)
	}

	first, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	second, err := url.Parse(client.urls[1])
	if err != nil {
		t.Fatal(err)
	}
	if first.Query().Get("cate") != "50ETF" || first.Query().Get("date") != "2024-06" {
		t.Fatalf("first query = %v", first.Query())
	}
	if second.Query().Get("cate") != "XD50ETF" || second.Query().Get("date") != "2024-06" {
		t.Fatalf("second query = %v", second.Query())
	}
}

func TestGetETFOptionMinuteBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":[{"i":"09:31","d":"2024-06-13","p":"1.23","v":"100","t":"200","a":"1.21"},{"i":"09:32","p":"1.24","v":"-","t":"","a":"1.22"}]}})`,
	}}

	rows, err := GetETFOptionMinute(context.Background(), client, "https://sina.test/minute", "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Date != "2024-06-13" || rows[1].Date != "2024-06-13" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Price == nil || *rows[0].Price != 1.23 || rows[0].Volume == nil || *rows[0].Volume != 100 || rows[0].OpenInterest == nil || *rows[0].OpenInterest != 200 {
		t.Fatalf("first row = %+v", rows[0])
	}
	if rows[1].Volume != nil || rows[1].OpenInterest != nil || rows[1].AvgPrice == nil || *rows[1].AvgPrice != 1.22 {
		t.Fatalf("second row = %+v", rows[1])
	}

	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Query().Get("symbol") != "CON_OP_10009633" {
		t.Fatalf("query = %v", parsed.Query())
	}
}

func TestGetETFOptionMinuteReturnsEmptyRowsForNonArrayPayload(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":{"not":"array"}}})`,
	}}

	rows, err := GetETFOptionMinute(context.Background(), client, "https://sina.test/minute", "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty rows", rows)
	}
}

func TestGetETFOptionDailyKlineUsesPathCallbackAndParsesRows(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`cb([{"d":"2024-06-12","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"},{"d":"2024-06-13","o":"-","h":"","l":"1.05","c":"1.18","v":"900"}])`,
	}}

	rows, err := GetETFOptionDailyKline(context.Background(), client, "https://sina.test/jsonp_v2.php/{callback}/StockOptionDaylineService.getSymbolInfo", "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Date != "2024-06-12" || rows[1].Date != "2024-06-13" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Open == nil || *rows[0].Open != 1.10 || rows[0].Close == nil || *rows[0].Close != 1.20 || rows[0].Volume == nil || *rows[0].Volume != 1000 {
		t.Fatalf("first row = %+v", rows[0])
	}
	if rows[1].Open != nil || rows[1].High != nil || rows[1].Low == nil || *rows[1].Low != 1.05 {
		t.Fatalf("second row = %+v", rows[1])
	}

	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Path != "/jsonp_v2.php/ceheng_jsonp/StockOptionDaylineService.getSymbolInfo" {
		t.Fatalf("path = %q", parsed.Path)
	}
	if parsed.Query().Get("symbol") != "CON_OP_10009633" {
		t.Fatalf("query = %v", parsed.Query())
	}
}

func TestGetETFOption5DayMinuteReturnsEmptyRowsForNonArrayPayload(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":{"not":"array"}}})`,
	}}

	rows, err := GetETFOption5DayMinute(context.Background(), client, "https://sina.test/5day", "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty rows", rows)
	}
}

func TestGetETFOption5DayMinuteSkipsNonArrayDayGroups(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":[{"not":"array"},[{"i":"09:31","d":"2024-06-13","p":"1.20","v":"30","t":"40","a":"1.19"}],null]}})`,
	}}

	rows, err := GetETFOption5DayMinute(context.Background(), client, "https://sina.test/5day", "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Date != "2024-06-13" {
		t.Fatalf("rows = %+v, want one parsed day group", rows)
	}
}

func TestGetETFOptionDailyKlineReturnsEmptyRowsForNonArrayPayload(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`cb({"result":{"data":null}})`,
	}}

	rows, err := GetETFOptionDailyKline(context.Background(), client, "https://sina.test/jsonp_v2.php/{callback}/StockOptionDaylineService.getSymbolInfo", "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty rows", rows)
	}
}

func TestGetETFOption5DayMinuteFlattensDayGroups(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":[[{"i":"09:31","d":"2024-06-12","p":"1.10","v":"10","t":"20","a":"1.09"}],[{"i":"09:31","d":"2024-06-13","p":"1.20","v":"30","t":"40","a":"1.19"},{"i":"09:32","p":"1.21","v":"31","t":"41","a":"1.20"}]]}})`,
	}}

	rows, err := GetETFOption5DayMinute(context.Background(), client, "https://sina.test/5day", "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Date != "2024-06-12" || rows[1].Date != "2024-06-13" || rows[2].Date != "2024-06-13" {
		t.Fatalf("dates = %+v", rows)
	}
	if rows[2].Price == nil || *rows[2].Price != 1.21 || rows[2].OpenInterest == nil || *rows[2].OpenInterest != 41 {
		t.Fatalf("last row = %+v", rows[2])
	}

	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Query().Get("symbol") != "CON_OP_10009633" {
		t.Fatalf("query = %v", parsed.Query())
	}
}
