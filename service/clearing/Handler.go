package clearing

import (
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BideWong/iStock/model"
	"github.com/BideWong/iStock/db"
	"github.com/BideWong/iStock/utils"
	"github.com/BideWong/iStock/conf"

	manager "github.com/BideWong/iStock/service"
)

// 清算系統
// 處理所有的 訂單的結算

// 與外界通信的 channel
//var clearing_que chan interface{}
//
//func GetChan() (chan interface{}){
//	return clearing_que
//}


func Handler(task interface{}) error{
	switch item := task.(type) {
	default:
		log.Info("task type can not handler: %T, %#v", task, task)
	case model.Tb_order_real:
		// 订单清算
		OrderHandler(item)
	case model.Tb_trade_detail:
		// 交易明细清算
		OrderDetailHandler(item)
	}

	return nil
}

// 订单处理
func OrderHandler(order_real model.Tb_order_real) {
	switch order_real.Order_status{
	default:
		log.Error("order.status not match handler. status=%d", order_real.Order_status)
	case model.ORDER_STATUS_FINISH:
		// 订单完成清算
		FinishOrderHandler(order_real)
	case model.ORDER_STATUS_REVOKE:
		// 订单撤单清算
		RevokeOrderHandler(order_real)
	}
}

// 撤单 清算处理
func RevokeOrderHandler(order_real model.Tb_order_real) {

	user := model.Tb_user{}

	err := db.DBSession.Where("user_id=?", order_real.User_id).First(&user).Error
	if err != nil {
		log.Error("user_id=%d 数据记录不存在.", order_real.User_id)
		return
	}

	if order_real.Trade_type == model.TRADE_TYPE_BUY {
		// 买订单撤单：
		//  冻结 金额， 佣金， 解冻

		freeze_money := order_real.Stock_price * float64(order_real.Stock_count)
		freeze_money = utils.Decimal(freeze_money, 2)  // 保留两位小数

		user.User_money += freeze_money
		user.User_all_capital += freeze_money

		db.DBSession.Save(&user)
	} else if order_real.Trade_type == model.TRADE_TYPE_SALE {

	}
}

// 成功 清算处理
func FinishOrderHandler(order_real model.Tb_order_real) {
	// 应该不会执行到这里

	log.Warn("FinishOrderHandler order:%#v", order_real)

	if order_real.Trade_type == model.TRADE_TYPE_BUY {
		// 买单 完成
	}else if order_real.Trade_type == model.TRADE_TYPE_SALE {
		// 卖单 完成
	}
}

// 订单明细 交易清算
func OrderDetailHandler(detail model.Tb_trade_detail) {

	if detail.Trade_type == model.TRADE_TYPE_SALE {
		// 步骤：
		//      計算印花稅， 保存交易詳細記錄到數據庫
		// 		計算用戶資產，资金回账户
		//      修改用戶的持股

		/************* 1) 計算印花稅， 保存交易詳細記錄到數據庫  */
		// 计算 总成交额
		trade_vol := detail.Stock_price * (float64)(detail.Stock_count)
		trade_vol = utils.Decimal(trade_vol, 2)
		// 计算 印花税
		detail.Stamp_tax = trade_vol * conf.Data.Trade.StampTax
		detail.Stamp_tax = utils.Decimal(detail.Stamp_tax, 2)

		// 保存新的紀錄到數據庫
		db.DBSession.Save(&detail)

		/************* 2) 計算用戶資產  */
		user := model.Tb_user{}
		if err:= db.DBSession.Where("user_id = ?", detail.User_id).First(&user).Error; err != nil {
			log.Error("user_id=%d 数据记录不存在.", detail.User_id)
			return
		}

		user.User_all_capital += trade_vol - detail.Stamp_tax
		user.User_money += trade_vol - detail.Stamp_tax

		db.DBSession.Save(&user)

		/************* 3) 修改 用戶的持股 */
		user_stocks := model.Tb_user_stock{}
		if err := db.DBSession.Where(&model.Tb_user_stock{
			User_id:detail.User_id,
			Stock_code:detail.Stock_code,
		}).First(&user_stocks).Error; err != nil {
			log.Error("Tb_user_stock select err:", err)
			return
		}

		// 修改 持倉股數
		user_stocks.Stock_count_can_sale -= detail.Stock_count
		// 修改持倉成本價
		// nothing
		db.DBSession.Save(&user_stocks)
	}else if detail.Trade_type == model.TRADE_TYPE_BUY {
		// 买成交：无需操作资金账户， 冻结资金无需减少，
		//			订单（order_real）中已经减少 委托股数， 撤单时以委托股数为准
		//			修改 用戶的持股

		// 1) 保存新的紀錄到數據庫
		db.DBSession.Save(&detail)

		// 2) 修改 持倉股數 和 持倉價格
		user_stocks := model.Tb_user_stock{}
		if err := db.DBSession.Where(&model.Tb_user_stock{
			User_id:detail.User_id,
			Stock_code:detail.Stock_code,
		}).First(&user_stocks).Error; err != nil {
			log.Error("Tb_user_stock select err:", err)
			return
		}

		// 3) 修改持倉價格
		user_stocks.Stock_price = (user_stocks.Stock_price * (float64)(user_stocks.Stock_count) +
			detail.Stock_price * (float64)(detail.Stock_count)) / (float64)(user_stocks.Stock_count + detail.Stock_count)
		// 取兩位小數
		user_stocks.Stock_price = utils.Decimal(user_stocks.Stock_price, 2)

		// 修改持倉股數
		user_stocks.Stock_count += detail.Stock_count

		// 存入數據庫
		db.DBSession.Save(&user_stocks)
	}
}


// 初始化函數
func Init() {
	clearing_que := manager.GetInstance().Clear_que

	go func() {
		for {
			task := <-clearing_que
			log.Info("recv a new task: %T, %#v", task, task)

			Handler(task)
		}
	}()

	log.Info("Clearing Init ok.")
}