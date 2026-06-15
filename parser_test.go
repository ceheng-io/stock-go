package stock

import "testing"

func TestRootReExportsParserUtilities(t *testing.T) {
	text, err := DecodeGBK([]byte{0xb9, 0xf3, 0xd6, 0xdd})
	if err != nil {
		t.Fatal(err)
	}
	if text != "贵州" {
		t.Fatalf("DecodeGBK = %q, want 贵州", text)
	}

	items := ParseResponse(`v_s_sh600519="1~贵州茅台";`)
	if len(items) != 1 || items[0].Key != "s_sh600519" || items[0].Fields[1] != "贵州茅台" {
		t.Fatalf("items = %#v", items)
	}
	if SafeNumber("bad") != 0 {
		t.Fatal("SafeNumber bad did not return 0")
	}
	if number := SafeNumberOrNull("5.5"); number == nil || *number != 5.5 {
		t.Fatalf("SafeNumberOrNull = %v, want 5.5", number)
	}
	if number := SafeNumberOrNull("bad"); number != nil {
		t.Fatalf("SafeNumberOrNull bad = %v, want nil", number)
	}
	if number := ToNumber("-"); number != nil {
		t.Fatalf("ToNumber dash = %v, want nil", number)
	}
}
