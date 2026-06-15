package utils

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestChunkArraySplitsValues(t *testing.T) {
	chunks, err := ChunkArray([]int{1, 2, 3, 4, 5}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 3 {
		t.Fatalf("len(chunks) = %d, want 3", len(chunks))
	}
	assertInts(t, chunks[0], []int{1, 2})
	assertInts(t, chunks[1], []int{3, 4})
	assertInts(t, chunks[2], []int{5})
}

func TestChunkArrayRejectsInvalidSize(t *testing.T) {
	if _, err := ChunkArray([]string{"a"}, 0); err == nil {
		t.Fatal("expected error for zero chunk size")
	}
	if _, err := ChunkArray([]string{"a"}, -1); err == nil {
		t.Fatal("expected error for negative chunk size")
	}
}

func TestAsyncPoolPreservesOrderAndLimitsConcurrency(t *testing.T) {
	var running int32
	var maxRunning int32
	started := make(chan struct{}, 4)
	release := make(chan struct{})
	tasks := make([]func(context.Context) (int, error), 4)
	for i := range tasks {
		value := i
		tasks[i] = func(context.Context) (int, error) {
			current := atomic.AddInt32(&running, 1)
			updateMax(&maxRunning, current)
			started <- struct{}{}
			<-release
			atomic.AddInt32(&running, -1)
			return value, nil
		}
	}

	done := make(chan struct{})
	var got []int
	var err error
	go func() {
		got, err = AsyncPool(context.Background(), tasks, 2, true)
		close(done)
	}()

	<-started
	<-started
	select {
	case <-started:
		t.Fatal("third task started before concurrency slot was released")
	default:
	}
	close(release)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("AsyncPool did not finish")
	}
	if err != nil {
		t.Fatal(err)
	}
	assertInts(t, got, []int{0, 1, 2, 3})
	if maxRunning != 2 {
		t.Fatalf("max concurrent tasks = %d, want 2", maxRunning)
	}
}

func TestAsyncPoolReturnsFirstError(t *testing.T) {
	wantErr := errors.New("boom")
	tasks := []func(context.Context) (int, error){
		func(context.Context) (int, error) { return 1, nil },
		func(context.Context) (int, error) { return 0, wantErr },
	}

	_, err := AsyncPool(context.Background(), tasks, 2, true)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestAsyncPoolStopsOnCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tasks := []func(context.Context) (int, error){
		func(context.Context) (int, error) { return 1, nil },
	}
	_, err := AsyncPool(ctx, tasks, 1, true)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestValidationHelpers(t *testing.T) {
	if err := AssertPositiveInteger(3, "size"); err != nil {
		t.Fatalf("AssertPositiveInteger valid = %v", err)
	}
	for _, value := range []int{0, -1} {
		if err := AssertPositiveInteger(value, "size"); err == nil {
			t.Fatalf("AssertPositiveInteger(%d) expected error", value)
		}
	}

	for _, period := range []string{"daily", "weekly", "monthly"} {
		if err := AssertKlinePeriod(period); err != nil {
			t.Fatalf("AssertKlinePeriod(%q) = %v", period, err)
		}
	}
	if err := AssertKlinePeriod("yearly"); err == nil {
		t.Fatal("AssertKlinePeriod invalid expected error")
	}

	for _, period := range []string{"1", "5", "15", "30", "60"} {
		if err := AssertMinutePeriod(period); err != nil {
			t.Fatalf("AssertMinutePeriod(%q) = %v", period, err)
		}
	}
	if err := AssertMinutePeriod("120"); err == nil {
		t.Fatal("AssertMinutePeriod invalid expected error")
	}

	for _, adjust := range []string{"", "none", "qfq", "hfq"} {
		if err := AssertAdjustType(adjust); err != nil {
			t.Fatalf("AssertAdjustType(%q) = %v", adjust, err)
		}
	}
	if err := AssertAdjustType("bad"); err == nil {
		t.Fatal("AssertAdjustType invalid expected error")
	}
}

func TestPeriodAndAdjustCodes(t *testing.T) {
	tests := map[string]string{
		"daily":   "101",
		"weekly":  "102",
		"monthly": "103",
	}
	for period, want := range tests {
		got, err := PeriodCode(period)
		if err != nil {
			t.Fatalf("PeriodCode(%q) = %v", period, err)
		}
		if got != want {
			t.Fatalf("PeriodCode(%q) = %q, want %q", period, got, want)
		}
	}
	if _, err := PeriodCode("yearly"); err == nil {
		t.Fatal("PeriodCode invalid expected error")
	}

	adjustTests := map[string]string{
		"":     "0",
		"none": "0",
		"qfq":  "1",
		"hfq":  "2",
	}
	for adjust, want := range adjustTests {
		got, err := AdjustCode(adjust)
		if err != nil {
			t.Fatalf("AdjustCode(%q) = %v", adjust, err)
		}
		if got != want {
			t.Fatalf("AdjustCode(%q) = %q, want %q", adjust, got, want)
		}
	}
	if _, err := AdjustCode("bad"); err == nil {
		t.Fatal("AdjustCode invalid expected error")
	}
}

func updateMax(target *int32, value int32) {
	for {
		current := atomic.LoadInt32(target)
		if value <= current {
			return
		}
		if atomic.CompareAndSwapInt32(target, current, value) {
			return
		}
	}
}

func assertInts(t *testing.T, got []int, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d; got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d; got=%v", i, got[i], want[i], got)
		}
	}
}
