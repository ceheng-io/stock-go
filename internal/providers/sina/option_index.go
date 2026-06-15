package sina

import (
	"context"
	"net/url"

	"github.com/ceheng.io/stock-go/types"
)

// IndexOptionProduct 是中金所股指期权产品。
type IndexOptionProduct = types.IndexOptionProduct

const (
	IndexOptionProductHO IndexOptionProduct = types.IndexOptionProductHO
	IndexOptionProductIO IndexOptionProduct = types.IndexOptionProductIO
	IndexOptionProductMO IndexOptionProduct = types.IndexOptionProductMO
)

type sinaOptionSpotResponse struct {
	Result struct {
		Data struct {
			Up   [][]string `json:"up"`
			Down [][]string `json:"down"`
		} `json:"data"`
	} `json:"result"`
}

type sinaOptionKlineItem struct {
	Date   string `json:"d"`
	Open   string `json:"o"`
	High   string `json:"h"`
	Low    string `json:"l"`
	Close  string `json:"c"`
	Volume string `json:"v"`
}

// GetIndexOptionSpot 获取新浪中金所股指期权 T 型报价。
func GetIndexOptionSpot(ctx context.Context, client JSONPClient, endpoint string, product IndexOptionProduct, contract string) (types.OptionTQuoteResult, error) {
	params := url.Values{}
	params.Set("type", "futures")
	params.Set("product", string(product))
	params.Set("exchange", "cffex")
	params.Set("pinzhong", contract)

	var payload sinaOptionSpotResponse
	if err := getSinaJSONP(ctx, client, endpoint, params, &payload); err != nil {
		return types.OptionTQuoteResult{}, err
	}
	return types.OptionTQuoteResult{
		Calls: parseOptionCallQuotes(payload.Result.Data.Up),
		Puts:  parseOptionPutQuotes(payload.Result.Data.Down),
	}, nil
}

// GetIndexOptionKline 获取新浪中金所股指期权合约日 K 线。
func GetIndexOptionKline(ctx context.Context, client JSONPClient, endpoint string, symbol string) ([]types.OptionKline, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	return getSinaOptionKlinesJSONP(ctx, client, endpoint, params)
}

func parseOptionCallQuotes(rows [][]string) []types.OptionTQuote {
	quotes := make([]types.OptionTQuote, 0, len(rows))
	for _, row := range rows {
		quotes = append(quotes, types.OptionTQuote{
			BuyVolume:    sinaNumber(sinaField(row, 0)),
			BuyPrice:     sinaNumber(sinaField(row, 1)),
			Price:        sinaNumber(sinaField(row, 2)),
			AskPrice:     sinaNumber(sinaField(row, 3)),
			AskVolume:    sinaNumber(sinaField(row, 4)),
			OpenInterest: sinaNumber(sinaField(row, 5)),
			Change:       sinaNumber(sinaField(row, 6)),
			StrikePrice:  sinaNumber(sinaField(row, 7)),
			Symbol:       sinaField(row, 8),
		})
	}
	return quotes
}

func parseOptionPutQuotes(rows [][]string) []types.OptionTQuote {
	quotes := make([]types.OptionTQuote, 0, len(rows))
	for _, row := range rows {
		quotes = append(quotes, types.OptionTQuote{
			BuyVolume:    sinaNumber(sinaField(row, 0)),
			BuyPrice:     sinaNumber(sinaField(row, 1)),
			Price:        sinaNumber(sinaField(row, 2)),
			AskPrice:     sinaNumber(sinaField(row, 3)),
			AskVolume:    sinaNumber(sinaField(row, 4)),
			OpenInterest: sinaNumber(sinaField(row, 5)),
			Change:       sinaNumber(sinaField(row, 6)),
			Symbol:       sinaField(row, 7),
		})
	}
	return quotes
}

func parseOptionKlines(items []sinaOptionKlineItem) []types.OptionKline {
	rows := make([]types.OptionKline, 0, len(items))
	for _, item := range items {
		rows = append(rows, types.OptionKline{
			Date:   item.Date,
			Open:   sinaNumber(item.Open),
			High:   sinaNumber(item.High),
			Low:    sinaNumber(item.Low),
			Close:  sinaNumber(item.Close),
			Volume: sinaNumber(item.Volume),
		})
	}
	return rows
}

func sinaField(fields []string, index int) string {
	if index < 0 || index >= len(fields) {
		return ""
	}
	return fields[index]
}
