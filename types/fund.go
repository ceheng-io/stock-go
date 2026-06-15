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
	Year     string
	Page     any
	AllPages bool
	FundType string
	Rank     FundDividendRank
	Sort     FundSortDirection
	Code     string
}

// FundEstimate is a public fund's latest net-value estimate.
type FundEstimate struct {
	Code                   string
	Name                   *string
	NavDate                *string
	Nav                    *float64
	EstimatedNav           *float64
	EstimatedChangePercent *float64
	EstimateTime           *string
}

// FundNavPoint is one public fund net-value history point.
type FundNavPoint struct {
	Date        string
	Timestamp   *int64
	Nav         float64
	AccNav      *float64
	DailyReturn *float64
	UnitMoney   string
}

// FundNavHistory contains a public fund's net-value history.
type FundNavHistory struct {
	Code  string
	Name  *string
	Items []FundNavPoint
}

// FundRankPoint is one public fund similar-type rank history point.
type FundRankPoint struct {
	Date       string
	Timestamp  *int64
	Rank       *float64
	Total      *float64
	Percentile *float64
}

// FundRankHistory contains a public fund's similar-type rank history.
type FundRankHistory struct {
	Code  string
	Name  *string
	Items []FundRankPoint
}

// FundDividend is one public fund dividend distribution row.
type FundDividend struct {
	Code             string
	Name             string
	EquityRecordDate *string
	ExDividendDate   *string
	DividendPerShare *float64
	PayDate          *string
}

// FundDividendListResult contains public fund dividend distribution rows.
type FundDividendListResult struct {
	Items       []FundDividend
	TotalPages  int
	PageSize    int
	CurrentPage int
}
