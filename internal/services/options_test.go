package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/internal/providers/sina"
)

type optionsClientStub struct {
	lastURL   string
	urls      []string
	payload   map[string]any
	responses []string
}

func (o *optionsClientStub) GetJSON(_ context.Context, requestURL string, target any) error {
	o.lastURL = requestURL
	o.urls = append(o.urls, requestURL)
	body, _ := json.Marshal(o.payload)
	return json.Unmarshal(body, target)
}

func (o *optionsClientStub) GetText(_ context.Context, requestURL string) (string, error) {
	o.lastURL = requestURL
	o.urls = append(o.urls, requestURL)
	if len(o.responses) == 0 {
		return `callback({})`, nil
	}
	text := o.responses[0]
	o.responses = o.responses[1:]
	return text, nil
}

func TestOptionsServiceCFFEXQuotesAndLHB(t *testing.T) {
	client := &optionsClientStub{payload: map[string]any{
		"list": []map[string]any{{"dm": "io2501C4000", "name": "沪深300购2501", "p": 123.4}},
	}}
	service := NewOptionsService(client, OptionsURLs{
		CFFEXQuotes: "https://em.test/option",
		LHB:         "https://em.test/lhb",
	})

	quotes, err := service.CFFEXQuotes(context.Background(), eastmoney.CFFEXOptionQuotesOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 || quotes[0].Code != "io2501C4000" || quotes[0].Price == nil || *quotes[0].Price != 123.4 {
		t.Fatalf("quotes = %+v", quotes)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "510050", "TRADE_DATE": "2024-12-16", "MEMBER_RANK": 1.0}},
		},
	}
	rows, err := service.LHB(context.Background(), "510050", "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Symbol != "510050" || rows[0].Rank != 1 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestOptionsServiceETFOptionMonthsAndExpireDay(t *testing.T) {
	client := &optionsClientStub{responses: []string{
		`callback({"result":{"data":{"contractMonth":["全部","2024-06"],"stockId":"510050","cateId":"50ETF","cateList":["50ETF"]}}})`,
		`callback({"result":{"data":{"expireDay":"2024-06-26","remainderDays":8,"stockId":"510050","other":{"name":"上证50ETF期权"}}}})`,
	}}
	service := NewOptionsService(client, OptionsURLs{
		ETFMonths: "https://sina.test/list",
		ETFExpire: "https://sina.test/expire",
	})

	months, err := service.ETFOptionMonths(context.Background(), sina.ETFOptionCate50ETF)
	if err != nil {
		t.Fatal(err)
	}
	if len(months.Months) != 1 || months.Months[0] != "2024-06" || months.StockID != "510050" {
		t.Fatalf("months = %+v", months)
	}

	expire, err := service.ETFOptionExpireDay(context.Background(), sina.ETFOptionCate50ETF, "2024-06")
	if err != nil {
		t.Fatal(err)
	}
	if expire.ExpireDay != "2024-06-26" || expire.RemainderDays != 8 || expire.Name != "上证50ETF期权" {
		t.Fatalf("expire = %+v", expire)
	}
}

func TestOptionsServiceETFOptionMinuteKlineAnd5Day(t *testing.T) {
	client := &optionsClientStub{responses: []string{
		`callback({"result":{"data":[{"i":"09:31","d":"2024-06-13","p":"1.23","v":"100","t":"200","a":"1.21"}]}})`,
		`cb([{"d":"2024-06-13","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"}])`,
		`callback({"result":{"data":[[{"i":"09:31","d":"2024-06-12","p":"1.10","v":"10","t":"20","a":"1.09"}],[{"i":"09:31","d":"2024-06-13","p":"1.20","v":"30","t":"40","a":"1.19"}]]}})`,
	}}
	service := NewOptionsService(client, OptionsURLs{
		ETFMinute: "https://sina.test/minute",
		ETFDaily:  "https://sina.test/jsonp_v2.php/{callback}/StockOptionDaylineService.getSymbolInfo",
		ETF5Day:   "https://sina.test/5day",
	})

	minutes, err := service.ETFOptionMinute(context.Background(), "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(minutes) != 1 || minutes[0].Price == nil || *minutes[0].Price != 1.23 {
		t.Fatalf("minutes = %+v", minutes)
	}

	klines, err := service.ETFOptionDailyKline(context.Background(), "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(klines) != 1 || klines[0].Close == nil || *klines[0].Close != 1.20 {
		t.Fatalf("klines = %+v", klines)
	}

	fiveDay, err := service.ETFOption5DayMinute(context.Background(), "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(fiveDay) != 2 || fiveDay[1].Date != "2024-06-13" || fiveDay[1].OpenInterest == nil || *fiveDay[1].OpenInterest != 40 {
		t.Fatalf("fiveDay = %+v", fiveDay)
	}
}

func TestOptionsServiceIndexOptionSpotAndKline(t *testing.T) {
	client := &optionsClientStub{responses: []string{
		`callback({"result":{"data":{"up":[["10","1.10","1.20","1.30","20","300","0.05","3600","io2504C3600"]],"down":[["11","2.10","2.20","2.30","21","301","-0.06","io2504P3600"]]}}})`,
		`cb([{"d":"2024-06-13","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"}])`,
	}}
	service := NewOptionsService(client, OptionsURLs{
		IndexSpot:  "https://sina.test/option",
		IndexKline: "https://sina.test/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline",
	})

	spot, err := service.IndexOptionSpot(context.Background(), sina.IndexOptionProductIO, "io2504")
	if err != nil {
		t.Fatal(err)
	}
	if len(spot.Calls) != 1 || spot.Calls[0].Symbol != "io2504C3600" || spot.Calls[0].StrikePrice == nil || *spot.Calls[0].StrikePrice != 3600 {
		t.Fatalf("spot = %+v", spot)
	}

	klines, err := service.IndexOptionKline(context.Background(), "io2504C3600")
	if err != nil {
		t.Fatal(err)
	}
	if len(klines) != 1 || klines[0].Close == nil || *klines[0].Close != 1.20 {
		t.Fatalf("klines = %+v", klines)
	}
}

func TestOptionsServiceCommodityOptionSpotAndKline(t *testing.T) {
	client := &optionsClientStub{responses: []string{
		`callback({"result":{"data":{"up":[["10","1.10","1.20","1.30","20","300","0.05","580","au2506C580"]],"down":[["11","2.10","2.20","2.30","21","301","-0.06","au2506P580"]]}}})`,
		`cb([{"d":"2024-06-13","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"}])`,
	}}
	service := NewOptionsService(client, OptionsURLs{
		CommoditySpot:  "https://sina.test/option",
		CommodityKline: "https://sina.test/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline",
	})

	spot, err := service.CommodityOptionSpot(context.Background(), "au", "au2506")
	if err != nil {
		t.Fatal(err)
	}
	if len(spot.Calls) != 1 || spot.Calls[0].Symbol != "au2506C580" || spot.Calls[0].StrikePrice == nil || *spot.Calls[0].StrikePrice != 580 {
		t.Fatalf("spot = %+v", spot)
	}

	klines, err := service.CommodityOptionKline(context.Background(), "au2506C580")
	if err != nil {
		t.Fatal(err)
	}
	if len(klines) != 1 || klines[0].Close == nil || *klines[0].Close != 1.20 {
		t.Fatalf("klines = %+v", klines)
	}
}
