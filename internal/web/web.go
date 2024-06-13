package web

import (
	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
)

const (
	CodeManageCashflows = "acct:ManageCashflows"
)

func RegisterEndpoints(rail miso.Rail) {
	common.LoadBuiltinPropagationKeys()
	auth.ExposeResourceInfo([]auth.Resource{{
		Code: CodeManageCashflows,
		Name: "Manage Personal Cashflows",
	}})

	miso.GroupRoute("/open/api/v1",
		miso.IPost("/cashflow/list", ApiListCashFlows).Resource(CodeManageCashflows),
		miso.Post("/cashflow/import/wechat", ApiImportWechatCashflows).Resource(CodeManageCashflows),
		miso.IPost("/cashflow/statistics/calculate", ApiCalcCashflowStats).Resource(CodeManageCashflows),
	)
}

func ApiListCashFlows(inb *miso.Inbound, req flow.ListCashFlowReq) (miso.PageRes[flow.ListCashFlowRes], error) {
	return flow.ListCashFlows(inb.Rail(), miso.GetMySQL(), common.GetUser(inb.Rail()), req)
}

func ApiImportWechatCashflows(inb *miso.Inbound) (any, error) {
	return nil, flow.ImportWechatCashflows(inb, miso.GetMySQL())
}

func ApiCalcCashflowStats(inb *miso.Inbound, req flow.ApiCalcCashflowStatsReq) (any, error) {
	user := common.GetUser(inb.Rail())
	return nil, flow.CalcCsahflowStats(inb.Rail(), miso.GetMySQL(), req, user)
}
