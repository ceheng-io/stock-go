package types

import "testing"

func TestDragonTigerPeriodConstantsMatchTSUnion(t *testing.T) {
	values := map[DragonTigerPeriod]string{
		DragonTigerPeriodOneMonth:   "1month",
		DragonTigerPeriodThreeMonth: "3month",
		DragonTigerPeriodSixMonth:   "6month",
		DragonTigerPeriodOneYear:    "1year",
	}

	assertStringEnumValues(t, values)
}

func TestDragonTigerDateOptionsMatchesTSShape(t *testing.T) {
	options := DragonTigerDateOptions{StartDate: "20241201", EndDate: "20241231"}
	if options.StartDate != "20241201" || options.EndDate != "20241231" {
		t.Fatalf("options = %+v", options)
	}
}
