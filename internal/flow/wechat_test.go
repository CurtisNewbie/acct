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
	p, err := ParseWechatCashflows(rail, "../../testdata/wechat_test.csv")
	if err != nil {
		t.Fatal(err)
	}
	for i, l := range p {
		t.Logf("%d - %+v", i, l)
	}
}
