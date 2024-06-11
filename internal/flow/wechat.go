package flow

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

const (
	WechatCategory = "WECHAT"
	WechatCurrency = "CNY"
)

func ParseWechatCashflows(rail miso.Rail, path string, user common.User) ([]SaveCashflowParam, error) {

	f, err := util.ReadWriteFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %v, %w", path, err)
	}
	defer f.Close()

	params := make([]SaveCashflowParam, 0, 30)
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
			t, err := time.Parse("2006-01-02 15:04:05", stranTime)
			if err != nil {
				rail.Errorf("failed to parse transaction time: '%v', %v", stranTime, err)
			} else {
				tranTime = util.ETime(t)
			}

			extram := map[string]string{}
			extram["交易类型"] = mapTryGet(titleMap, "交易类型", l)
			extram["交易类型"] = mapTryGet(titleMap, "交易类型", l)
			extrav, _ := util.SWriteJson(extram)

			amtv := mapTryGet(titleMap, "金额(元)", l)
			amtv, _ = strings.CutPrefix(amtv, "¥")

			p := SaveCashflowParam{
				UserNo:       user.UserNo,
				Direction:    dir,
				TransTime:    tranTime,
				TransId:      mapTryGet(titleMap, "交易单号", l),
				Counterparty: mapTryGet(titleMap, "交易对方", l),
				Amount:       amtv,
				Currency:     WechatCurrency,
				Extra:        extrav,
				Category:     WechatCategory,
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
