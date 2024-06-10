package flow

import (
	"os"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

var (
	importPool = util.NewAsyncPool(4, 50)
)

const (
	DirectionIn  = "IN"
	DirectionOut = "OUT"
)

type ListCashFlowReq struct {
	Paging         miso.Paging `desc:"Paging"`
	Direction      string      `desc:"Flow Direction: IN / OUT" valid:"member:IN|OUT|"`
	TransTimeStart *util.ETime `desc:"Transaction Time Range Start"`
	TransTimeEnd   *util.ETime `desc:"Transaction Time Range End"`
	TransId        string      `desc:"Transaction ID"`
	Category       string      `desc:"Category"`
}

type ListCashFlowRes struct {
	Direction    string     `desc:"Flow Direction: IN / OUT"`
	TransTime    util.ETime `desc:"Transaction Time"`
	TransId      string     `desc:"Transaction ID"`
	Counterparty string     `desc:"Counterparty of the transaction"`
	Amount       string     `desc:"Amount"`
	Currency     string     `desc:"Currency"`
	Extra        string     `desc:"Extra Information"`
	Category     string     `desc:"Category"`
	Remark       string     `desc:"Remark"`
	CTime        util.ETime `desc:"Create Time"`
}

func ListCashFlows(rail miso.Rail, db *gorm.DB, user common.User, req ListCashFlowReq) (miso.PageRes[ListCashFlowRes], error) {
	return miso.NewPageQuery[ListCashFlowRes]().
		WithPage(req.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Table(`cashflow`).
				Where("user_no = ?", user.UserNo).
				Where("deleted = 0")
			if req.Direction != "" {
				tx = tx.Where("direction = ?", req.Direction)
			}
			if req.TransId != "" {
				tx = tx.Where("trans_id = ?", req.TransId)
			}
			if req.Category != "" {
				tx = tx.Where("category = ?", req.Category)
			}
			if req.TransTimeStart != nil {
				tx = tx.Where("trans_time >= ?", req.TransTimeStart)
			}
			if req.TransTimeEnd != nil {
				tx = tx.Where("trans_time <= ?", req.TransTimeEnd)
			}
			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select("direction", "trans_time", "trans_id", "counterparty",
				"amount", "currency", "extra", "category", "remark", "ctime")
		}).
		Exec(rail, db)
}

func ImportWechatCashflows(inb *miso.Inbound, db *gorm.DB) (any, error) {
	rail := inb.Rail()
	_, r := inb.Unwrap()
	user := common.GetUser(rail)
	path, err := util.SaveTmpFile("/tmp", r.Body)
	if err != nil {
		return nil, err
	}
	importPool.Go(func() {
		defer os.Remove(path)
		rail := rail.NextSpan()
		rec, err := ParseWechatCashflows(rail, path, user)
		if err != nil {
			rail.Errorf("failed to parse wechat cashflows for %v, %v", user.Username, err)
			return
		}
		if err := SaveCashflows(rail, db, rec); err != nil {
			rail.Errorf("failed to save wechat cashflows for %v, %v", user.Username, err)
		}
		rail.Infof("Wechat cashflows (%d records) saved for %v", len(rec), user.Username)
	})

	return nil, nil
}

type SaveCashflowParam struct {
	UserNo       string
	Direction    string
	TransTime    util.ETime
	TransId      string
	Counterparty string
	Amount       string
	Currency     string
	Extra        string
	Category     string
	Remark       string
	Ctime        util.ETime
}

func SaveCashflows(rail miso.Rail, db *gorm.DB, records []SaveCashflowParam) error {
	if len(records) < 1 {
		return nil
	}
	return db.Table("cashflow").CreateInBatches(records, 200).Error
}
