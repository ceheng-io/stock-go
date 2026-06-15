package main

import (
	"net/http"
)

type Server struct {
	sdk SDK
	mux *http.ServeMux
}

func NewServer(sdk SDK) *Server {
	server := &Server{
		sdk: sdk,
		mux: http.NewServeMux(),
	}
	server.routes()
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/health", s.handleHealth)
	s.mux.HandleFunc("GET /api/search", s.handleSearch)
	s.mux.HandleFunc("GET /api/quotes/full", s.handleFullQuotes)
	s.mux.HandleFunc("GET /api/quotes/batch", s.handleBatchQuotes)
	s.mux.HandleFunc("GET /api/quotes/a-share", s.handleAllAShareQuotes)
	s.mux.HandleFunc("GET /api/fund-flow/quotes", s.handleQuoteFundFlow)
	s.mux.HandleFunc("GET /api/panel-large-order", s.handlePanelLargeOrder)
	s.mux.HandleFunc("GET /api/timeline/today", s.handleTodayTimeline)
	s.mux.HandleFunc("GET /api/kline/history", s.handleHistoryKline)
	s.mux.HandleFunc("GET /api/kline/minute", s.handleMinuteKline)
	s.mux.HandleFunc("GET /api/kline/indicators", s.handleKlineWithIndicators)
	s.mux.HandleFunc("GET /api/boards/industry", s.handleIndustryList)
	s.mux.HandleFunc("GET /api/boards/concept", s.handleConceptList)
	s.mux.HandleFunc("GET /api/boards/{type}/{code}/spot", s.handleBoardSpot)
	s.mux.HandleFunc("GET /api/boards/{type}/{code}/constituents", s.handleBoardConstituents)
	s.mux.HandleFunc("GET /api/boards/{type}/{code}/kline", s.handleBoardKline)
	s.mux.HandleFunc("GET /api/boards/{type}/{code}/minute", s.handleBoardMinute)
	s.mux.HandleFunc("GET /api/fund-flow/individual", s.handleIndividualFundFlow)
	s.mux.HandleFunc("GET /api/fund-flow/market", s.handleMarketFundFlow)
	s.mux.HandleFunc("GET /api/fund-flow/rank", s.handleFundFlowRank)
	s.mux.HandleFunc("GET /api/fund-flow/sector-rank", s.handleSectorFundFlowRank)
	s.mux.HandleFunc("GET /api/fund-flow/sector-history", s.handleSectorFundFlowHistory)
	s.mux.HandleFunc("GET /api/northbound/minute", s.handleNorthboundMinute)
	s.mux.HandleFunc("GET /api/northbound/summary", s.handleNorthboundSummary)
	s.mux.HandleFunc("GET /api/northbound/holding-rank", s.handleNorthboundHoldingRank)
	s.mux.HandleFunc("GET /api/northbound/history", s.handleNorthboundHistory)
	s.mux.HandleFunc("GET /api/northbound/individual", s.handleNorthboundIndividual)
	s.mux.HandleFunc("GET /api/market-event/zt-pool", s.handleZTPool)
	s.mux.HandleFunc("GET /api/market-event/stock-changes", s.handleStockChanges)
	s.mux.HandleFunc("GET /api/market-event/board-changes", s.handleBoardChanges)
	s.mux.HandleFunc("GET /api/dragon-tiger/detail", s.handleDragonTigerDetail)
	s.mux.HandleFunc("GET /api/block-trade/detail", s.handleBlockTradeDetail)
	s.mux.HandleFunc("GET /api/margin/account", s.handleMarginAccount)
	s.mux.HandleFunc("GET /api/dividends", s.handleDividendDetail)
	s.mux.HandleFunc("GET /api/trading-calendar", s.handleTradingCalendar)
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
