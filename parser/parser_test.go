package parser

import "testing"

func TestDecodeGBK(t *testing.T) {
	got, err := DecodeGBK([]byte{0xb9, 0xf3, 0xd6, 0xdd, 0xc3, 0xa9, 0xcc, 0xa8})
	if err != nil {
		t.Fatal(err)
	}
	if got != "贵州茅台" {
		t.Fatalf("DecodeGBK = %q, want 贵州茅台", got)
	}
}

func TestParseResponse(t *testing.T) {
	items := ParseResponse(` v_s_sh600519="1~贵州茅台~600519"; bad; v_pv_none_match="1"; `)
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if items[0].Key != "s_sh600519" {
		t.Fatalf("key = %q, want s_sh600519", items[0].Key)
	}
	if len(items[0].Fields) != 3 || items[0].Fields[1] != "贵州茅台" {
		t.Fatalf("fields = %#v", items[0].Fields)
	}
	if items[1].Key != "pv_none_match" || items[1].Fields[0] != "1" {
		t.Fatalf("none match = %#v", items[1])
	}
}

func TestNumberHelpers(t *testing.T) {
	if got := SafeNumber("12.34"); got != 12.34 {
		t.Fatalf("SafeNumber = %v, want 12.34", got)
	}
	if got := SafeNumber(""); got != 0 {
		t.Fatalf("SafeNumber blank = %v, want 0", got)
	}
	if got := SafeNumber("bad"); got != 0 {
		t.Fatalf("SafeNumber bad = %v, want 0", got)
	}

	if got := SafeNumberOrNil("5.5"); got == nil || *got != 5.5 {
		t.Fatalf("SafeNumberOrNil = %v, want 5.5", got)
	}
	if got := SafeNumberOrNil(""); got != nil {
		t.Fatalf("SafeNumberOrNil blank = %v, want nil", got)
	}

	if got := ToNumber("-"); got != nil {
		t.Fatalf("ToNumber dash = %v, want nil", got)
	}
	if got := ToNumber("7.25"); got == nil || *got != 7.25 {
		t.Fatalf("ToNumber = %v, want 7.25", got)
	}
	if got := ToNumberSafe(8); got == nil || *got != 8 {
		t.Fatalf("ToNumberSafe int = %v, want 8", got)
	}
	if got := ToNumberSafe(nil); got != nil {
		t.Fatalf("ToNumberSafe nil = %v, want nil", got)
	}
}
