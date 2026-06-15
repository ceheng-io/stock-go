package stock_test

import (
	"testing"

	stock "github.com/ceheng.io/stock-go"
	stockerrors "github.com/ceheng.io/stock-go/errors"
)

func TestErrorsSubpackageReExportsTSErrorEntry(t *testing.T) {
	err := stockerrors.NewInvalidArgumentError("bad argument", map[string]any{"field": "code"})

	var sdkErr *stockerrors.SdkError = err
	var requestErr *stockerrors.RequestError = err
	var code stockerrors.SdkErrorCode = sdkErr.Code
	if code != stockerrors.CodeInvalidArgument {
		t.Fatalf("SdkErrorCode = %s, want %s", code, stockerrors.CodeInvalidArgument)
	}
	if !stockerrors.IsSdkError(requestErr) {
		t.Fatal("IsSdkError returned false")
	}
	if got := stockerrors.GetSdkErrorCode(requestErr); got != stock.CodeInvalidArgument {
		t.Fatalf("GetSdkErrorCode = %s, want %s", got, stock.CodeInvalidArgument)
	}
}
