package tencent

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/ceheng.io/stock-go/internal/core"
)

type batchQuoteClient struct {
	mu     sync.Mutex
	params []string
}

func (b *batchQuoteClient) GetTencentQuote(_ context.Context, params string) ([]core.TencentQuoteItem, error) {
	b.mu.Lock()
	b.params = append(b.params, params)
	b.mu.Unlock()

	items := []core.TencentQuoteItem{}
	for _, code := range strings.Split(params, ",") {
		fields := make([]string, 80)
		key := code
		if strings.HasPrefix(code, "hk") {
			fields = make([]string, 50)
			fields[0] = "100"
			fields[1] = "腾讯控股"
			fields[2] = strings.TrimPrefix(code, "hk")
		} else if strings.HasPrefix(code, "us") {
			fields = make([]string, 50)
			fields[0] = "200"
			fields[1] = "APPLE"
			fields[2] = strings.TrimPrefix(code, "us")
		} else {
			fields[0] = "1"
			fields[1] = "贵州茅台"
			fields[2] = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(code, "sh"), "sz"), "bj")
		}
		fields[3] = "1.23"
		items = append(items, core.TencentQuoteItem{Key: key, Fields: fields})
	}
	return items, nil
}

func TestGetAllQuotesByCodesChunksAndReportsProgress(t *testing.T) {
	client := &batchQuoteClient{}
	progress := []int{}

	quotes, err := GetAllQuotesByCodes(context.Background(), client, []string{"sh600519", "sz000001", "bj830799"}, BatchOptions{
		BatchSize:   2,
		Concurrency: 2,
		OnProgress: func(completed, total int) {
			progress = append(progress, completed*10+total)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 3 {
		t.Fatalf("len(quotes) = %d, want 3", len(quotes))
	}
	if len(client.params) != 2 {
		t.Fatalf("params = %#v, want 2 chunks", client.params)
	}
	if len(progress) != 2 || progress[len(progress)-1] != 22 {
		t.Fatalf("progress = %#v, want final 22", progress)
	}
}

func TestGetAllHKAndUSQuotesByCodes(t *testing.T) {
	client := &batchQuoteClient{}

	hk, err := GetAllHKQuotesByCodes(context.Background(), client, []string{"00700", "09988"}, BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 2 || hk[0].Code != "00700" {
		t.Fatalf("hk = %+v", hk)
	}

	us, err := GetAllUSQuotesByCodes(context.Background(), client, []string{"AAPL", "MSFT"}, BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 2 || us[0].Code != "AAPL" {
		t.Fatalf("us = %+v", us)
	}
}

func TestBatchOptionsValidation(t *testing.T) {
	client := &batchQuoteClient{}
	if _, err := GetAllQuotesByCodes(context.Background(), client, []string{"sh600519"}, BatchOptions{BatchSize: -1}); err == nil {
		t.Fatal("BatchSize -1 expected error")
	}
	if _, err := GetAllQuotesByCodes(context.Background(), client, []string{"sh600519"}, BatchOptions{Concurrency: -1}); err == nil {
		t.Fatal("Concurrency -1 expected error")
	}
}
