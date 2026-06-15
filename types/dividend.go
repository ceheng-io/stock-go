package types

// DividendDetail is a stock dividend and bonus issue detail row.
type DividendDetail struct {
	Code                string
	Name                string
	ReportDate          *string
	PlanNoticeDate      *string
	DisclosureDate      *string
	AssignTransferRatio *float64
	BonusRatio          *float64
	TransferRatio       *float64
	DividendPretax      *float64
	DividendDesc        *string
	DividendYield       *float64
	EPS                 *float64
	BPS                 *float64
	CapitalReserve      *float64
	UnassignedProfit    *float64
	NetProfitYoY        *float64
	TotalShares         *float64
	EquityRecordDate    *string
	ExDividendDate      *string
	PayDate             *string
	AssignProgress      *string
	NoticeDate          *string
}
