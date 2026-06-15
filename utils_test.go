package stock

import (
	"context"
	"testing"
)

func TestRootReExportsUtils(t *testing.T) {
	chunks, err := ChunkArray([]string{"a", "b", "c"}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 2 || chunks[0][0] != "a" || chunks[1][0] != "c" {
		t.Fatalf("chunks = %#v", chunks)
	}

	tasks := []func(context.Context) (int, error){
		func(context.Context) (int, error) { return 1, nil },
		func(context.Context) (int, error) { return 2, nil },
	}
	values, err := AsyncPool(context.Background(), tasks, 2, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Fatalf("values = %#v", values)
	}

	if err := AssertPositiveInteger(1, "n"); err != nil {
		t.Fatalf("AssertPositiveInteger = %v", err)
	}
	if err := AssertKlinePeriod("daily"); err != nil {
		t.Fatalf("AssertKlinePeriod = %v", err)
	}
	if err := AssertMinutePeriod("5"); err != nil {
		t.Fatalf("AssertMinutePeriod = %v", err)
	}
	if err := AssertAdjustType("qfq"); err != nil {
		t.Fatalf("AssertAdjustType = %v", err)
	}
	if err := AssertAdjustType(string(AdjustNone)); err != nil {
		t.Fatalf("AssertAdjustType(AdjustNone) = %v", err)
	}
	period, err := PeriodCode("weekly")
	if err != nil {
		t.Fatalf("PeriodCode = %v", err)
	}
	if period != "102" {
		t.Fatalf("PeriodCode = %q, want 102", period)
	}
	tsPeriod, err := GetPeriodCode("monthly")
	if err != nil {
		t.Fatalf("GetPeriodCode = %v", err)
	}
	if tsPeriod != "103" {
		t.Fatalf("GetPeriodCode = %q, want 103", tsPeriod)
	}
	adjust, err := AdjustCode("hfq")
	if err != nil {
		t.Fatalf("AdjustCode = %v", err)
	}
	if adjust != "2" {
		t.Fatalf("AdjustCode = %q, want 2", adjust)
	}
	tsAdjust, err := GetAdjustCode("qfq")
	if err != nil {
		t.Fatalf("GetAdjustCode = %v", err)
	}
	if tsAdjust != "1" {
		t.Fatalf("GetAdjustCode = %q, want 1", tsAdjust)
	}
	noneAdjust, err := GetAdjustCode(string(AdjustNone))
	if err != nil {
		t.Fatalf("GetAdjustCode(AdjustNone) = %v", err)
	}
	if noneAdjust != "0" {
		t.Fatalf("GetAdjustCode(AdjustNone) = %q, want 0", noneAdjust)
	}
}

func TestRootUtilsInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "assert positive integer",
			call: func() error { return AssertPositiveInteger(0, "n") },
		},
		{
			name: "assert kline period",
			call: func() error { return AssertKlinePeriod("yearly") },
		},
		{
			name: "assert minute period",
			call: func() error { return AssertMinutePeriod("2") },
		},
		{
			name: "assert adjust type",
			call: func() error { return AssertAdjustType("bad") },
		},
		{
			name: "chunk size",
			call: func() error {
				_, err := ChunkArray([]int{1, 2, 3}, 0)
				return err
			},
		},
		{
			name: "async pool concurrency",
			call: func() error {
				_, err := AsyncPool(context.Background(), []func(context.Context) (int, error){
					func(context.Context) (int, error) { return 1, nil },
				}, 0, true)
				return err
			},
		},
		{
			name: "period code",
			call: func() error {
				_, err := GetPeriodCode("yearly")
				return err
			},
		},
		{
			name: "adjust code",
			call: func() error {
				_, err := GetAdjustCode("bad")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}
