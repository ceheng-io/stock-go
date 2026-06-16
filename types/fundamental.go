package types

// StockIssueInfo contains listing and IPO information for an A-share company.
type StockIssueInfo struct {
	FoundDate        *string  `json:"foundDate"`
	ListingDate      *string  `json:"listingDate"`
	IssueWay         string   `json:"issueWay"`
	ParValue         *float64 `json:"parValue"`
	TotalIssueShares *float64 `json:"totalIssueShares"`
	IssuePrice       *float64 `json:"issuePrice"`
	TotalFunds       *float64 `json:"totalFunds"`
	NetRaiseFunds    *float64 `json:"netRaiseFunds"`
	OpenPrice        *float64 `json:"openPrice"`
	ClosePrice       *float64 `json:"closePrice"`
	TurnoverRate     *float64 `json:"turnoverRate"`
}

// StockProfile contains company profile and issue information.
type StockProfile struct {
	SecuCode                 string         `json:"secuCode"`
	Code                     string         `json:"code"`
	Name                     string         `json:"name"`
	OrgCode                  string         `json:"orgCode"`
	OrgName                  string         `json:"orgName"`
	OrgNameEN                string         `json:"orgNameEn"`
	FormerName               string         `json:"formerName"`
	ACode                    string         `json:"aCode"`
	AName                    string         `json:"aName"`
	HCode                    string         `json:"hCode"`
	HName                    string         `json:"hName"`
	SecurityType             string         `json:"securityType"`
	Industry                 string         `json:"industry"`
	TradeMarket              string         `json:"tradeMarket"`
	CSRCIndustry             string         `json:"csrcIndustry"`
	President                string         `json:"president"`
	LegalRepresentative      string         `json:"legalRepresentative"`
	Secretary                string         `json:"secretary"`
	Chairman                 string         `json:"chairman"`
	SecuritiesRepresentative string         `json:"securitiesRepresentative"`
	IndependentDirectors     string         `json:"independentDirectors"`
	Tel                      string         `json:"tel"`
	Email                    string         `json:"email"`
	Fax                      string         `json:"fax"`
	Website                  string         `json:"website"`
	Address                  string         `json:"address"`
	RegisteredAddress        string         `json:"registeredAddress"`
	Province                 string         `json:"province"`
	AddressPostcode          string         `json:"addressPostcode"`
	RegisteredCapital        *float64       `json:"registeredCapital"`
	RegistrationNumber       string         `json:"registrationNumber"`
	EmployeeCount            *float64       `json:"employeeCount"`
	LawFirm                  string         `json:"lawFirm"`
	AccountingFirm           string         `json:"accountingFirm"`
	Profile                  string         `json:"profile"`
	BusinessScope            string         `json:"businessScope"`
	Issue                    StockIssueInfo `json:"issue"`
}

// FinancialReportPeriod identifies the financial indicator report scope.
type FinancialReportPeriod string

const (
	FinancialReportPeriodAll    FinancialReportPeriod = "all"
	FinancialReportPeriodAnnual FinancialReportPeriod = "annual"
)

// FinancialIndicatorOptions configures stock financial indicator requests.
type FinancialIndicatorOptions struct {
	Period FinancialReportPeriod `json:"period"`
}

// FinancialIndicator contains one stock financial indicator row.
type FinancialIndicator struct {
	SecuCode                  string   `json:"secuCode"`
	Code                      string   `json:"code"`
	Name                      string   `json:"name"`
	ReportDate                *string  `json:"reportDate"`
	ReportType                string   `json:"reportType"`
	ReportDateName            string   `json:"reportDateName"`
	NoticeDate                *string  `json:"noticeDate"`
	UpdateDate                *string  `json:"updateDate"`
	Currency                  string   `json:"currency"`
	BasicEPS                  *float64 `json:"basicEps"`
	DeductBasicEPS            *float64 `json:"deductBasicEps"`
	DilutedEPS                *float64 `json:"dilutedEps"`
	BPS                       *float64 `json:"bps"`
	CapitalReservePerShare    *float64 `json:"capitalReservePerShare"`
	UnassignedProfitPerShare  *float64 `json:"unassignedProfitPerShare"`
	OperatingCashFlowPerShare *float64 `json:"operatingCashFlowPerShare"`
	TotalRevenue              *float64 `json:"totalRevenue"`
	GrossProfit               *float64 `json:"grossProfit"`
	ParentNetProfit           *float64 `json:"parentNetProfit"`
	DeductParentNetProfit     *float64 `json:"deductParentNetProfit"`
	TotalRevenueYoY           *float64 `json:"totalRevenueYoY"`
	ParentNetProfitYoY        *float64 `json:"parentNetProfitYoY"`
	DeductParentNetProfitYoY  *float64 `json:"deductParentNetProfitYoY"`
	ROEWeighted               *float64 `json:"roeWeighted"`
	ROEDeductWeighted         *float64 `json:"roeDeductWeighted"`
	ROA                       *float64 `json:"roa"`
	NetMargin                 *float64 `json:"netMargin"`
	GrossMargin               *float64 `json:"grossMargin"`
	AssetLiabilityRatio       *float64 `json:"assetLiabilityRatio"`
	ROIC                      *float64 `json:"roic"`
	StaffCount                *float64 `json:"staffCount"`
}

// AnnouncementOptions configures stock announcement list requests.
type AnnouncementOptions struct {
	PageSize  int `json:"pageSize"`
	PageIndex int `json:"pageIndex"`
}

// StockAnnouncementCode identifies a security mentioned by an announcement.
type StockAnnouncementCode struct {
	AnnouncementType string `json:"announcementType"`
	InnerCode        string `json:"innerCode"`
	MarketCode       string `json:"marketCode"`
	ShortName        string `json:"shortName"`
	StockCode        string `json:"stockCode"`
}

// StockAnnouncementColumn identifies the announcement category.
type StockAnnouncementColumn struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// StockAnnouncement contains one stock announcement list row.
type StockAnnouncement struct {
	ArtCode     string                    `json:"artCode"`
	Title       string                    `json:"title"`
	TitleCH     string                    `json:"titleCh"`
	TitleEN     string                    `json:"titleEn"`
	NoticeDate  *string                   `json:"noticeDate"`
	DisplayTime *string                   `json:"displayTime"`
	SortDate    *string                   `json:"sortDate"`
	Columns     []StockAnnouncementColumn `json:"columns"`
	Codes       []StockAnnouncementCode   `json:"codes"`
}

// StockAnnouncementResult contains paginated stock announcement rows.
type StockAnnouncementResult struct {
	List      []StockAnnouncement `json:"list"`
	PageIndex int                 `json:"pageIndex"`
	PageSize  int                 `json:"pageSize"`
	Total     int                 `json:"total"`
}

// StockAnnouncementAttachment contains one announcement attachment.
type StockAnnouncementAttachment struct {
	URL  string   `json:"url"`
	Type string   `json:"type"`
	Size *float64 `json:"size"`
	Seq  *float64 `json:"seq"`
}

// StockAnnouncementDetail contains announcement body text and attachments.
type StockAnnouncementDetail struct {
	ArtCode       string                        `json:"artCode"`
	Title         string                        `json:"title"`
	NoticeDate    *string                       `json:"noticeDate"`
	AttachURL     string                        `json:"attachUrl"`
	AttachURLWeb  string                        `json:"attachUrlWeb"`
	AttachSize    string                        `json:"attachSize"`
	AttachType    string                        `json:"attachType"`
	NoticeContent string                        `json:"noticeContent"`
	Attachments   []StockAnnouncementAttachment `json:"attachments"`
}
