package stock

import (
	"testing"

	"github.com/ceheng.io/stock-go/types"
)

func TestGenerateSearchExternalLinksForTencentMarketIDs(t *testing.T) {
	tests := []struct {
		name          string
		result        types.SearchResult
		wantEastmoney string
		wantXueqiu    string
	}{
		{
			name:          "A-share Shanghai",
			result:        types.SearchResult{Code: "sh600519", Market: "1"},
			wantEastmoney: "https://quote.eastmoney.com/sh600519.html",
			wantXueqiu:    "https://xueqiu.com/S/SH600519",
		},
		{
			name:          "A-share index",
			result:        types.SearchResult{Code: "sh000001", Market: "1"},
			wantEastmoney: "https://quote.eastmoney.com/zs000001.html",
			wantXueqiu:    "https://xueqiu.com/S/SH000001",
		},
		{
			name:          "HK stock pads code",
			result:        types.SearchResult{Code: "hk700", Market: "116"},
			wantEastmoney: "https://quote.eastmoney.com/hk/00700.html",
			wantXueqiu:    "https://xueqiu.com/S/00700",
		},
		{
			name:          "US stock strips Tencent prefix",
			result:        types.SearchResult{Code: "105.BABA", Market: "105"},
			wantEastmoney: "https://quote.eastmoney.com/us/BABA.html",
			wantXueqiu:    "https://xueqiu.com/S/BABA",
		},
		{
			name:          "global index maps aliases",
			result:        types.SearchResult{Code: "100.IXIC", Market: "100"},
			wantEastmoney: "https://quote.eastmoney.com/gb/zsNDX.html",
			wantXueqiu:    "https://xueqiu.com/S/.IXIC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links := GenerateSearchExternalLinks(tt.result)
			if len(links) != 2 {
				t.Fatalf("len(links) = %d, want 2", len(links))
			}
			if links[0].Name != "东方财富" || links[0].URL != tt.wantEastmoney {
				t.Fatalf("Eastmoney link = %+v, want %q", links[0], tt.wantEastmoney)
			}
			if links[1].Name != "雪球" || links[1].URL != tt.wantXueqiu {
				t.Fatalf("Xueqiu link = %+v, want %q", links[1], tt.wantXueqiu)
			}
		})
	}
}

func TestGenerateSearchExternalLinksFallsBackToSearch(t *testing.T) {
	links := GenerateSearchExternalLinks(types.SearchResult{Code: " mystery code ", Market: "unknown"})

	if links[0].URL != "https://so.eastmoney.com/web/s?keyword=mystery+code" {
		t.Fatalf("Eastmoney fallback = %q", links[0].URL)
	}
	if links[1].URL != "https://xueqiu.com/k?q=mystery+code" {
		t.Fatalf("Xueqiu fallback = %q", links[1].URL)
	}
}
