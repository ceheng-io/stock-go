package types

// MarginAccountItem is a daily margin account statistics row.
type MarginAccountItem struct {
	Date                   string   `json:"date"`
	FinBalance             *float64 `json:"finBalance"`
	LoanBalance            *float64 `json:"loanBalance"`
	FinBuyAmount           *float64 `json:"finBuyAmount"`
	LoanSellAmount         *float64 `json:"loanSellAmount"`
	InvestorCount          *float64 `json:"investorCount"`
	LiabilityInvestorCount *float64 `json:"liabilityInvestorCount"`
	TotalGuarantee         *float64 `json:"totalGuarantee"`
	AvgGuaranteeRatio      *float64 `json:"avgGuaranteeRatio"`
}

// MarginTargetItem is a stock margin target detail row.
type MarginTargetItem struct {
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	Date            string   `json:"date"`
	FinBalance      *float64 `json:"finBalance"`
	FinBuyAmount    *float64 `json:"finBuyAmount"`
	FinRepayAmount  *float64 `json:"finRepayAmount"`
	LoanBalance     *float64 `json:"loanBalance"`
	LoanSellVolume  *float64 `json:"loanSellVolume"`
	LoanRepayVolume *float64 `json:"loanRepayVolume"`
}
