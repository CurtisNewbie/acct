package flow

import (
	"time"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

const (
	AggTypeYearly  = "YEARLY"
	AggTypeMonthly = "MONTHLY"
	AggTypeWeekly  = "WEEKLY"
)

var (
	RangeFormatMap = map[string]string{
		AggTypeYearly:  `2006`,
		AggTypeMonthly: `200601`,
		AggTypeWeekly:  `20060102`,
	}
)

type ApiCalcCashflowStatsReq struct {
	AggType  string `desc:"Aggregation Type." valid:"member:YEARLY|MONTLY|WEEKLY"`
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)." valid:"notEmpty"`
}

func (r ApiCalcCashflowStatsReq) Validate() (time.Time, error) {
	pat, ok := RangeFormatMap[r.AggType]
	if !ok {
		return time.Time{}, miso.NewErrf("Invalid AggType")
	}

	parsed, err := time.ParseInLocation(pat, r.AggRange, time.Local)
	if err != nil {
		return time.Time{}, miso.NewErrf("Invalid AppRange '%s' for %s aggregate type", r.AggRange, r.AggType).
			WithInternalMsg("%v", err)
	}
	return parsed, nil
}

func CalcCsahflowStats(rail miso.Rail, db *gorm.DB, req ApiCalcCashflowStatsReq, user common.User) error {
	_, err := req.Validate()
	if err != nil {
		return err
	}

	// TODO
	return nil
}
