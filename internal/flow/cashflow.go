package flow

import (
	"os"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

var (
	importPool    = util.NewAsyncPool(4, 50)
	categoryConfs map[string]CategoryConf
)

const (
	DirectionIn  = "IN"
	DirectionOut = "OUT"
)

type CategoryConf struct {
	Code string
	Name string
}

func LoadCategoryConfs(rail miso.Rail) {
	var cate []CategoryConf
	miso.UnmarshalFromPropKey("acct.category.builtin", &cate)
	categoryConfs = make(map[string]CategoryConf, len(cate))
	for i, v := range cate {
		categoryConfs[v.Code] = cate[i]
	}
	rail.Debugf("Loaded conf: %#v", categoryConfs)
}

type ListCashFlowReq struct {
	Paging         miso.Paging `desc:"Paging"`
	Direction      string      `desc:"Flow Direction: IN / OUT" valid:"member:IN|OUT|"`
	TransTimeStart *util.ETime `desc:"Transaction Time Range Start"`
	TransTimeEnd   *util.ETime `desc:"Transaction Time Range End"`
	TransId        string      `desc:"Transaction ID"`
	Category       string      `desc:"Category Code"`
}

type ListCashFlowRes struct {
	Direction    string     `desc:"Flow Direction: IN / OUT"`
	TransTime    util.ETime `desc:"Transaction Time"`
	TransId      string     `desc:"Transaction ID"`
	Counterparty string     `desc:"Counterparty of the transaction"`
	Amount       string     `desc:"Amount"`
	Currency     string     `desc:"Currency"`
	Extra        string     `desc:"Extra Information"`
	Category     string     `desc:"Category Code"`
	CategoryName string     `desc:"Category Name"`
	Remark       string     `desc:"Remark"`
	CreatedAt    util.ETime `desc:"Create Time"`
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
				"amount", "currency", "extra", "category", "remark", "created_at").
				Order("trans_time desc")
		}).
		ForEach(func(t ListCashFlowRes) ListCashFlowRes {
			if v, ok := categoryConfs[t.Category]; ok {
				t.CategoryName = v.Name
			}
			return t
		}).
		Exec(rail, db)
}

func ImportWechatCashflows(inb *miso.Inbound, db *gorm.DB) error {
	rail := inb.Rail()
	user := common.GetUser(rail)
	rail.Infof("User %v importing wechat cashflows", user.Username)

	_, r := inb.Unwrap()
	path, err := util.SaveTmpFile("/tmp", r.Body)
	if err != nil {
		return err
	}
	rail.Infof("Wechat cashflows saved to temp file: %v", path)

	importPool.Go(func() {
		rail := rail.NextSpan()
		defer func() {
			os.Remove(path)
			rail.Infof("Temp file removed, %v", path)
		}()

		rec, err := ParseWechatCashflows(rail, path)
		if err != nil {
			rail.Errorf("failed to parse wechat cashflows for %v, %v", user.Username, err)
			return
		}
		if len(rec) > 0 {
			if err := SaveCashflows(rail, db, rec, user.UserNo); err != nil {
				rail.Errorf("failed to save wechat cashflows for %v, %v", user.Username, err)
			}
		}
		rail.Infof("Wechat cashflows (%d records) saved for %v", len(rec), user.Username)
	})

	return nil
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
	CreatedAt    util.ETime
}

func SaveCashflows(rail miso.Rail, db *gorm.DB, records []SaveCashflowParam, userNo string) error {
	if len(records) < 1 {
		return nil
	}

	lock := userCashflowLock(rail, userNo)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	transIdSet := util.NewSet[string]()
	for i := range records {
		v := records[i]
		transIdSet.Add(v.TransId)
		v.UserNo = userNo
		records[i] = v
	}

	// find those that already exist and skip them
	var existingTransId []string
	err := db.Raw(`SELECT trans_id FROM cashflow WHERE user_no = ? AND trans_id IN ? AND deleted = 0`, userNo, transIdSet.CopyKeys()).
		Scan(&existingTransId).Error
	if err != nil {
		return err
	}
	for _, ti := range existingTransId {
		rail.Debugf("Transaction %v for user %v already exists, ignored", ti, userNo)
		transIdSet.Del(ti)
	}

	records = util.Filter(records, func(p SaveCashflowParam) bool { return transIdSet.Has(p.TransId) })
	return db.Table("cashflow").CreateInBatches(records, 200).Error
}

func userCashflowLock(rail miso.Rail, userNo string) *miso.RLock {
	return miso.NewRLockf(rail, "acct:cashflow:user:%v", userNo)
}
