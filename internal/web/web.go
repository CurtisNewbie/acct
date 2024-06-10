package web

import (
	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
)

func RegisterEndpoints(rail miso.Rail) error {
	miso.GroupRoute("/open/api/v1",
		miso.IPost("/cashflow/list", ApiListCashFlows),
		miso.RawPost("/cashflow/import/wechat", ApiImportWechatCashflows),
	)

	return nil
}

func ApiListCashFlows(inb *miso.Inbound, req flow.ListCashFlowReq) (miso.PageRes[flow.ListCashFlowRes], error) {
	return flow.ListCashFlows(inb.Rail(), miso.GetMySQL(), common.GetUser(inb.Rail()), req)
}

func ApiImportWechatCashflows(inb *miso.Inbound) {
	flow.ImportWechatCashflows(inb, miso.GetMySQL())
}
