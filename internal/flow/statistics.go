package flow

const (
	AggTypeYearly  = "YEARLY"
	AggTypeMonthly = "MONTHLY"
	AggTypeWeekly  = "WEEKLY"
)

type ApiCalcCashflowStatsReq struct {
	AggType  string `desc:"Aggregation Type." valid:"member:YEARLY|MONTLY|WEEKLY"`
	AggRange string `desc:"Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD)." valid:"notEmpty"`
}
