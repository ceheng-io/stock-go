package main

import (
	"net/http"
	"strings"

	stock "github.com/ceheng.io/stock-go"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	if keyword == "" {
		writeError(w, http.StatusBadRequest, "keyword is required")
		return
	}
	if s.sdk == nil {
		writeError(w, http.StatusServiceUnavailable, "sdk is not configured")
		return
	}
	results, err := s.sdk.Search(r.Context(), keyword)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (s *Server) ensureSDK(w http.ResponseWriter) bool {
	if s.sdk == nil {
		writeError(w, http.StatusServiceUnavailable, "sdk is not configured")
		return false
	}
	return true
}

func (s *Server) handleFullQuotes(w http.ResponseWriter, r *http.Request) {
	codes, ok := queryCodes(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "codes is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetFullQuotes(r.Context(), codes)
	writeResult(w, data, err)
}

func (s *Server) handleBatchQuotes(w http.ResponseWriter, r *http.Request) {
	codes, ok := queryCodes(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "codes is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetAllQuotesByCodes(r.Context(), codes, asBatchOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleAllAShareQuotes(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetAllAShareQuotes(r.Context(), asBatchOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleQuoteFundFlow(w http.ResponseWriter, r *http.Request) {
	codes, ok := queryCodes(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "codes is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetFundFlow(r.Context(), codes)
	writeResult(w, data, err)
}

func (s *Server) handlePanelLargeOrder(w http.ResponseWriter, r *http.Request) {
	codes, ok := queryCodes(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "codes is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetPanelLargeOrder(r.Context(), codes)
	writeResult(w, data, err)
}

func (s *Server) handleTodayTimeline(w http.ResponseWriter, r *http.Request) {
	code, ok := requiredQuery(r, "code")
	if !ok {
		writeError(w, http.StatusBadRequest, "code is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetTodayTimeline(r.Context(), code)
	writeResult(w, data, err)
}

func (s *Server) handleHistoryKline(w http.ResponseWriter, r *http.Request) {
	symbol, ok := requiredQuery(r, "symbol")
	if !ok {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetHistoryKline(r.Context(), symbol, asHistoryOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleMinuteKline(w http.ResponseWriter, r *http.Request) {
	symbol, ok := requiredQuery(r, "symbol")
	if !ok {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetMinuteKline(r.Context(), symbol, asMinuteOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleKlineWithIndicators(w http.ResponseWriter, r *http.Request) {
	symbol, ok := requiredQuery(r, "symbol")
	if !ok {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	options := stock.KlineWithIndicatorsOptions{
		Period:     stock.KlinePeriod(queryString(r, "period", "")),
		Adjust:     stock.AdjustType(queryString(r, "adjust", "")),
		StartDate:  queryString(r, "startDate", ""),
		EndDate:    queryString(r, "endDate", ""),
		Indicators: defaultIndicatorOptions(),
	}
	data, err := s.sdk.GetKlineWithIndicators(r.Context(), symbol, options)
	writeResult(w, data, err)
}

func defaultIndicatorOptions() stock.IndicatorOptions {
	return stock.IndicatorOptions{
		MA:   &stock.MAOptions{Periods: []int{5, 10, 20, 60}},
		MACD: &stock.MACDOptions{Short: 12, Long: 26, Signal: 9},
		BOLL: &stock.BOLLOptions{Period: 20, StdDev: 2},
		KDJ:  &stock.KDJOptions{Period: 9, KPeriod: 3, DPeriod: 3},
		RSI:  &stock.RSIOptions{Periods: []int{6, 12, 24}},
		OBV:  &stock.OBVOptions{MAPeriod: 10},
		ROC:  &stock.ROCOptions{Period: 12, SignalPeriod: 6},
		DMI:  &stock.DMIOptions{Period: 14, ADXPeriod: 14},
		SAR:  &stock.SAROptions{AFStart: 0.02, AFIncrement: 0.02, AFMax: 0.2},
		KC:   &stock.KCOptions{EMAPeriod: 20, ATRPeriod: 10, Multiplier: 2},
	}
}

func (s *Server) handleIndustryList(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetIndustryList(r.Context())
	writeResult(w, data, err)
}

func (s *Server) handleConceptList(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetConceptList(r.Context())
	writeResult(w, data, err)
}

func (s *Server) handleBoardSpot(w http.ResponseWriter, r *http.Request) {
	boardType, code := r.PathValue("type"), r.PathValue("code")
	if !s.ensureSDK(w) {
		return
	}
	var data []stock.BoardSpot
	var err error
	if boardType == "industry" {
		data, err = s.sdk.GetIndustrySpot(r.Context(), code)
	} else if boardType == "concept" {
		data, err = s.sdk.GetConceptSpot(r.Context(), code)
	} else {
		writeError(w, http.StatusBadRequest, "unsupported board type")
		return
	}
	writeResult(w, data, err)
}

func (s *Server) handleBoardConstituents(w http.ResponseWriter, r *http.Request) {
	boardType, code := r.PathValue("type"), r.PathValue("code")
	if !s.ensureSDK(w) {
		return
	}
	var data []stock.BoardConstituent
	var err error
	if boardType == "industry" {
		data, err = s.sdk.GetIndustryConstituents(r.Context(), code)
	} else if boardType == "concept" {
		data, err = s.sdk.GetConceptConstituents(r.Context(), code)
	} else {
		writeError(w, http.StatusBadRequest, "unsupported board type")
		return
	}
	writeResult(w, data, err)
}

func (s *Server) handleBoardKline(w http.ResponseWriter, r *http.Request) {
	boardType, code := r.PathValue("type"), r.PathValue("code")
	if !s.ensureSDK(w) {
		return
	}
	options := asHistoryOptions(r)
	var data []stock.BoardKline
	var err error
	if boardType == "industry" {
		data, err = s.sdk.GetIndustryKline(r.Context(), code, stock.IndustryBoardKlineOptions(options))
	} else if boardType == "concept" {
		data, err = s.sdk.GetConceptKline(r.Context(), code, stock.ConceptBoardKlineOptions(options))
	} else {
		writeError(w, http.StatusBadRequest, "unsupported board type")
		return
	}
	writeResult(w, data, err)
}

func (s *Server) handleBoardMinute(w http.ResponseWriter, r *http.Request) {
	boardType, code := r.PathValue("type"), r.PathValue("code")
	if !s.ensureSDK(w) {
		return
	}
	options := asMinuteOptions(r)
	var data stock.BoardMinuteKlineResult
	var err error
	if boardType == "industry" {
		data, err = s.sdk.GetIndustryMinuteKline(r.Context(), code, stock.IndustryBoardMinuteKlineOptions(options))
	} else if boardType == "concept" {
		data, err = s.sdk.GetConceptMinuteKline(r.Context(), code, stock.ConceptBoardMinuteKlineOptions(options))
	} else {
		writeError(w, http.StatusBadRequest, "unsupported board type")
		return
	}
	writeResult(w, data, err)
}

func (s *Server) handleIndividualFundFlow(w http.ResponseWriter, r *http.Request) {
	symbol, ok := requiredQuery(r, "symbol")
	if !ok {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetIndividualFundFlow(r.Context(), symbol, asFundFlowOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleMarketFundFlow(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetMarketFundFlow(r.Context())
	writeResult(w, data, err)
}

func (s *Server) handleFundFlowRank(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetFundFlowRank(r.Context(), asFundFlowRankOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleSectorFundFlowRank(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetSectorFundFlowRank(r.Context(), asFundFlowRankOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleSectorFundFlowHistory(w http.ResponseWriter, r *http.Request) {
	symbol, ok := requiredQuery(r, "symbol")
	if !ok {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetSectorFundFlowHistory(r.Context(), symbol, asFundFlowOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleNorthboundMinute(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetNorthboundMinute(r.Context(), stock.NorthboundDirection(queryString(r, "direction", "")))
	writeResult(w, data, err)
}

func (s *Server) handleNorthboundSummary(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetNorthboundFlowSummary(r.Context())
	writeResult(w, data, err)
}

func (s *Server) handleNorthboundHoldingRank(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	options := stock.NorthboundHoldingRankOptions{
		Market: stock.NorthboundMarket(queryString(r, "market", "")),
		Period: stock.NorthboundRankPeriod(queryString(r, "period", "")),
		Date:   queryString(r, "date", ""),
	}
	data, err := s.sdk.GetNorthboundHoldingRank(r.Context(), options)
	writeResult(w, data, err)
}

func (s *Server) handleNorthboundHistory(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetNorthboundHistory(
		r.Context(),
		stock.NorthboundDirection(queryString(r, "direction", "")),
		asNorthboundHistoryOptions(r),
	)
	writeResult(w, data, err)
}

func (s *Server) handleNorthboundIndividual(w http.ResponseWriter, r *http.Request) {
	symbol, ok := requiredQuery(r, "symbol")
	if !ok {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetNorthboundIndividual(r.Context(), symbol, asNorthboundHistoryOptions(r))
	writeResult(w, data, err)
}

func (s *Server) handleZTPool(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	args := []any{stock.ZTPoolType(queryString(r, "type", ""))}
	if date := queryString(r, "date", ""); date != "" {
		args = append(args, date)
	}
	data, err := s.sdk.GetZTPool(r.Context(), args...)
	writeResult(w, data, err)
}

func (s *Server) handleStockChanges(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetStockChanges(r.Context(), stock.StockChangeType(queryString(r, "type", "")))
	writeResult(w, data, err)
}

func (s *Server) handleBoardChanges(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetBoardChanges(r.Context())
	writeResult(w, data, err)
}

func (s *Server) handleDragonTigerDetail(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetDragonTigerDetail(r.Context(), stock.DragonTigerDateOptions{
		StartDate: queryString(r, "startDate", ""),
		EndDate:   queryString(r, "endDate", ""),
	})
	writeResult(w, data, err)
}

func (s *Server) handleBlockTradeDetail(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetBlockTradeDetail(r.Context(), stock.BlockTradeDateOptions{
		StartDate: queryString(r, "startDate", ""),
		EndDate:   queryString(r, "endDate", ""),
	})
	writeResult(w, data, err)
}

func (s *Server) handleMarginAccount(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetMarginAccountInfo(r.Context())
	writeResult(w, data, err)
}

func (s *Server) handleDividendDetail(w http.ResponseWriter, r *http.Request) {
	symbol, ok := requiredQuery(r, "symbol")
	if !ok {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetDividendDetail(r.Context(), symbol)
	writeResult(w, data, err)
}

func (s *Server) handleTradingCalendar(w http.ResponseWriter, r *http.Request) {
	if !s.ensureSDK(w) {
		return
	}
	data, err := s.sdk.GetTradingCalendar(r.Context())
	writeResult(w, data, err)
}
