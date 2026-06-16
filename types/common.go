package types

// Market identifies a supported market.
type Market string

const (
	MarketCN     Market = "CN"
	MarketHK     Market = "HK"
	MarketUS     Market = "US"
	MarketGlobal Market = "GLOBAL"
)

// ExternalLink 是外部财经站点链接。
type ExternalLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
