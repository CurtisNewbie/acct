package web

import (
	"net/http"

	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
)

func RegisterEndpoints(rail miso.Rail) {
	miso.GroupRoute("/open/api/v1",
		miso.IPost("/cashflow/list", ApiListCashFlows),
		miso.RawPost("/cashflow/import/wechat", ApiImportWechatCashflows),
	)
}

func ApiListCashFlows(inb *miso.Inbound, req flow.ListCashFlowReq) (miso.PageRes[flow.ListCashFlowRes], error) {
	return flow.ListCashFlows(inb.Rail(), miso.GetMySQL(), common.GetUser(inb.Rail()), req)
}

func ApiImportWechatCashflows(inb *miso.Inbound) {
	err := flow.ImportWechatCashflows(inb, miso.GetMySQL())
	if err != nil {
		inb.Errorf("Failed to import wechat cashflows, %v", err)
		inb.Status(http.StatusInternalServerError)
	}
}
