package eastmoney

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ceheng-io/stock-go/internal/core"
	"github.com/ceheng-io/stock-go/types"
)

// FundClient is the minimal client interface required by Eastmoney fund providers.
type FundClient interface {
	GetText(context.Context, string) (string, error)
}

type fundGZPayload struct {
	FundCode string `json:"fundcode"`
	Name     string `json:"name"`
	NavDate  string `json:"jzrq"`
	Nav      string `json:"dwjz"`
	Estimate string `json:"gsz"`
	Change   string `json:"gszzl"`
	Time     string `json:"gztime"`
}

type fundNavTrendItem struct {
	Timestamp   *int64  `json:"x"`
	Nav         float64 `json:"y"`
	DailyReturn any     `json:"equityReturn"`
	UnitMoney   string  `json:"unitMoney"`
}

type fundRankItem struct {
	Timestamp *int64 `json:"x"`
	Rank      any    `json:"y"`
	Total     any    `json:"sc"`
}

// GetFundEstimate fetches a public fund's latest net-value estimate.
func GetFundEstimate(ctx context.Context, client FundClient, endpoint string, code string) (types.FundEstimate, error) {
	requestURL := strings.TrimRight(endpoint, "/") + "/" + url.PathEscape(code) + ".js?rt=" + url.QueryEscape(time.Now().Format("20060102150405.000000000"))
	text, err := client.GetText(ctx, requestURL)
	if err != nil {
		return types.FundEstimate{}, err
	}
	var raw fundGZPayload
	if strings.TrimSpace(text) != "" {
		_ = core.ExtractJSONP(text, &raw)
	}
	return parseFundEstimate(raw, code), nil
}

// GetFundNavHistory fetches a public fund's full net-value history.
func GetFundNavHistory(ctx context.Context, client FundClient, endpoint string, code string) (types.FundNavHistory, error) {
	text, err := fetchFundPingzhongText(ctx, client, endpoint, code)
	if err != nil {
		return types.FundNavHistory{}, err
	}
	fundCode, fundName := fundIdentity(text, code)
	var trend []fundNavTrendItem
	_ = core.ExtractJSVar(text, "Data_netWorthTrend", &trend)
	var accumulated [][]float64
	_ = core.ExtractJSVar(text, "Data_ACWorthTrend", &accumulated)
	accumulatedByTime := make(map[int64]float64, len(accumulated))
	for _, row := range accumulated {
		if len(row) >= 2 {
			accumulatedByTime[int64(row[0])] = row[1]
		}
	}
	items := make([]types.FundNavPoint, 0, len(trend))
	for _, point := range trend {
		var accNav *float64
		var date string
		if point.Timestamp != nil {
			if value, ok := accumulatedByTime[*point.Timestamp]; ok {
				value := value
				accNav = &value
			}
			date = timestampDate(*point.Timestamp)
		}
		items = append(items, types.FundNavPoint{
			Date:        date,
			Timestamp:   point.Timestamp,
			Nav:         point.Nav,
			AccNav:      accNav,
			DailyReturn: nullableNumberAny(point.DailyReturn),
			UnitMoney:   point.UnitMoney,
		})
	}
	return types.FundNavHistory{Code: fundCode, Name: fundName, Items: items}, nil
}

// GetFundRankHistory fetches a public fund's similar-type rank history.
func GetFundRankHistory(ctx context.Context, client FundClient, endpoint string, code string) (types.FundRankHistory, error) {
	text, err := fetchFundPingzhongText(ctx, client, endpoint, code)
	if err != nil {
		return types.FundRankHistory{}, err
	}
	fundCode, fundName := fundIdentity(text, code)
	var ranks []fundRankItem
	_ = core.ExtractJSVar(text, "Data_rateInSimilarType", &ranks)
	var percentiles [][]float64
	_ = core.ExtractJSVar(text, "Data_rateInSimilarPersent", &percentiles)
	percentileByTime := make(map[int64]float64, len(percentiles))
	for _, row := range percentiles {
		if len(row) >= 2 {
			percentileByTime[int64(row[0])] = row[1]
		}
	}
	items := make([]types.FundRankPoint, 0, len(ranks))
	for _, point := range ranks {
		var percentile *float64
		var date string
		if point.Timestamp != nil {
			if value, ok := percentileByTime[*point.Timestamp]; ok {
				value := value
				percentile = &value
			}
			date = timestampDate(*point.Timestamp)
		}
		items = append(items, types.FundRankPoint{
			Date:       date,
			Timestamp:  point.Timestamp,
			Rank:       nullableNumberAny(point.Rank),
			Total:      nullableNumberAny(point.Total),
			Percentile: percentile,
		})
	}
	return types.FundRankHistory{Code: fundCode, Name: fundName, Items: items}, nil
}

func fetchFundPingzhongText(ctx context.Context, client FundClient, endpoint string, code string) (string, error) {
	requestURL := strings.TrimRight(endpoint, "/") + "/" + url.PathEscape(code) + ".js"
	return client.GetText(ctx, requestURL)
}

func fundIdentity(text string, fallbackCode string) (string, *string) {
	code := fallbackCode
	_ = core.ExtractJSVar(text, "fS_code", &code)
	var name string
	_ = core.ExtractJSVar(text, "fS_name", &name)
	return strings.TrimSpace(code), nullableString(name)
}

func parseFundEstimate(raw fundGZPayload, fallbackCode string) types.FundEstimate {
	code := strings.TrimSpace(raw.FundCode)
	if code == "" {
		code = fallbackCode
	}
	return types.FundEstimate{
		Code:                   code,
		Name:                   nullableString(raw.Name),
		NavDate:                nullableString(raw.NavDate),
		Nav:                    nullableNumber(raw.Nav),
		EstimatedNav:           nullableNumber(raw.Estimate),
		EstimatedChangePercent: nullableNumber(raw.Change),
		EstimateTime:           nullableString(raw.Time),
	}
}

func nullableString(text string) *string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	return &text
}

func nullableNumber(text string) *float64 {
	text = strings.TrimSpace(text)
	if text == "" || text == "--" {
		return nil
	}
	return toNumberFromAny(text)
}

func nullableNumberAny(value any) *float64 {
	switch typed := value.(type) {
	case nil:
		return nil
	case string:
		return nullableNumber(typed)
	case float64:
		return &typed
	case int:
		number := float64(typed)
		return &number
	case int64:
		number := float64(typed)
		return &number
	default:
		return nil
	}
}

func timestampDate(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.UnixMilli(timestamp).UTC().Format("2006-01-02")
}

func toInt64(value any) *int64 {
	switch typed := value.(type) {
	case nil:
		return nil
	case int64:
		return &typed
	case int:
		number := int64(typed)
		return &number
	case float64:
		number := int64(typed)
		return &number
	case string:
		if typed == "" {
			return nil
		}
		parsed, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return nil
		}
		return &parsed
	default:
		return nil
	}
}
