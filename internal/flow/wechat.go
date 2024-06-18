package flow

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/encoding"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

const (
	WechatCategory = "WECHAT"
	WechatCurrency = "CNY"
)

func ParseWechatCashflows(rail miso.Rail, path string) ([]NewCashflow, error) {

	f, err := util.ReadWriteFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %v, %w", path, err)
	}
	defer f.Close()

	params := make([]NewCashflow, 0, 30)
	titleMap := make(map[string]int, 10)
	start := false

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	for {
		l, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("failed to read csv file, %v, %w", path, err)
		}
		if len(l) < 1 {
			continue
		}
		if !start {
			first := l[0]
			if strings.Contains(first, "微信支付账单明细列表") {
				start = true
				continue
			}
		}

		if !start {
			rail.Debugf("not started yet, l: %+v", l)
			continue
		}

		if len(titleMap) < 1 {
			for i, v := range l {
				titleMap[v] = i
			}
		} else {
			var dir string
			v := mapTryGet(titleMap, "收/支", l)
			if v == "支出" {
				dir = DirectionOut
			} else {
				dir = DirectionIn
			}

			var stranTime string = mapTryGet(titleMap, "交易时间", l)
			var tranTime util.ETime
			t, err := time.ParseInLocation("2006-01-02 15:04:05", stranTime, time.FixedZone("", 8))
			if err != nil {
				rail.Errorf("failed to parse transaction time: '%v', %v", stranTime, err)
			} else {
				tranTime = util.ETime(t)
			}

			extram := map[string]string{}
			extram["交易类型"] = mapTryGet(titleMap, "交易类型", l)
			extram["商户单号"] = mapTryGet(titleMap, "商户单号", l)

			paymentMethod := mapTryGet(titleMap, "支付方式", l)
			extram["支付方式"] = paymentMethod
			extrav, _ := encoding.SWriteJson(extram)

			amtv := mapTryGet(titleMap, "金额(元)", l)
			amtv, _ = strings.CutPrefix(amtv, "¥")

			p := NewCashflow{
				Direction:     dir,
				PaymentMethod: paymentMethod,
				TransTime:     tranTime,
				TransId:       mapTryGet(titleMap, "交易单号", l),
				Counterparty:  mapTryGet(titleMap, "交易对方", l),
				Amount:        amtv,
				Currency:      WechatCurrency,
				Extra:         extrav,
				Remark:        "",
			}
			params = append(params, p)
		}
	}

	return params, nil

}

func mapTryGet(m map[string]int, s string, l []string) string {
	i, ok := m[s]
	if !ok {
		return ""
	}
	if i >= len(l) {
		return ""
	}
	return strings.TrimSpace(l[i])
}
