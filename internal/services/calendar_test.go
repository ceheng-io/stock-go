package services

import (
	"context"
	"testing"
	"time"
)

type calendarQuoteStub struct {
	days []string
}

func (c calendarQuoteStub) TradingCalendar(context.Context) ([]string, error) {
	return append([]string(nil), c.days...), nil
}

func TestCalendarServiceTradingDays(t *testing.T) {
	service := NewCalendarService(calendarQuoteStub{days: []string{
		"2024-06-12",
		"2024-06-13",
		"2024-06-17",
	}})

	ok, err := service.IsTradingDay(context.Background(), "20240613")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected 2024-06-13 to be trading day")
	}
	ok, err = service.IsTradingDay(context.Background(), "2024-06-15")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected 2024-06-15 to be non-trading day")
	}

	next, err := service.NextTradingDay(context.Background(), "2024-06-13")
	if err != nil {
		t.Fatal(err)
	}
	if next != "2024-06-17" {
		t.Fatalf("next = %q, want 2024-06-17", next)
	}
	prev, err := service.PrevTradingDay(context.Background(), "2024-06-15")
	if err != nil {
		t.Fatal(err)
	}
	if prev != "2024-06-13" {
		t.Fatalf("prev = %q, want 2024-06-13", prev)
	}
}

func TestCalendarServiceMarketStatus(t *testing.T) {
	service := NewCalendarService(calendarQuoteStub{})

	open := time.Date(2024, 6, 13, 2, 0, 0, 0, time.UTC) // Asia/Shanghai 10:00 Thu
	if got := service.MarketStatus(MarketA, open); got != MarketStatusOpen {
		t.Fatalf("A open = %s", got)
	}
	lunch := time.Date(2024, 6, 13, 4, 0, 0, 0, time.UTC) // Asia/Shanghai 12:00 Thu
	if got := service.MarketStatus(MarketA, lunch); got != MarketStatusLunchBreak {
		t.Fatalf("A lunch = %s", got)
	}
	closed := time.Date(2024, 6, 15, 2, 0, 0, 0, time.UTC) // Saturday
	if got := service.MarketStatus(MarketA, closed); got != MarketStatusClosed {
		t.Fatalf("A closed = %s", got)
	}
}
