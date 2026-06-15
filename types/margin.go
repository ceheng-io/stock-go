package types

// MarginAccountItem is a daily margin account statistics row.
type MarginAccountItem struct {
	Date                   string
	FinBalance             *float64
	LoanBalance            *float64
	FinBuyAmount           *float64
	LoanSellAmount         *float64
	InvestorCount          *float64
	LiabilityInvestorCount *float64
	TotalGuarantee         *float64
	AvgGuaranteeRatio      *float64
}

// MarginTargetItem is a stock margin target detail row.
type MarginTargetItem struct {
	Code            string
	Name            string
	Date            string
	FinBalance      *float64
	FinBuyAmount    *float64
	FinRepayAmount  *float64
	LoanBalance     *float64
	LoanSellVolume  *float64
	LoanRepayVolume *float64
}
