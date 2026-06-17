package types

import (
	"reflect"
	"testing"
)

func TestKlineTimeUsesTransportTimestamp(t *testing.T) {
	if got := reflect.TypeOf(Kline{}.Time).Kind(); got != reflect.Int64 {
		t.Fatalf("Kline.Time kind = %s, want int64", got)
	}
}
