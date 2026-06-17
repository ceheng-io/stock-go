package services

import (
	"context"
	"testing"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
)

type fundClientStub struct {
	text string
}

func (f *fundClientStub) GetText(_ context.Context, _ string) (string, error) {
	return f.text, nil
}

func TestFundServiceEstimate(t *testing.T) {
	client := &fundClientStub{text: `jsonpgz({"fundcode":"110011","name":"易方达中小盘","dwjz":"3.5000","gsz":"3.5600","gszzl":"1.71"});`}
	service := NewFundService(client, FundURLs{GZ: "https://fundgz.test/js"})

	row, err := service.Estimate(context.Background(), "110011")
	if err != nil {
		t.Fatal(err)
	}
	if row.Code != "110011" || row.EstimatedNav == nil || *row.EstimatedNav != 3.56 {
		t.Fatalf("row = %+v", row)
	}
}

func TestFundServiceNavAndRankHistory(t *testing.T) {
	client := &fundClientStub{text: `
var fS_code = "110011";
var fS_name = "易方达中小盘";
var Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"equityReturn":"1.2","unitMoney":""}];
var Data_ACWorthTrend = [[1702857600000,5.6]];
`}
	service := NewFundService(client, FundURLs{GZ: "https://fundgz.test/js", Pingzhong: "https://fund.test/pingzhongdata"})

	nav, err := service.NavHistory(context.Background(), "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(nav.Items) != 1 || nav.Items[0].AccNav == nil || *nav.Items[0].AccNav != 5.6 {
		t.Fatalf("nav = %+v", nav)
	}

	client.text = `
var fS_code = "110011";
var fS_name = "易方达中小盘";
var Data_rateInSimilarType = [{"x":1702857600000,"y":"12","sc":"300"}];
var Data_rateInSimilarPersent = [[1702857600000,4.0]];
`
	rank, err := service.RankHistory(context.Background(), "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(rank.Items) != 1 || rank.Items[0].Rank == nil || *rank.Items[0].Rank != 12 {
		t.Fatalf("rank = %+v", rank)
	}
}

func TestFundServiceDividendList(t *testing.T) {
	client := &fundClientStub{text: `var pageinfo = [1,20,1]; var jjfh_data = [["110011","易方达中小盘","2024-12-16","2024-12-17","0.12","2024-12-18","混合型"]];`}
	service := NewFundService(client, FundURLs{DataIndex: "https://fund.test/Data/funddataIndex_Interface.aspx"})

	result, err := service.DividendList(context.Background(), eastmoney.FundDividendListOptions{Year: "2024"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 || result.Items[0].Code != "110011" {
		t.Fatalf("result = %+v", result)
	}
}
