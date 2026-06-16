package types

// StockIssueInfo contains listing and IPO information for an A-share company.
type StockIssueInfo struct {
	FoundDate        *string
	ListingDate      *string
	IssueWay         string
	ParValue         *float64
	TotalIssueShares *float64
	IssuePrice       *float64
	TotalFunds       *float64
	NetRaiseFunds    *float64
	OpenPrice        *float64
	ClosePrice       *float64
	TurnoverRate     *float64
}

// StockProfile contains company profile and issue information.
type StockProfile struct {
	SecuCode                 string
	Code                     string
	Name                     string
	OrgCode                  string
	OrgName                  string
	OrgNameEN                string
	FormerName               string
	ACode                    string
	AName                    string
	HCode                    string
	HName                    string
	SecurityType             string
	Industry                 string
	TradeMarket              string
	CSRCIndustry             string
	President                string
	LegalRepresentative      string
	Secretary                string
	Chairman                 string
	SecuritiesRepresentative string
	IndependentDirectors     string
	Tel                      string
	Email                    string
	Fax                      string
	Website                  string
	Address                  string
	RegisteredAddress        string
	Province                 string
	AddressPostcode          string
	RegisteredCapital        *float64
	RegistrationNumber       string
	EmployeeCount            *float64
	LawFirm                  string
	AccountingFirm           string
	Profile                  string
	BusinessScope            string
	Issue                    StockIssueInfo
}

// FinancialReportPeriod identifies the financial indicator report scope.
type FinancialReportPeriod string

const (
	FinancialReportPeriodAll    FinancialReportPeriod = "all"
	FinancialReportPeriodAnnual FinancialReportPeriod = "annual"
)

// FinancialIndicatorOptions configures stock financial indicator requests.
type FinancialIndicatorOptions struct {
	Period FinancialReportPeriod
}

// FinancialIndicator contains one stock financial indicator row.
type FinancialIndicator struct {
	SecuCode                  string
	Code                      string
	Name                      string
	ReportDate                *string
	ReportType                string
	ReportDateName            string
	NoticeDate                *string
	UpdateDate                *string
	Currency                  string
	BasicEPS                  *float64
	DeductBasicEPS            *float64
	DilutedEPS                *float64
	BPS                       *float64
	CapitalReservePerShare    *float64
	UnassignedProfitPerShare  *float64
	OperatingCashFlowPerShare *float64
	TotalRevenue              *float64
	GrossProfit               *float64
	ParentNetProfit           *float64
	DeductParentNetProfit     *float64
	TotalRevenueYoY           *float64
	ParentNetProfitYoY        *float64
	DeductParentNetProfitYoY  *float64
	ROEWeighted               *float64
	ROEDeductWeighted         *float64
	ROA                       *float64
	NetMargin                 *float64
	GrossMargin               *float64
	AssetLiabilityRatio       *float64
	ROIC                      *float64
	StaffCount                *float64
}

// AnnouncementOptions configures stock announcement list requests.
type AnnouncementOptions struct {
	PageSize  int
	PageIndex int
}

// StockAnnouncementCode identifies a security mentioned by an announcement.
type StockAnnouncementCode struct {
	AnnouncementType string
	InnerCode        string
	MarketCode       string
	ShortName        string
	StockCode        string
}

// StockAnnouncementColumn identifies the announcement category.
type StockAnnouncementColumn struct {
	Code string
	Name string
}

// StockAnnouncement contains one stock announcement list row.
type StockAnnouncement struct {
	ArtCode     string
	Title       string
	TitleCH     string
	TitleEN     string
	NoticeDate  *string
	DisplayTime *string
	SortDate    *string
	Columns     []StockAnnouncementColumn
	Codes       []StockAnnouncementCode
}

// StockAnnouncementResult contains paginated stock announcement rows.
type StockAnnouncementResult struct {
	List      []StockAnnouncement
	PageIndex int
	PageSize  int
	Total     int
}

// StockAnnouncementAttachment contains one announcement attachment.
type StockAnnouncementAttachment struct {
	URL  string
	Type string
	Size *float64
	Seq  *float64
}

// StockAnnouncementDetail contains announcement body text and attachments.
type StockAnnouncementDetail struct {
	ArtCode       string
	Title         string
	NoticeDate    *string
	AttachURL     string
	AttachURLWeb  string
	AttachSize    string
	AttachType    string
	NoticeContent string
	Attachments   []StockAnnouncementAttachment
}
