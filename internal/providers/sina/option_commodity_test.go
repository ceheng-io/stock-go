package sina

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestGetCommodityOptionSpotBuildsMappedRequestAndParsesTQuotes(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":{"up":[["10","1.10","1.20","1.30","20","300","0.05","580","au2506C580"]],"down":[["11","2.10","2.20","2.30","21","301","-0.06","au2506P580"]]}}})`,
	}}

	result, err := GetCommodityOptionSpot(context.Background(), client, "https://sina.test/option", "au", "au2506")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Calls) != 1 || len(result.Puts) != 1 {
		t.Fatalf("result = %+v", result)
	}
	if result.Calls[0].Symbol != "au2506C580" || result.Calls[0].StrikePrice == nil || *result.Calls[0].StrikePrice != 580 {
		t.Fatalf("call = %+v", result.Calls[0])
	}
	if result.Puts[0].Symbol != "au2506P580" || result.Puts[0].StrikePrice != nil {
		t.Fatalf("put = %+v", result.Puts[0])
	}

	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	if query.Get("type") != "futures" || query.Get("product") != "au_o" || query.Get("exchange") != "shfe" || query.Get("pinzhong") != "au2506" {
		t.Fatalf("query = %v", query)
	}
}

func TestGetCommodityOptionSpotPreservesUppercaseCZCEVariety(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`callback({"result":{"data":{"up":[],"down":[]}}})`,
	}}

	if _, err := GetCommodityOptionSpot(context.Background(), client, "https://sina.test/option", "SR", "SR501"); err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(client.urls[0])
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Query().Get("product") != "SR_o" || parsed.Query().Get("exchange") != "czce" {
		t.Fatalf("query = %v", parsed.Query())
	}
}

func TestGetCommodityOptionSpotReturnsUnknownVarietyError(t *testing.T) {
	client := &fakeETFOptionClient{}

	_, err := GetCommodityOptionSpot(context.Background(), client, "https://sina.test/option", "unknown", "unknown2506")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `unknown commodity option variety "unknown"`) {
		t.Fatalf("err = %v", err)
	}
	if len(client.urls) != 0 {
		t.Fatalf("urls = %#v", client.urls)
	}
}

func TestGetCommodityOptionKlineUsesPathCallbackAndParsesRows(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`cb([{"d":"2024-06-12","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"},{"d":"2024-06-13","o":"-","h":"1.28","l":"","c":"1.18","v":"900"}])`,
	}}

	rows, err := GetCommodityOptionKline(context.Background(), client, "https://sina.test/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline", "au2506C580")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Date != "2024-06-12" || rows[1].Date != "2024-06-13" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Close == nil || *rows[0].Close != 1.20 || rows[0].Volume == nil || *rows[0].Volume != 1000 {
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
	if parsed.Query().Get("symbol") != "au2506C580" {
		t.Fatalf("query = %v", parsed.Query())
	}
}

func TestGetCommodityOptionKlineReturnsEmptyRowsForNonArrayPayload(t *testing.T) {
	client := &fakeETFOptionClient{responses: []string{
		`cb({"result":{"data":null}})`,
	}}

	rows, err := GetCommodityOptionKline(context.Background(), client, "https://sina.test/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline", "au2506C580")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty rows", rows)
	}
}
