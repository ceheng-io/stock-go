package sina

import (
	"context"
	"net/url"
	"testing"
)

func TestGetIndexOptionSpotBuildsRequestAndParsesTQuotes(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"status":{"code":0},"data":{"up":[["10","1.10","1.20","1.30","20","300","0.05","3600","io2504C3600"]],"down":[["11","2.10","2.20","2.30","21","301","-0.06","io2504P3600"]]}}})`,
	}}

	result, err := GetIndexOptionSpot(context.Background(), client, "https://sina.test/option", IndexOptionProductIO, "io2504")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Calls) != 1 || len(result.Puts) != 1 {
		t.Fatalf("result = %+v", result)
	}
	call := result.Calls[0]
	if call.Symbol != "io2504C3600" || call.Price == nil || *call.Price != 1.20 || call.StrikePrice == nil || *call.StrikePrice != 3600 {
		t.Fatalf("call = %+v", call)
	}
	put := result.Puts[0]
	if put.Symbol != "io2504P3600" || put.Change == nil || *put.Change != -0.06 || put.StrikePrice != nil {
		t.Fatalf("put = %+v", put)
	}

	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	if query.Get("type") != "futures" || query.Get("product") != "io" || query.Get("exchange") != "cffex" || query.Get("pinzhong") != "io2504" {
		t.Fatalf("query = %v", query)
	}
}

func TestGetIndexOptionKlineUsesPathCallbackAndParsesRows(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`cb([{"d":"2024-06-12","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"},{"d":"2024-06-13","o":"","h":"1.28","l":"-","c":"1.18","v":"900"}])`,
	}}

	rows, err := GetIndexOptionKline(context.Background(), client, "https://sina.test/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline", "io2504C3600")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Date != "2024-06-12" || rows[1].Date != "2024-06-13" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Open == nil || *rows[0].Open != 1.10 || rows[0].Volume == nil || *rows[0].Volume != 1000 {
		t.Fatalf("first row = %+v", rows[0])
	}
	if rows[1].Open != nil || rows[1].Low != nil || rows[1].High == nil || *rows[1].High != 1.28 {
		t.Fatalf("second row = %+v", rows[1])
	}

	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Path != "/jsonp.php/ceheng_jsonp/FutureOptionAllService.getOptionDayline" {
		t.Fatalf("path = %q", parsed.Path)
	}
	if parsed.Query().Get("symbol") != "io2504C3600" {
		t.Fatalf("query = %v", parsed.Query())
	}
}

func TestGetIndexOptionKlineReturnsEmptyRowsForNonArrayPayload(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`cb({"result":{"data":null}})`,
	}}

	rows, err := GetIndexOptionKline(context.Background(), client, "https://sina.test/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline", "io2504C3600")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty rows", rows)
	}
}
