package stock

import "testing"

func TestRootReExportsDomainOptionAliases(t *testing.T) {
	var _ MarketType = MarketA
	var _ SupportedMarket = MarketType(MarketHK)
	var _ GetAShareCodeListOptions = GetAShareCodeListOptions{Simple: true, Market: AShareMarketSH}
	var _ GetAllAShareQuotesOptions = GetAllAShareQuotesOptions{
		BatchSize:   100,
		Concurrency: 2,
		Market:      AShareMarketSZ,
	}
	var _ GetAllHKQuotesOptions = GetAllHKQuotesOptions{BatchSize: 100, Concurrency: 2}
	var _ GetUSCodeListOptions = GetUSCodeListOptions{Simple: true, Market: USMarketNASDAQ}
	var _ GetAllUSQuotesOptions = GetAllUSQuotesOptions{
		BatchSize:   100,
		Concurrency: 2,
		Market:      USMarketNYSE,
	}
	var _ HistoryKlineOptions = HKKlineOptions{}
	var _ HistoryKlineOptions = USKlineOptions{}
	var _ MinuteKlineOptions = HKMinuteKlineOptions{}
	var _ MinuteKlineOptions = USMinuteKlineOptions{}
	var _ HistoryKlineOptions = IndustryBoardKlineOptions{}
	var _ HistoryKlineOptions = ConceptBoardKlineOptions{}
	var _ MinuteKlineOptions = IndustryBoardMinuteKlineOptions{}
	var _ MinuteKlineOptions = ConceptBoardMinuteKlineOptions{}
}
