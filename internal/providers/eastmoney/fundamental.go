package eastmoney

import (
	"context"
	"net/url"
	"strings"

	"github.com/ceheng-io/stock-go/types"
)

const (
	stockProfilePath            = "/PC_HSF10/CompanySurvey/PageAjax"
	financialIndicatorsPath     = "/PC_HSF10/NewFinanceAnalysis/ZYZBAjaxNew"
	announcementListDefaultSize = 20
)

// FundamentalClient is the minimal client interface required by stock fundamental providers.
type FundamentalClient interface {
	GetJSON(context.Context, string, any) error
}

type stockProfileResponse struct {
	Basic []boardDynamicItem `json:"jbzl"`
	Issue []boardDynamicItem `json:"fxxg"`
}

type financialIndicatorsResponse struct {
	Pages int                `json:"pages"`
	Data  []boardDynamicItem `json:"data"`
}

type announcementListResponse struct {
	Data struct {
		List      []announcementItem `json:"list"`
		PageIndex int                `json:"page_index"`
		PageSize  int                `json:"page_size"`
		TotalHits int                `json:"total_hits"`
	} `json:"data"`
}

type announcementItem struct {
	ArtCode     string                      `json:"art_code"`
	Title       string                      `json:"title"`
	TitleCH     string                      `json:"title_ch"`
	TitleEN     string                      `json:"title_en"`
	NoticeDate  string                      `json:"notice_date"`
	DisplayTime string                      `json:"display_time"`
	SortDate    string                      `json:"sort_date"`
	Columns     []announcementColumnPayload `json:"columns"`
	Codes       []announcementCodePayload   `json:"codes"`
}

type announcementColumnPayload struct {
	Code string `json:"column_code"`
	Name string `json:"column_name"`
}

type announcementCodePayload struct {
	AnnouncementType string `json:"ann_type"`
	InnerCode        string `json:"inner_code"`
	MarketCode       string `json:"market_code"`
	ShortName        string `json:"short_name"`
	StockCode        string `json:"stock_code"`
}

type announcementDetailResponse struct {
	Data announcementDetailPayload `json:"data"`
}

type announcementDetailPayload struct {
	ArtCode       string                      `json:"art_code"`
	Title         string                      `json:"title"`
	NoticeDate    string                      `json:"notice_date"`
	AttachURL     string                      `json:"attach_url"`
	AttachURLWeb  string                      `json:"attach_url_web"`
	AttachSize    string                      `json:"attach_size"`
	AttachType    string                      `json:"attach_type"`
	NoticeContent string                      `json:"notice_content"`
	AttachList    []announcementAttachmentRaw `json:"attach_list"`
}

type announcementAttachmentRaw struct {
	URL  string `json:"attach_url"`
	Type string `json:"attach_type"`
	Size any    `json:"attach_size"`
	Seq  any    `json:"seq"`
}

// GetStockProfile fetches Eastmoney F10 company profile data.
func GetStockProfile(ctx context.Context, client FundamentalClient, endpoint string, symbol string) (types.StockProfile, error) {
	params := url.Values{}
	params.Set("code", f10Code(symbol))
	var payload stockProfileResponse
	if err := client.GetJSON(ctx, strings.TrimRight(endpoint, "/")+stockProfilePath+"?"+params.Encode(), &payload); err != nil {
		return types.StockProfile{}, err
	}
	basic := firstDynamicItem(payload.Basic)
	issue := firstDynamicItem(payload.Issue)
	return parseStockProfile(basic, issue), nil
}

// GetFinancialIndicators fetches stock financial indicator rows.
func GetFinancialIndicators(ctx context.Context, client FundamentalClient, endpoint string, symbol string, options types.FinancialIndicatorOptions) ([]types.FinancialIndicator, error) {
	params := url.Values{}
	params.Set("type", financialIndicatorType(options.Period))
	params.Set("code", f10Code(symbol))
	var payload financialIndicatorsResponse
	if err := client.GetJSON(ctx, strings.TrimRight(endpoint, "/")+financialIndicatorsPath+"?"+params.Encode(), &payload); err != nil {
		return nil, err
	}
	rows := make([]types.FinancialIndicator, 0, len(payload.Data))
	for _, item := range payload.Data {
		rows = append(rows, parseFinancialIndicator(item))
	}
	return rows, nil
}

// GetStockAnnouncements fetches paginated Eastmoney stock announcements.
func GetStockAnnouncements(ctx context.Context, client FundamentalClient, endpoint string, symbol string, options types.AnnouncementOptions) (types.StockAnnouncementResult, error) {
	params := url.Values{}
	params.Set("sr", "-1")
	params.Set("page_size", intString(defaultPositive(options.PageSize, announcementListDefaultSize)))
	params.Set("page_index", intString(defaultPositive(options.PageIndex, 1)))
	params.Set("ann_type", "A")
	params.Set("client_source", "web")
	params.Set("stock_list", pureAshareCode(symbol))
	var payload announcementListResponse
	if err := client.GetJSON(ctx, strings.TrimRight(endpoint, "/")+"?"+params.Encode(), &payload); err != nil {
		return types.StockAnnouncementResult{}, err
	}
	rows := make([]types.StockAnnouncement, 0, len(payload.Data.List))
	for _, item := range payload.Data.List {
		rows = append(rows, parseAnnouncement(item))
	}
	return types.StockAnnouncementResult{
		List:      rows,
		PageIndex: payload.Data.PageIndex,
		PageSize:  payload.Data.PageSize,
		Total:     payload.Data.TotalHits,
	}, nil
}

// GetStockAnnouncementDetail fetches one Eastmoney announcement body and attachments.
func GetStockAnnouncementDetail(ctx context.Context, client FundamentalClient, endpoint string, artCode string) (types.StockAnnouncementDetail, error) {
	params := url.Values{}
	params.Set("art_code", strings.TrimSpace(artCode))
	params.Set("client_source", "web")
	var payload announcementDetailResponse
	if err := client.GetJSON(ctx, strings.TrimRight(endpoint, "/")+"?"+params.Encode(), &payload); err != nil {
		return types.StockAnnouncementDetail{}, err
	}
	return parseAnnouncementDetail(payload.Data), nil
}

func parseStockProfile(basic boardDynamicItem, issue boardDynamicItem) types.StockProfile {
	return types.StockProfile{
		SecuCode:                 stringValue(basic["SECUCODE"]),
		Code:                     stringValue(basic["SECURITY_CODE"]),
		Name:                     stringValue(basic["SECURITY_NAME_ABBR"]),
		OrgCode:                  stringValue(basic["ORG_CODE"]),
		OrgName:                  stringValue(basic["ORG_NAME"]),
		OrgNameEN:                stringValue(basic["ORG_NAME_EN"]),
		FormerName:               stringValue(basic["FORMERNAME"]),
		ACode:                    stringValue(basic["STR_CODEA"]),
		AName:                    stringValue(basic["STR_NAMEA"]),
		HCode:                    stringValue(basic["STR_CODEH"]),
		HName:                    stringValue(basic["STR_NAMEH"]),
		SecurityType:             stringValue(basic["SECURITY_TYPE"]),
		Industry:                 stringValue(basic["EM2016"]),
		TradeMarket:              stringValue(basic["TRADE_MARKET"]),
		CSRCIndustry:             stringValue(basic["INDUSTRYCSRC1"]),
		President:                stringValue(basic["PRESIDENT"]),
		LegalRepresentative:      stringValue(basic["LEGAL_PERSON"]),
		Secretary:                stringValue(basic["SECRETARY"]),
		Chairman:                 stringValue(basic["CHAIRMAN"]),
		SecuritiesRepresentative: stringValue(basic["SECPRESENT"]),
		IndependentDirectors:     stringValue(basic["INDEDIRECTORS"]),
		Tel:                      stringValue(basic["ORG_TEL"]),
		Email:                    stringValue(basic["ORG_EMAIL"]),
		Fax:                      stringValue(basic["ORG_FAX"]),
		Website:                  stringValue(basic["ORG_WEB"]),
		Address:                  stringValue(basic["ADDRESS"]),
		RegisteredAddress:        stringValue(basic["REG_ADDRESS"]),
		Province:                 stringValue(basic["PROVINCE"]),
		AddressPostcode:          stringValue(basic["ADDRESS_POSTCODE"]),
		RegisteredCapital:        toNumberFromAny(basic["REG_CAPITAL"]),
		RegistrationNumber:       stringValue(basic["REG_NUM"]),
		EmployeeCount:            toNumberFromAny(basic["EMP_NUM"]),
		LawFirm:                  stringValue(basic["LAW_FIRM"]),
		AccountingFirm:           stringValue(basic["ACCOUNTFIRM_NAME"]),
		Profile:                  strings.TrimSpace(stringValue(basic["ORG_PROFILE"])),
		BusinessScope:            strings.TrimSpace(stringValue(basic["BUSINESS_SCOPE"])),
		Issue:                    parseStockIssueInfo(issue),
	}
}

func parseStockIssueInfo(item boardDynamicItem) types.StockIssueInfo {
	return types.StockIssueInfo{
		FoundDate:        nullableDatacenterDate(item["FOUND_DATE"]),
		ListingDate:      nullableDatacenterDate(item["LISTING_DATE"]),
		IssueWay:         stringValue(item["ISSUE_WAY"]),
		ParValue:         toNumberFromAny(item["PAR_VALUE"]),
		TotalIssueShares: toNumberFromAny(item["TOTAL_ISSUE_NUM"]),
		IssuePrice:       toNumberFromAny(item["ISSUE_PRICE"]),
		TotalFunds:       toNumberFromAny(item["TOTAL_FUNDS"]),
		NetRaiseFunds:    toNumberFromAny(item["NET_RAISE_FUNDS"]),
		OpenPrice:        toNumberFromAny(item["OPEN_PRICE"]),
		ClosePrice:       toNumberFromAny(item["CLOSE_PRICE"]),
		TurnoverRate:     toNumberFromAny(item["TURNOVERRATE"]),
	}
}

func parseFinancialIndicator(item boardDynamicItem) types.FinancialIndicator {
	return types.FinancialIndicator{
		SecuCode:                  stringValue(item["SECUCODE"]),
		Code:                      stringValue(item["SECURITY_CODE"]),
		Name:                      stringValue(item["SECURITY_NAME_ABBR"]),
		ReportDate:                nullableDatacenterDate(item["REPORT_DATE"]),
		ReportType:                stringValue(item["REPORT_TYPE"]),
		ReportDateName:            stringValue(item["REPORT_DATE_NAME"]),
		NoticeDate:                nullableDatacenterDate(item["NOTICE_DATE"]),
		UpdateDate:                nullableDatacenterDate(item["UPDATE_DATE"]),
		Currency:                  stringValue(item["CURRENCY"]),
		BasicEPS:                  toNumberFromAny(item["EPSJB"]),
		DeductBasicEPS:            toNumberFromAny(item["EPSKCJB"]),
		DilutedEPS:                toNumberFromAny(item["EPSXS"]),
		BPS:                       toNumberFromAny(item["BPS"]),
		CapitalReservePerShare:    toNumberFromAny(item["MGZBGJ"]),
		UnassignedProfitPerShare:  toNumberFromAny(item["MGWFPLR"]),
		OperatingCashFlowPerShare: toNumberFromAny(item["MGJYXJJE"]),
		TotalRevenue:              toNumberFromAny(item["TOTALOPERATEREVE"]),
		GrossProfit:               toNumberFromAny(item["MLR"]),
		ParentNetProfit:           toNumberFromAny(item["PARENTNETPROFIT"]),
		DeductParentNetProfit:     toNumberFromAny(item["KCFJCXSYJLR"]),
		TotalRevenueYoY:           toNumberFromAny(item["TOTALOPERATEREVETZ"]),
		ParentNetProfitYoY:        toNumberFromAny(item["PARENTNETPROFITTZ"]),
		DeductParentNetProfitYoY:  toNumberFromAny(item["KCFJCXSYJLRTZ"]),
		ROEWeighted:               toNumberFromAny(item["ROEJQ"]),
		ROEDeductWeighted:         toNumberFromAny(item["ROEKCJQ"]),
		ROA:                       toNumberFromAny(item["ZZCJLL"]),
		NetMargin:                 toNumberFromAny(item["XSJLL"]),
		GrossMargin:               toNumberFromAny(item["XSMLL"]),
		AssetLiabilityRatio:       toNumberFromAny(item["ZCFZL"]),
		ROIC:                      toNumberFromAny(item["ROIC"]),
		StaffCount:                toNumberFromAny(item["STAFF_NUM"]),
	}
}

func parseAnnouncement(item announcementItem) types.StockAnnouncement {
	columns := make([]types.StockAnnouncementColumn, 0, len(item.Columns))
	for _, column := range item.Columns {
		columns = append(columns, types.StockAnnouncementColumn{Code: column.Code, Name: column.Name})
	}
	codes := make([]types.StockAnnouncementCode, 0, len(item.Codes))
	for _, code := range item.Codes {
		codes = append(codes, types.StockAnnouncementCode{
			AnnouncementType: code.AnnouncementType,
			InnerCode:        code.InnerCode,
			MarketCode:       code.MarketCode,
			ShortName:        code.ShortName,
			StockCode:        code.StockCode,
		})
	}
	return types.StockAnnouncement{
		ArtCode:     item.ArtCode,
		Title:       item.Title,
		TitleCH:     item.TitleCH,
		TitleEN:     item.TitleEN,
		NoticeDate:  nullableDatacenterDate(item.NoticeDate),
		DisplayTime: nullableDatacenterString(normalizeDateTime(item.DisplayTime)),
		SortDate:    nullableDatacenterString(normalizeDateTime(item.SortDate)),
		Columns:     columns,
		Codes:       codes,
	}
}

func parseAnnouncementDetail(item announcementDetailPayload) types.StockAnnouncementDetail {
	attachments := make([]types.StockAnnouncementAttachment, 0, len(item.AttachList))
	for _, attachment := range item.AttachList {
		attachments = append(attachments, types.StockAnnouncementAttachment{
			URL:  attachment.URL,
			Type: attachment.Type,
			Size: toNumberFromAny(attachment.Size),
			Seq:  toNumberFromAny(attachment.Seq),
		})
	}
	return types.StockAnnouncementDetail{
		ArtCode:       item.ArtCode,
		Title:         item.Title,
		NoticeDate:    nullableDatacenterDate(item.NoticeDate),
		AttachURL:     item.AttachURL,
		AttachURLWeb:  item.AttachURLWeb,
		AttachSize:    item.AttachSize,
		AttachType:    item.AttachType,
		NoticeContent: strings.TrimSpace(item.NoticeContent),
		Attachments:   attachments,
	}
}

func financialIndicatorType(period types.FinancialReportPeriod) string {
	if period == types.FinancialReportPeriodAll {
		return "0"
	}
	return "1"
}

func f10Code(symbol string) string {
	code := pureAshareCode(symbol)
	market := strings.ToUpper(symbolMarket(symbol))
	if market == "" {
		switch {
		case strings.HasPrefix(code, "6"):
			market = "SH"
		case strings.HasPrefix(code, "4"), strings.HasPrefix(code, "8"), strings.HasPrefix(code, "92"):
			market = "BJ"
		default:
			market = "SZ"
		}
	}
	return market + code
}

func pureAshareCode(symbol string) string {
	code := strings.TrimSpace(symbol)
	code = strings.TrimPrefix(code, ".")
	lower := strings.ToLower(code)
	for _, prefix := range []string{"sh", "sz", "bj"} {
		if strings.HasPrefix(lower, prefix) {
			code = code[len(prefix):]
			code = strings.TrimPrefix(code, ".")
			break
		}
	}
	if index := strings.Index(code, "."); index >= 0 {
		code = code[:index]
	}
	return strings.TrimSpace(code)
}

func symbolMarket(symbol string) string {
	value := strings.TrimSpace(symbol)
	lower := strings.ToLower(value)
	for _, prefix := range []string{"sh", "sz", "bj"} {
		if strings.HasPrefix(lower, prefix) {
			return prefix
		}
	}
	if index := strings.LastIndex(value, "."); index >= 0 && index+1 < len(value) {
		market := strings.ToLower(value[index+1:])
		if market == "sh" || market == "sz" || market == "bj" {
			return market
		}
	}
	return ""
}

func firstDynamicItem(items []boardDynamicItem) boardDynamicItem {
	if len(items) == 0 || items[0] == nil {
		return boardDynamicItem{}
	}
	return items[0]
}

func defaultPositive(value int, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func intString(value int) string {
	return strings.TrimSpace(stringValue(value))
}

func normalizeDateTime(value string) string {
	text := strings.TrimSpace(value)
	if len(text) >= 19 {
		return text[:19]
	}
	return text
}
