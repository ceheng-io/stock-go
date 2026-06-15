package stock

import "testing"

func TestNewExposesIndicatorService(t *testing.T) {
	client := New()

	if client.Indicator == nil {
		t.Fatal("client.Indicator is nil")
	}
}
