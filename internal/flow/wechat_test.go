package flow

import (
	"testing"

	"github.com/curtisnewbie/miso/miso"
)

func TestParseWechatCashflows(t *testing.T) {
	rail := miso.EmptyRail()
	if err := miso.LoadConfigFromFile("../../conf.yml", rail); err != nil {
		t.Fatal(err)
	}
	p, err := ParseWechatCashflows(rail, "")
	if err != nil {
		t.Fatal(err)
	}
	for i, l := range p {
		t.Logf("%d - %+v", i, l)
	}

	miso.InitMySQLFromProp(rail)
	err = SaveCashflows(rail, miso.GetMySQL(), p, "")
	if err != nil {
		t.Fatal(err)
	}
}
