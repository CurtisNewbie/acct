package flow

import "testing"

func TestReqValidate(t *testing.T) {
	tab := [][]string{
		{AggTypeYearly, "2024", "y"},
		{AggTypeYearly, "20241", "n"},
		{AggTypeMonthly, "202402", "y"},
		{AggTypeMonthly, "2024", "n"},
		{AggTypeWeekly, "20240203", "y"},
		{AggTypeWeekly, "2024010", "n"},
		{AggTypeWeekly, "202401", "n"},
	}
	for _, r := range tab {
		req := ApiCalcCashflowStatsReq{
			AggType:  r[0],
			AggRange: r[1],
		}
		ti, err := req.Validate()
		actual := err == nil
		expected := r[2] == "y"

		if expected != actual {
			if err != nil {
				t.Fatal(err)
			} else {
				t.Fatalf("actual: %v != expected: %v", actual, expected)
			}
		}
		if actual {
			t.Logf("Time: %v", ti)
		}
	}
}
