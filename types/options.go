package types

// IndexOptionProduct 是中金所股指期权产品。
type IndexOptionProduct string

const (
	IndexOptionProductHO IndexOptionProduct = "ho"
	IndexOptionProductIO IndexOptionProduct = "io"
	IndexOptionProductMO IndexOptionProduct = "mo"
)

// ETFOptionCate 是新浪 ETF 期权分类。
type ETFOptionCate string

const (
	ETFOptionCate50ETF           ETFOptionCate = "50ETF"
	ETFOptionCate300ETF          ETFOptionCate = "300ETF"
	ETFOptionCate500ETF          ETFOptionCate = "500ETF"
	ETFOptionCateKechuang50      ETFOptionCate = "科创50"
	ETFOptionCateKechuangBoard50 ETFOptionCate = "科创板50"
)

// ETFOptionMonth 是新浪 ETF 期权可用月份信息。
type ETFOptionMonth struct {
	Months   []string `json:"months"`
	StockID  string   `json:"stockId"`
	CateID   string   `json:"cateId"`
	CateList []string `json:"cateList"`
}

// ETFOptionExpireDay 是新浪 ETF 期权到期日信息。
type ETFOptionExpireDay struct {
	ExpireDay     string `json:"expireDay"`
	RemainderDays int    `json:"remainderDays"`
	StockID       string `json:"stockId"`
	Name          string `json:"name"`
}

// OptionKline 是期权日 K 线。
type OptionKline struct {
	Date   string   `json:"date"`
	Open   *float64 `json:"open"`
	High   *float64 `json:"high"`
	Low    *float64 `json:"low"`
	Close  *float64 `json:"close"`
	Volume *float64 `json:"volume"`
}

// OptionMinute 是期权分钟行情。
type OptionMinute struct {
	Time         string   `json:"time"`
	Date         string   `json:"date"`
	Price        *float64 `json:"price"`
	Volume       *float64 `json:"volume"`
	OpenInterest *float64 `json:"openInterest"`
	AvgPrice     *float64 `json:"avgPrice"`
}

// OptionTQuote 是期权 T 型报价项。
type OptionTQuote struct {
	Symbol       string   `json:"symbol"`
	BuyVolume    *float64 `json:"buyVolume"`
	BuyPrice     *float64 `json:"buyPrice"`
	Price        *float64 `json:"price"`
	AskPrice     *float64 `json:"askPrice"`
	AskVolume    *float64 `json:"askVolume"`
	OpenInterest *float64 `json:"openInterest"`
	Change       *float64 `json:"change"`
	StrikePrice  *float64 `json:"strikePrice"`
}

// OptionTQuoteResult 是期权 T 型报价结果。
type OptionTQuoteResult struct {
	Calls []OptionTQuote `json:"calls"`
	Puts  []OptionTQuote `json:"puts"`
}

// CFFEXOptionQuote 是中金所期权实时行情。
type CFFEXOptionQuote struct {
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Price         *float64 `json:"price"`
	Change        *float64 `json:"change"`
	ChangePercent *float64 `json:"changePercent"`
	Volume        *float64 `json:"volume"`
	Amount        *float64 `json:"amount"`
	OpenInterest  *float64 `json:"openInterest"`
	StrikePrice   *float64 `json:"strikePrice"`
	RemainDays    *float64 `json:"remainDays"`
	DailyChange   *float64 `json:"dailyChange"`
	PrevSettle    *float64 `json:"prevSettle"`
	Open          *float64 `json:"open"`
}

// OptionLHBItem 是期权龙虎榜数据行。
type OptionLHBItem struct {
	TradeType          string   `json:"tradeType"`
	Date               string   `json:"date"`
	Symbol             string   `json:"symbol"`
	TargetName         string   `json:"targetName"`
	Rank               int      `json:"rank"`
	MemberName         string   `json:"memberName"`
	SellVolume         *float64 `json:"sellVolume"`
	SellVolumeChange   *float64 `json:"sellVolumeChange"`
	NetSellVolume      *float64 `json:"netSellVolume"`
	SellVolumeRatio    *float64 `json:"sellVolumeRatio"`
	BuyVolume          *float64 `json:"buyVolume"`
	BuyVolumeChange    *float64 `json:"buyVolumeChange"`
	NetBuyVolume       *float64 `json:"netBuyVolume"`
	BuyVolumeRatio     *float64 `json:"buyVolumeRatio"`
	SellPosition       *float64 `json:"sellPosition"`
	SellPositionChange *float64 `json:"sellPositionChange"`
	NetSellPosition    *float64 `json:"netSellPosition"`
	SellPositionRatio  *float64 `json:"sellPositionRatio"`
	BuyPosition        *float64 `json:"buyPosition"`
	BuyPositionChange  *float64 `json:"buyPositionChange"`
	NetBuyPosition     *float64 `json:"netBuyPosition"`
	BuyPositionRatio   *float64 `json:"buyPositionRatio"`
}
