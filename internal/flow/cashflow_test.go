package flow

import (
	"testing"
	"time"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

func TestListCashFlows(t *testing.T) {
	rail := miso.EmptyRail()
	if err := miso.LoadConfigFromFile("../../conf.yml", rail); err != nil {
		t.Fatal(err)
	}
	miso.SetLogLevel("debug")
	miso.InitMySQLFromProp(rail)
	LoadCategoryConfs(rail)

	l, err := ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("1. l: %+v", l)

	l, err = ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{Direction: "OUT"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("2. l: %+v", l)

	l, err = ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{Direction: "IN"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("3. l: %+v", l)

	l, err = ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransId: "123"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("4. l: %+v", l)

	l, err = ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransId: "444"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("5. l: %+v", l)

	l, err = ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{Category: "WECHAT"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("6. l: %+v", l)

	tt := util.ETime(time.Now().Add(-time.Hour * 24))
	l, err = ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransTimeStart: &tt})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("7. l: %+v", l)

	tt = util.ETime(time.Now().Add(time.Hour * 24))
	l, err = ListCashFlows(rail, miso.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransTimeStart: &tt})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("8. l: %+v", l)
}
