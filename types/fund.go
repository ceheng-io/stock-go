package types

// FundDividendRank identifies a fund dividend list sort column.
type FundDividendRank string

const (
	FundDividendRankCode             FundDividendRank = "BZDM"
	FundDividendRankName             FundDividendRank = "ABBNAME"
	FundDividendRankEquityRecordDate FundDividendRank = "DJR"
	FundDividendRankExDividendDate   FundDividendRank = "FSRQ"
	FundDividendRankDividendPerShare FundDividendRank = "FHFCZ"
	FundDividendRankPayDate          FundDividendRank = "FFR"
)

// FundSortDirection identifies a fund list sort direction.
type FundSortDirection string

const (
	FundSortAsc  FundSortDirection = "asc"
	FundSortDesc FundSortDirection = "desc"
)

// FundDividendListOptions configures public fund dividend list requests.
type FundDividendListOptions struct {
	Year     string            `json:"year"`
	Page     any               `json:"page"`
	AllPages bool              `json:"allPages"`
	FundType string            `json:"fundType"`
	Rank     FundDividendRank  `json:"rank"`
	Sort     FundSortDirection `json:"sort"`
	Code     string            `json:"code"`
}

// FundEstimate is a public fund's latest net-value estimate.
type FundEstimate struct {
	Code                   string   `json:"code"`
	Name                   *string  `json:"name"`
	NavDate                *string  `json:"navDate"`
	Nav                    *float64 `json:"nav"`
	EstimatedNav           *float64 `json:"estimatedNav"`
	EstimatedChangePercent *float64 `json:"estimatedChangePercent"`
	EstimateTime           *string  `json:"estimateTime"`
}

// FundNavPoint is one public fund net-value history point.
type FundNavPoint struct {
	Date        string   `json:"date"`
	Timestamp   *int64   `json:"timestamp"`
	Nav         float64  `json:"nav"`
	AccNav      *float64 `json:"accNav"`
	DailyReturn *float64 `json:"dailyReturn"`
	UnitMoney   string   `json:"unitMoney"`
}

// FundNavHistory contains a public fund's net-value history.
type FundNavHistory struct {
	Code  string         `json:"code"`
	Name  *string        `json:"name"`
	Items []FundNavPoint `json:"items"`
}

// FundRankPoint is one public fund similar-type rank history point.
type FundRankPoint struct {
	Date       string   `json:"date"`
	Timestamp  *int64   `json:"timestamp"`
	Rank       *float64 `json:"rank"`
	Total      *float64 `json:"total"`
	Percentile *float64 `json:"percentile"`
}

// FundRankHistory contains a public fund's similar-type rank history.
type FundRankHistory struct {
	Code  string          `json:"code"`
	Name  *string         `json:"name"`
	Items []FundRankPoint `json:"items"`
}

// FundDividend is one public fund dividend distribution row.
type FundDividend struct {
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	EquityRecordDate *string  `json:"equityRecordDate"`
	ExDividendDate   *string  `json:"exDividendDate"`
	DividendPerShare *float64 `json:"dividendPerShare"`
	PayDate          *string  `json:"payDate"`
}

// FundDividendListResult contains public fund dividend distribution rows.
type FundDividendListResult struct {
	Items       []FundDividend `json:"items"`
	TotalPages  int            `json:"totalPages"`
	PageSize    int            `json:"pageSize"`
	CurrentPage int            `json:"currentPage"`
}
