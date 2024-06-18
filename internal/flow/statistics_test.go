package flow

import (
	"testing"
	"time"

	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

func TestReqValidate(t *testing.T) {
	tab := [][]string{
		{AggTypeYearly, "2024", "y"},
		{AggTypeYearly, "20241", "n"},
		{AggTypeMonthly, "202402", "y"},
		{AggTypeMonthly, "2024", "n"},
		{AggTypeWeekly, "20240204", "y"},
		{AggTypeWeekly, "20240203", "n"},
		{AggTypeWeekly, "2024010", "n"},
		{AggTypeWeekly, "202401", "n"},
	}
	for _, r := range tab {
		ti, err := ParseAggRangeTime(r[0], r[1])
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

func TestOnCalcCashflowStatsEvent(t *testing.T) {
	rail := miso.EmptyRail()
	if err := miso.LoadConfigFromFile("../../conf.yml", rail); err != nil {
		t.Fatal(err)
	}
	miso.SetLogLevel("debug")
	err := miso.InitMySQLFromProp(rail)
	if err != nil {
		t.Fatal(err)
	}
	_, err = miso.InitRedisFromProp(rail)
	if err != nil {
		t.Fatal(err)
	}

	tab := [][]string{
		{AggTypeYearly, "2024"},
		{AggTypeMonthly, "202406"},
		{AggTypeWeekly, "20240602"},
	}
	for _, r := range tab {
		typ := r[0]
		rng := r[1]
		ti, err := ParseAggRangeTime(typ, rng)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("rng: %v, ti: %v", rng, ti)
		err = OnCalcCashflowStatsEvent(rail, CalcCashflowStatsEvent{
			UserNo:   "UE1049787455160320075953",
			AggType:  typ,
			AggRange: rng,
			AggTime:  util.ETime(ti),
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestOnCashflowChanged(t *testing.T) {
	rail := miso.EmptyRail()
	if err := miso.LoadConfigFromFile("../../conf.yml", rail); err != nil {
		t.Fatal(err)
	}
	miso.SetLogLevel("debug")
	err := miso.InitMySQLFromProp(rail)
	if err != nil {
		t.Fatal(err)
	}
	_, err = miso.InitRedisFromProp(rail)
	if err != nil {
		t.Fatal(err)
	}
	err = rabbit.StartRabbitMqClient(rail)
	if err != nil {
		t.Fatal(err)
	}

	userNo := "UE1049787455160320075953"
	var tranTimes []util.ETime
	err = miso.GetMySQL().Raw(`SELECT trans_time FROM cashflow WHERE user_no = ?`, userNo).Scan(&tranTimes).Error
	if err != nil {
		t.Fatal(err)
	}

	err = OnCashflowChanged(rail, util.MapTo(tranTimes, func(et util.ETime) CashflowChange { return CashflowChange{et} }), userNo)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
}
