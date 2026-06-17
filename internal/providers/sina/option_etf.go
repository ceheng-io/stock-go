package sina

import (
	"context"
	"net/url"

	"github.com/ceheng-io/stock-go/internal/core"
	"github.com/ceheng-io/stock-go/types"
)

// JSONPClient 是新浪 provider 所需的最小请求客户端接口。
type JSONPClient interface {
	GetText(context.Context, string) (string, error)
}

// ETFOptionCate 是新浪 ETF 期权分类。
type ETFOptionCate = types.ETFOptionCate

const (
	ETFOptionCate50ETF           ETFOptionCate = types.ETFOptionCate50ETF
	ETFOptionCate300ETF          ETFOptionCate = types.ETFOptionCate300ETF
	ETFOptionCate500ETF          ETFOptionCate = types.ETFOptionCate500ETF
	ETFOptionCateKechuang50      ETFOptionCate = types.ETFOptionCateKechuang50
	ETFOptionCateKechuangBoard50 ETFOptionCate = types.ETFOptionCateKechuangBoard50
)

type etfOptionMonthsResponse struct {
	Result struct {
		Data struct {
			ContractMonth []string `json:"contractMonth"`
			StockID       string   `json:"stockId"`
			CateID        string   `json:"cateId"`
			CateList      []string `json:"cateList"`
		} `json:"data"`
	} `json:"result"`
}

type etfOptionExpireResponse struct {
	Result struct {
		Data struct {
			ExpireDay     string `json:"expireDay"`
			RemainderDays int    `json:"remainderDays"`
			StockID       string `json:"stockId"`
			Other         struct {
				Name string `json:"name"`
			} `json:"other"`
		} `json:"data"`
	} `json:"result"`
}

// GetETFOptionMonths 获取新浪 ETF 期权可用月份。
func GetETFOptionMonths(ctx context.Context, client JSONPClient, endpoint string, cate ETFOptionCate) (types.ETFOptionMonth, error) {
	params := url.Values{}
	params.Set("exchange", "null")
	params.Set("cate", string(cate))

	var payload etfOptionMonthsResponse
	if err := getSinaJSONP(ctx, client, endpoint, params, &payload); err != nil {
		return types.ETFOptionMonth{}, err
	}
	data := payload.Result.Data
	months := data.ContractMonth
	if len(months) > 1 {
		months = months[1:]
	}
	return types.ETFOptionMonth{
		Months:   cloneStringSlice(months),
		StockID:  data.StockID,
		CateID:   data.CateID,
		CateList: cloneStringSlice(data.CateList),
	}, nil
}

func cloneStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return append([]string(nil), values...)
}

// GetETFOptionExpireDay 获取新浪 ETF 期权到期日。
func GetETFOptionExpireDay(ctx context.Context, client JSONPClient, endpoint string, cate ETFOptionCate, month string) (types.ETFOptionExpireDay, error) {
	expire, err := fetchETFOptionExpireDay(ctx, client, endpoint, cate, month)
	if err != nil {
		return types.ETFOptionExpireDay{}, err
	}
	if expire.RemainderDays < 0 {
		return fetchETFOptionExpireDay(ctx, client, endpoint, ETFOptionCate("XD"+string(cate)), month)
	}
	return expire, nil
}

func fetchETFOptionExpireDay(ctx context.Context, client JSONPClient, endpoint string, cate ETFOptionCate, month string) (types.ETFOptionExpireDay, error) {
	params := url.Values{}
	params.Set("exchange", "null")
	params.Set("cate", string(cate))
	params.Set("date", month)

	var payload etfOptionExpireResponse
	if err := getSinaJSONP(ctx, client, endpoint, params, &payload); err != nil {
		return types.ETFOptionExpireDay{}, err
	}
	data := payload.Result.Data
	return types.ETFOptionExpireDay{
		ExpireDay:     data.ExpireDay,
		RemainderDays: data.RemainderDays,
		StockID:       data.StockID,
		Name:          data.Other.Name,
	}, nil
}

func getSinaJSONP(ctx context.Context, client JSONPClient, endpoint string, params url.Values, target any) error {
	text, err := client.GetText(ctx, endpoint+"?"+params.Encode())
	if err != nil {
		return err
	}
	return core.ExtractJSONP(text, target)
}
