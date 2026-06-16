package types

// DividendDetail is a stock dividend and bonus issue detail row.
type DividendDetail struct {
	Code                string   `json:"code"`
	Name                string   `json:"name"`
	ReportDate          *string  `json:"reportDate"`
	PlanNoticeDate      *string  `json:"planNoticeDate"`
	DisclosureDate      *string  `json:"disclosureDate"`
	AssignTransferRatio *float64 `json:"assignTransferRatio"`
	BonusRatio          *float64 `json:"bonusRatio"`
	TransferRatio       *float64 `json:"transferRatio"`
	DividendPretax      *float64 `json:"dividendPretax"`
	DividendDesc        *string  `json:"dividendDesc"`
	DividendYield       *float64 `json:"dividendYield"`
	EPS                 *float64 `json:"eps"`
	BPS                 *float64 `json:"bps"`
	CapitalReserve      *float64 `json:"capitalReserve"`
	UnassignedProfit    *float64 `json:"unassignedProfit"`
	NetProfitYoY        *float64 `json:"netProfitYoY"`
	TotalShares         *float64 `json:"totalShares"`
	EquityRecordDate    *string  `json:"equityRecordDate"`
	ExDividendDate      *string  `json:"exDividendDate"`
	PayDate             *string  `json:"payDate"`
	AssignProgress      *string  `json:"assignProgress"`
	NoticeDate          *string  `json:"noticeDate"`
}
