package timeutil

import "testing"

func TestParseMarketTimeUsesMarketTimeZones(t *testing.T) {
	tests := []struct {
		name  string
		local string
		tz    MarketTz
		want  int64
	}{
		{
			name:  "CN compact quote time",
			local: "20240613093000",
			tz:    MarketTZ.CN,
			want:  1718242200000,
		},
		{
			name:  "HK separated time",
			local: "2024-06-13 09:30",
			tz:    MarketTZ.HK,
			want:  1718242200000,
		},
		{
			name:  "US summer time",
			local: "2024-06-13 09:30",
			tz:    MarketTZ.US,
			want:  1718285400000,
		},
		{
			name:  "US winter time",
			local: "2024-12-13 09:30",
			tz:    MarketTZ.US,
			want:  1734100200000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseMarketTime(tt.local, tt.tz)
			if !ok {
				t.Fatalf("ParseMarketTime(%q, %q) ok = false", tt.local, tt.tz)
			}
			if got != tt.want {
				t.Fatalf("timestamp = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBuildTimeMeta(t *testing.T) {
	meta := BuildTimeMeta("2024-06-13", MarketTZ.CN)
	if meta.Timestamp == nil || *meta.Timestamp != 1718208000000 {
		t.Fatalf("timestamp = %v, want 1718208000000", meta.Timestamp)
	}
	if meta.TZ != MarketTZ.CN {
		t.Fatalf("tz = %q, want %q", meta.TZ, MarketTZ.CN)
	}

	invalid := BuildTimeMeta("not-time", MarketTZ.CN)
	if invalid.Timestamp != nil {
		t.Fatalf("invalid timestamp = %v, want nil", invalid.Timestamp)
	}
}

func TestBuildTimeMetaFromDateAndTime(t *testing.T) {
	meta := BuildTimeMetaFromDateAndTime("20240613", "09:31:00", MarketTZ.CN)
	if meta.Timestamp == nil || *meta.Timestamp != 1718242260000 {
		t.Fatalf("timestamp = %v, want 1718242260000", meta.Timestamp)
	}

	invalid := BuildTimeMetaFromDateAndTime("2024/06/13", "09:31", MarketTZ.CN)
	if invalid.Timestamp != nil {
		t.Fatalf("invalid timestamp = %v, want nil", invalid.Timestamp)
	}
}

func TestFormatInTz(t *testing.T) {
	tests := []struct {
		name  string
		epoch *int64
		tz    MarketTz
		want  string
	}{
		{
			name:  "CN",
			epoch: int64Ptr(1718242200000),
			tz:    MarketTZ.CN,
			want:  "2024-06-13 09:30",
		},
		{
			name:  "US summer",
			epoch: int64Ptr(1718285400000),
			tz:    MarketTZ.US,
			want:  "2024-06-13 09:30",
		},
		{
			name:  "nil",
			epoch: nil,
			tz:    MarketTZ.CN,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatInTz(tt.epoch, tt.tz); got != tt.want {
				t.Fatalf("FormatInTz = %q, want %q", got, tt.want)
			}
		})
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}

func TestParseMarketTimeRejectsInvalidInput(t *testing.T) {
	if got, ok := ParseMarketTime("", MarketTZ.CN); ok || got != 0 {
		t.Fatalf("ParseMarketTime empty = (%d, %v), want (0, false)", got, ok)
	}
	if got, ok := ParseMarketTime("2024-99-99", MarketTZ.CN); ok || got != 0 {
		t.Fatalf("ParseMarketTime invalid date = (%d, %v), want (0, false)", got, ok)
	}
	if got, ok := ParseMarketTime("2024-06-13", "Bad/Zone"); ok || got != 0 {
		t.Fatalf("ParseMarketTime invalid tz = (%d, %v), want (0, false)", got, ok)
	}
}
