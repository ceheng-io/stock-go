package tencent

import (
	"context"
	"regexp"
	"strings"
)

// AShareMarket is an A-share exchange or board filter.
type AShareMarket string

const (
	AShareMarketSH AShareMarket = "sh"
	AShareMarketSZ AShareMarket = "sz"
	AShareMarketBJ AShareMarket = "bj"
	AShareMarketKC AShareMarket = "kc"
	AShareMarketCY AShareMarket = "cy"
)

// CodeListOptions configures CN code-list fetching.
type CodeListOptions struct {
	Simple bool
	Market AShareMarket
}

// USMarket is a US exchange filter.
type USMarket string

const (
	USMarketNASDAQ USMarket = "NASDAQ"
	USMarketNYSE   USMarket = "NYSE"
	USMarketAMEX   USMarket = "AMEX"
)

// USCodeListOptions configures US code-list fetching.
type USCodeListOptions struct {
	Simple bool
	Market USMarket
}

// CodeListClient is the minimal client interface required by code-list providers.
type CodeListClient interface {
	GetJSON(context.Context, string, any) error
	GetText(context.Context, string) (string, error)
	AShareListURL() string
	USListURL() string
	HKListURL() string
	FundListURL() string
}

type codeListResponse struct {
	Success bool     `json:"success"`
	List    []string `json:"list"`
}

func fetchCodeList(ctx context.Context, client CodeListClient, requestURL string) ([]string, error) {
	var payload codeListResponse
	if err := client.GetJSON(ctx, requestURL, &payload); err != nil {
		return nil, err
	}
	return append([]string(nil), payload.List...), nil
}

// GetAShareCodeList returns A-share codes.
func GetAShareCodeList(ctx context.Context, client CodeListClient, options CodeListOptions) ([]string, error) {
	codes, err := fetchCodeList(ctx, client, client.AShareListURL())
	if err != nil {
		return nil, err
	}
	if options.Market != "" {
		filtered := codes[:0]
		for _, code := range codes {
			if matchAShareMarket(code, options.Market) {
				filtered = append(filtered, code)
			}
		}
		codes = filtered
	}
	if options.Simple {
		for i, code := range codes {
			codes[i] = stripASharePrefix(code)
		}
	}
	return append([]string(nil), codes...), nil
}

var aSharePrefixPattern = regexp.MustCompile(`^(sh|sz|bj)`)
var usMarketCodePrefixPattern = regexp.MustCompile(`^\d{3}\.`)

func stripASharePrefix(code string) string {
	return aSharePrefixPattern.ReplaceAllString(code, "")
}

func matchAShareMarket(code string, market AShareMarket) bool {
	pureCode := stripASharePrefix(code)
	switch market {
	case AShareMarketSH:
		return strings.HasPrefix(pureCode, "6")
	case AShareMarketSZ:
		return strings.HasPrefix(pureCode, "0") || strings.HasPrefix(pureCode, "3")
	case AShareMarketBJ:
		return strings.HasPrefix(pureCode, "4") || strings.HasPrefix(pureCode, "8") || strings.HasPrefix(pureCode, "92")
	case AShareMarketKC:
		return strings.HasPrefix(pureCode, "688")
	case AShareMarketCY:
		return strings.HasPrefix(pureCode, "30")
	default:
		return true
	}
}

// GetUSCodeList returns US codes.
func GetUSCodeList(ctx context.Context, client CodeListClient, options USCodeListOptions) ([]string, error) {
	codes, err := fetchCodeList(ctx, client, client.USListURL())
	if err != nil {
		return nil, err
	}
	if options.Market != "" {
		prefix := usMarketPrefix(options.Market)
		filtered := codes[:0]
		for _, code := range codes {
			if strings.HasPrefix(code, prefix) {
				filtered = append(filtered, code)
			}
		}
		codes = filtered
	}
	if options.Simple {
		for i, code := range codes {
			codes[i] = usMarketCodePrefixPattern.ReplaceAllString(code, "")
		}
	}
	return append([]string(nil), codes...), nil
}

func usMarketPrefix(market USMarket) string {
	switch market {
	case USMarketNASDAQ:
		return "105."
	case USMarketNYSE:
		return "106."
	case USMarketAMEX:
		return "107."
	default:
		return ""
	}
}

// GetHKCodeList returns HK codes.
func GetHKCodeList(ctx context.Context, client CodeListClient) ([]string, error) {
	return fetchCodeList(ctx, client, client.HKListURL())
}

// GetFundCodeList returns fund codes.
func GetFundCodeList(ctx context.Context, client CodeListClient) ([]string, error) {
	text, err := client.GetText(ctx, client.FundListURL())
	if err != nil {
		return nil, err
	}
	parts := strings.Split(text, ",")
	if len(parts) > 0 {
		parts = parts[1:]
	}
	codes := make([]string, 0, len(parts))
	for _, part := range parts {
		code := strings.TrimSpace(part)
		if code != "" {
			codes = append(codes, code)
		}
	}
	return codes, nil
}
