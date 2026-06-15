package tencent

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeCodeClient struct {
	lists map[string][]string
	texts map[string]string
}

func (f fakeCodeClient) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := struct {
		Success bool     `json:"success"`
		List    []string `json:"list"`
	}{Success: true, List: f.lists[requestURL]}
	b, _ := json.Marshal(payload)
	return json.Unmarshal(b, target)
}

func (f fakeCodeClient) GetText(_ context.Context, requestURL string) (string, error) {
	return f.texts[requestURL], nil
}

func (f fakeCodeClient) AShareListURL() string { return "a" }
func (f fakeCodeClient) USListURL() string     { return "us" }
func (f fakeCodeClient) HKListURL() string     { return "hk" }
func (f fakeCodeClient) FundListURL() string   { return "fund" }

func TestGetAShareCodeListMarketFilters(t *testing.T) {
	client := fakeCodeClient{lists: map[string][]string{
		"a": {"sh600000", "sz000001", "sz300750", "bj830799", "bj870204", "bj430047", "bj920819", "sh900901"},
	}}

	bj, err := GetAShareCodeList(context.Background(), client, CodeListOptions{Market: AShareMarketBJ})
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, bj, []string{"bj830799", "bj870204", "bj430047", "bj920819"})

	sh, err := GetAShareCodeList(context.Background(), client, CodeListOptions{Market: AShareMarketSH})
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, sh, []string{"sh600000"})

	sz, err := GetAShareCodeList(context.Background(), client, CodeListOptions{Market: AShareMarketSZ, Simple: true})
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, sz, []string{"000001", "300750"})
}

func TestGetUSCodeListFiltersAndSimplifies(t *testing.T) {
	client := fakeCodeClient{lists: map[string][]string{
		"us": {"105.AAPL", "106.BABA", "107.TEST"},
	}}

	got, err := GetUSCodeList(context.Background(), client, USCodeListOptions{Market: USMarketNASDAQ, Simple: true})
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, got, []string{"AAPL"})
}

func TestGetHKAndFundCodeList(t *testing.T) {
	client := fakeCodeClient{lists: map[string][]string{
		"hk": {"00700", "09988"},
	}, texts: map[string]string{"fund": ",110011,000001"}}

	hk, err := GetHKCodeList(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, hk, []string{"00700", "09988"})

	fund, err := GetFundCodeList(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, fund, []string{"110011", "000001"})
}

func TestGetFundCodeListParsesCommaTextLikeTypeScript(t *testing.T) {
	client := fakeCodeClient{texts: map[string]string{
		"fund": ",110011,,000001, ",
	}}

	got, err := GetFundCodeList(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, got, []string{"110011", "000001"})
}

func assertStrings(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
}
