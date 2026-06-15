package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

// BoardClient is the request client interface required by BoardService.
type BoardClient interface {
	eastmoney.BoardClient
}

// BoardURLs contains Eastmoney board endpoints.
type BoardURLs struct {
	IndustryList         string
	IndustrySpot         string
	IndustryConstituents string
	IndustryKline        string
	IndustryTrends       string
	ConceptList          string
	ConceptSpot          string
	ConceptConstituents  string
	ConceptKline         string
	ConceptTrends        string
}

// BoardService orchestrates Eastmoney board providers.
type BoardService struct {
	client BoardClient
	urls   BoardURLs
}

// NewBoardService creates a BoardService.
func NewBoardService(client BoardClient, urls BoardURLs) *BoardService {
	return &BoardService{client: client, urls: urls}
}

// IndustryList returns industry board list rows.
func (s *BoardService) IndustryList(ctx context.Context) ([]types.Board, error) {
	return eastmoney.GetIndustryList(ctx, s.client, s.urls.IndustryList)
}

// ConceptList returns concept board list rows.
func (s *BoardService) ConceptList(ctx context.Context) ([]types.Board, error) {
	return eastmoney.GetConceptList(ctx, s.client, s.urls.ConceptList)
}

// IndustrySpot returns industry board spot metrics.
func (s *BoardService) IndustrySpot(ctx context.Context, boardCode string) ([]types.BoardSpot, error) {
	boardCode, err := s.resolveIndustryBoardCode(ctx, boardCode)
	if err != nil {
		return nil, err
	}
	return eastmoney.GetBoardSpot(ctx, s.client, boardCode, s.urls.IndustrySpot)
}

// ConceptSpot returns concept board spot metrics.
func (s *BoardService) ConceptSpot(ctx context.Context, boardCode string) ([]types.BoardSpot, error) {
	boardCode, err := s.resolveConceptBoardCode(ctx, boardCode)
	if err != nil {
		return nil, err
	}
	return eastmoney.GetBoardSpot(ctx, s.client, boardCode, s.urls.ConceptSpot)
}

// IndustryConstituents returns industry board constituent stocks.
func (s *BoardService) IndustryConstituents(ctx context.Context, boardCode string) ([]types.BoardConstituent, error) {
	boardCode, err := s.resolveIndustryBoardCode(ctx, boardCode)
	if err != nil {
		return nil, err
	}
	return eastmoney.GetBoardConstituents(ctx, s.client, boardCode, s.urls.IndustryConstituents)
}

// ConceptConstituents returns concept board constituent stocks.
func (s *BoardService) ConceptConstituents(ctx context.Context, boardCode string) ([]types.BoardConstituent, error) {
	boardCode, err := s.resolveConceptBoardCode(ctx, boardCode)
	if err != nil {
		return nil, err
	}
	return eastmoney.GetBoardConstituents(ctx, s.client, boardCode, s.urls.ConceptConstituents)
}

// IndustryKline returns industry board historical K-line rows.
func (s *BoardService) IndustryKline(ctx context.Context, boardCode string, options eastmoney.HistoryKlineOptions) ([]types.BoardKline, error) {
	boardCode, err := s.resolveIndustryBoardCode(ctx, boardCode)
	if err != nil {
		return nil, err
	}
	return eastmoney.GetBoardKline(ctx, s.client, boardCode, s.urls.IndustryKline, options)
}

// ConceptKline returns concept board historical K-line rows.
func (s *BoardService) ConceptKline(ctx context.Context, boardCode string, options eastmoney.HistoryKlineOptions) ([]types.BoardKline, error) {
	boardCode, err := s.resolveConceptBoardCode(ctx, boardCode)
	if err != nil {
		return nil, err
	}
	return eastmoney.GetBoardKline(ctx, s.client, boardCode, s.urls.ConceptKline, options)
}

// IndustryMinute returns industry board minute timeline or K-line rows.
func (s *BoardService) IndustryMinute(ctx context.Context, boardCode string, options eastmoney.MinuteKlineOptions) (types.BoardMinuteKlineResult, error) {
	boardCode, err := s.resolveIndustryBoardCode(ctx, boardCode)
	if err != nil {
		return types.BoardMinuteKlineResult{}, err
	}
	return eastmoney.GetBoardMinuteKline(ctx, s.client, boardCode, s.urls.IndustryKline, s.urls.IndustryTrends, options)
}

// ConceptMinute returns concept board minute timeline or K-line rows.
func (s *BoardService) ConceptMinute(ctx context.Context, boardCode string, options eastmoney.MinuteKlineOptions) (types.BoardMinuteKlineResult, error) {
	boardCode, err := s.resolveConceptBoardCode(ctx, boardCode)
	if err != nil {
		return types.BoardMinuteKlineResult{}, err
	}
	return eastmoney.GetBoardMinuteKline(ctx, s.client, boardCode, s.urls.ConceptKline, s.urls.ConceptTrends, options)
}

func (s *BoardService) resolveIndustryBoardCode(ctx context.Context, symbol string) (string, error) {
	return s.resolveBoardCode(ctx, symbol, s.IndustryList, "Industry board not found")
}

func (s *BoardService) resolveConceptBoardCode(ctx context.Context, symbol string) (string, error) {
	return s.resolveBoardCode(ctx, symbol, s.ConceptList, "Concept board not found")
}

func (s *BoardService) resolveBoardCode(ctx context.Context, symbol string, list func(context.Context) ([]types.Board, error), messagePrefix string) (string, error) {
	trimmed := strings.TrimSpace(symbol)
	if strings.HasPrefix(strings.ToUpper(trimmed), "BK") {
		return trimmed, nil
	}
	boards, err := list(ctx)
	if err != nil {
		return "", err
	}
	for _, board := range boards {
		if board.Name == symbol || board.Name == trimmed {
			return board.Code, nil
		}
	}
	return "", core.NewCodedError("NOT_FOUND", fmt.Sprintf("%s: %s", messagePrefix, symbol), nil)
}
