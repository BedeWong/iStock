package account

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"fmt"
	"github.com/BedeWong/iStock/utils"
	"github.com/BedeWong/iStock/conf"
	"github.com/BedeWong/iStock/model"
	"github.com/BedeWong/iStock/db"
	"github.com/gpmgo/gopm/modules/log"
	"errors"
)

type Handler struct {
}

// 确认身份:
//  uid+req_time_stamp  的MD5值等於sign的值
func (this *Handler)CheckIdentity(uid, req_time_stamp, sign string) (bool, error){
	// (NOTE) hzq: 先关闭请求检查
	return true, nil

	h := md5.New()
	h.Write([]byte(uid + req_time_stamp))
	cipherStr := h.Sum(nil)
	result := hex.EncodeToString(cipherStr) // 输出加密结果

	if result != sign {
		return false, nil
	}

	return true, nil
}

// 计算 总金额， 印花税， 過戶費， 交易佣金
func (this *Handler)CalcTax(
	uid int,
	trade_type int,
	stock_code,
	stock_name string,
	stock_price float64,
	stock_count int,
	)(float64, float64, float64, float64, error){

	var amount, stamp_tax, transfer_tax, brokerage float64

	// 交易總金額
	amount = utils.Decimal(stock_price*(float64(stock_count)), 2) // 保留兩位小數

	// 賣出 需要計算印花稅
	//if trade_type == model.TRADE_TYPE_SALE {
	//	stamp_tax = amount * conf.Data.Trade.StampTax
	//	stamp_tax = utils.Decimal(stamp_tax, 2)
	//}

	// 過戶費:
	if strings.ToLower(stock_code[:2]) == "sh" {
		transfer_tax = amount * conf.Data.Trade.TransferFeeSH
	}else if strings.ToLower(stock_code[:2]) == "sz" {
		transfer_tax = amount * conf.Data.Trade.TransferFeeSZ
	}
	transfer_tax = utils.Decimal(transfer_tax, 2)
	if transfer_tax < 1.0 {
		transfer_tax = 1.0
	}

	// 佣金：
	brokerage = amount * conf.Data.Trade.Brokerage
	brokerage = utils.Decimal(brokerage, 2)

	// 总金额: 卖出不做计算
	if trade_type == model.TRADE_TYPE_SALE {
		amount = 0
	}
	log.Debug(
		"CalcTax amout: %f, stamp_tax: %f, transfer_tax: %f, brokerage: %f",
		amount, stamp_tax, transfer_tax, brokerage)
	return amount, stamp_tax, transfer_tax, brokerage, nil
}

// 扣算 金額，
// 印花税： 卖方出， 每笔交易体现  -- 不在这里计算  modified 2018年10月24日22:20:24  by BedeWong
// 过户费：单次交易计算， 后续不在扣除， 撤单不在退回。
// 总金额： 卖出0 卖出时，逐笔结算； 买入时统一冻结
func (this *Handler)DeductUserTax(userid int,
	amount, stamp_tax, transfer_tax, brokerage float64) (err error) {
	user := model.Tb_user_assets{}

	found := db.DBSession.Where(
		"user_id = ?", userid).First(&user).RecordNotFound()
	if found == true {
		err = errors.New(fmt.Sprintf("用户[%d]找不到.", userid))
		log.Error("CheckAccountMoney:err:%s", err)
		return err
	}

	user.User_money -= amount + transfer_tax + brokerage
	if user.User_money <= 0 {
		err = errors.New(fmt.Sprintf("用户[%d]的余额不足.", userid))
		log.Info(
			"CheckAccountMoney:err:%s, 當前可用資產：%f", err, user.User_money)
		return err
	}

	db.DBSession.Save(&user)
	log.Debug("CheckAccountMoney user: %v", user)
	return nil
}
