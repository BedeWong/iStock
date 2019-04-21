// 清算模块，處理所有的 訂單的結算
//
// @author: bedewong
// @create at:
// @update at: 2019年4月20日
// @change log:

package clearing

import (
	"github.com/BedeWong/iStock/conf"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/model"
	"github.com/BedeWong/iStock/db"
	"github.com/BedeWong/iStock/utils"
	manager "github.com/BedeWong/iStock/service"
	"github.com/BedeWong/iStock/service/order"
)


func handleCmd(task interface{}) error{
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
		// todo..
		FinishOrderHandler(order_real)
	case model.ORDER_STATUS_REVOKE:
		// 订单撤单清算
		RevokeOrderHandler(order_real)
	}
}


// 撤销普通委托单清算
func RevokeGeneralOrder(order_real model.Tb_order_real) {
	user := model.Tb_user_assets{}

	err := db.DBSession.Where("user_id=?", order_real.User_id).First(&user).Error
	if err != nil {
		log.Error("user_id=%d 数据记录不存在.", order_real.User_id)
		return
	}

	if order_real.Trade_type == model.TRADE_TYPE_BUY {
		// 买订单撤单：
		//  冻结的印花税， 佣金， 解冻

		freeze_money := order_real.Stock_price * float64(order_real.Stock_count)
		freeze_money = utils.Decimal(freeze_money, 2)  // 保留两位小数

		user.User_money += freeze_money

		db.DBSession.Save(&user)
	} else if order_real.Trade_type == model.TRADE_TYPE_SALE {

	}

	order.SetOederStatusRevoke(order_real.Order_id)
}


func RevokeContestOrder(order_real model.Tb_order_real) {
	user_id := order_real.User_id
	contest_id := order_real.Contest_id

	if order_real.Trade_type == model.TRADE_TYPE_BUY {
		// 买订单撤单：
		//  冻结的印花税， 佣金， 解冻
		freeze_money := order_real.Stock_price * float64(order_real.Stock_count)
		freeze_money = utils.Decimal(freeze_money, 2) // 保留两位小数

		db.DBSession.Exec(
			"update tb_contest_detail set c_money=c_money+? where " +
				"u_id=? and c_id=? and c_status=0",
				freeze_money, user_id, contest_id)
	}

	order.SetOederStatusRevoke(order_real.Order_id)
}


// 撤单 清算处理
func RevokeOrderHandler(order_real model.Tb_order_real) {
	if order_real.Contest_id > 0 {
		// 比赛撤单清算
		RevokeContestOrder(order_real)
	}else {
		// 练习撤单清算
		RevokeGeneralOrder(order_real)
	}
}


// 订单完成 清算处理
func FinishOrderHandler(order_real model.Tb_order_real) {
	// 应该不会执行到这里
	log.Warn("FinishOrderHandler order:%#v", order_real)

	order.SetOederStatusFinished(order_real.Order_id)
	if order_real.Trade_type == model.TRADE_TYPE_BUY {
		// 买单 完成
	}else if order_real.Trade_type == model.TRADE_TYPE_SALE {
		// 卖单 完成
		// 扣除
	}
}


// 保存逐笔交易数据
func saveTradeDetail(detail *model.Tb_trade_detail) {
	if detail.Trade_type == model.TRADE_TYPE_SALE {
		// 计算 总成交额
		trade_vol := detail.Stock_price * (float64)(detail.Stock_count)
		trade_vol = utils.Decimal(trade_vol, 2)
		// 计算 印花税
		detail.Stamp_tax = trade_vol * conf.GetConfig().Trade.StampTax
		detail.Stamp_tax = utils.Decimal(detail.Stamp_tax, 2)

		log.Info("saveTradeDetail Sale: %.2f, Stamp_tax: %.2f", trade_vol, detail.Stamp_tax)
	}else if detail.Trade_type == model.TRADE_TYPE_BUY {
		log.Info("saveTradeDetail Buy: %#v", detail)
	}
	// 1) 保存新的紀錄到數據庫
	db.DBSession.Save(&detail)
}


// 保存用户持仓信息
func saveUserPosition(detail *model.Tb_trade_detail) {
	// 修改 用戶的持股
	user_stocks := model.Tb_user_position{}
	if detail.Trade_type == model.TRADE_TYPE_SALE {
		if err := db.DBSession.Where(&model.Tb_user_position{
			User_id:    detail.User_id,
			Stock_code: detail.Stock_code,
		}).First(&user_stocks).Error; err != nil {
			log.Error("saveUserPosition Tb_user_stock select err:", err)
			return
		}
		// 修改 持倉股數
		user_stocks.Stock_count_can_sale -= detail.Stock_count
		// 修改持倉成本價
		// nothing
	}else if detail.Trade_type == model.TRADE_TYPE_BUY {
		// 修改 持倉股數 和 持倉價格
		if err := db.DBSession.Where(&model.Tb_user_position{
			User_id:detail.User_id,
			Stock_code:detail.Stock_code,
		}).First(&user_stocks).Error; err != nil {
			log.Error("saveUserPosition Tb_user_stock select err:", err)
		}

		// 用户信息
		user_stocks.User_id = detail.User_id
		// 修改持倉價格
		// 取兩位小數
		user_stocks.Stock_price =
			(user_stocks.Stock_price * (float64)(user_stocks.Stock_count) +
				detail.Stock_price * (float64)(detail.Stock_count)) /
				(float64)(user_stocks.Stock_count + detail.Stock_count)
		user_stocks.Stock_price = utils.Decimal(user_stocks.Stock_price, 2)
		// 修改持倉股數
		user_stocks.Stock_count += detail.Stock_count
		// 股票名
		user_stocks.Stock_code = detail.Stock_code
		user_stocks.Stock_name = detail.Stock_name

		log.Info("saveUserPosition user_stocks: %#v", user_stocks)
	}
	if user_stocks.User_id > 0 {
		db.DBSession.Save(&user_stocks)
	}
}


// 用户资产清算
func saveUserAssets(detail *model.Tb_trade_detail) {
	user := model.Tb_user_assets{}
	if detail.Trade_type == model.TRADE_TYPE_SALE {
		if err:= db.DBSession.Where("user_id = ?", detail.User_id).First(&user).Error; err != nil {
			log.Error("user_id=%d 数据记录不存在.", detail.User_id)
			return
		}

		trade_vol := detail.Stock_price * (float64)(detail.Stock_count)
		trade_vol = utils.Decimal(trade_vol, 2)
		user.User_money += trade_vol - detail.Stamp_tax
	}else if detail.Trade_type == model.TRADE_TYPE_BUY {

	}
	if user.User_id > 0 {
		db.DBSession.Save(&user)
	}
}


// 常规逐笔交易清算
func orderDetailGeneralClearing(detail model.Tb_trade_detail) {
	// 持久化交易明细到数据到数据库
	saveTradeDetail(&detail)
	// 常规交易处理

	//  持久化用户持仓
	saveUserPosition(&detail)
	// 清算资产
	saveUserAssets(&detail)
}


// 保存用户比赛的持仓信息
func saveUserContestPosition(detail *model.Tb_trade_detail) {
	// 修改 用戶的持股
	user_stocks := model.Tb_user_contest_position{}
	if detail.Trade_type == model.TRADE_TYPE_SALE {
		if err := db.DBSession.Where(&model.Tb_user_contest_position{
			User_id:    detail.User_id,
			Stock_code: detail.Stock_code,
			Contest_id: detail.Contest_id,
		}).First(&user_stocks).Error; err != nil {
			log.Error("saveUserContestPosition Tb_user_contest_position select err:", err)
			return
		}
		// 修改 持倉股數
		user_stocks.Stock_count_can_sale -= detail.Stock_count
		// 修改持倉成本價
		// nothing
	}else if detail.Trade_type == model.TRADE_TYPE_BUY {
		// 修改 持倉股數 和 持倉價格
		if err := db.DBSession.Where(&model.Tb_user_contest_position{
			User_id:detail.User_id,
			Stock_code:detail.Stock_code,
		}).First(&user_stocks).Error; err != nil {
			log.Info("saveUserContestPosition Tb_user_contest_position select err:", err)
		}

		// 用户信息
		user_stocks.User_id = detail.User_id
		// 修改持倉價格
		// 取兩位小數
		user_stocks.Stock_price =
			(user_stocks.Stock_price * (float64)(user_stocks.Stock_count) +
				detail.Stock_price * (float64)(detail.Stock_count)) /
				(float64)(user_stocks.Stock_count + detail.Stock_count)
		user_stocks.Stock_price = utils.Decimal(user_stocks.Stock_price, 2)
		// 修改持倉股數
		user_stocks.Stock_count += detail.Stock_count
		// 股票名
		user_stocks.Stock_name = detail.Stock_name
		user_stocks.Stock_code = detail.Stock_code

		// 比赛场次信息
		user_stocks.Contest_id = detail.Contest_id
		log.Info("saveUserContestPosition Tb_user_contest_position: %#v", user_stocks)
	}
	if user_stocks.User_id > 0 {
		db.DBSession.Save(&user_stocks)
	}
}


// 清算用户比赛的资产信息
func saveUserContestAssets(detail *model.Tb_trade_detail) {
	type assetsInfo struct {
		c_money float64
	}
	var info = assetsInfo{}

	if detail.Trade_type == model.TRADE_TYPE_SALE {
		err := db.DBSession.Raw(
			"select c_money from tb_contest_detail where " +
				"u_id=? and c_id=? ", detail.User_id, detail.Contest_id,
		).Scan(&info).Error
		if err != nil {
			log.Error(
				"saveUserContestAssets user_id: %d, c_id: %d 数据记录不存在.",
				detail.User_id, detail.Contest_id)
			return
		}

		trade_vol := detail.Stock_price * (float64)(detail.Stock_count)
		trade_vol = utils.Decimal(trade_vol, 2)
		info.c_money += trade_vol - detail.Stamp_tax
	}else if detail.Trade_type == model.TRADE_TYPE_BUY {

	}

	err := db.DBSession.Exec(
		"update tb_contest_detail set c_money=? where u_id=? and c_id=? ",
		info.c_money, detail.User_id, detail.Contest_id,
	).Error
	if err != nil {
		log.Error("saveUserContestAssets update failed user_id: %d, c_id: %d",
			detail.User_id, detail.Contest_id)
	}
}


// 比赛逐笔交易清算
func orderDetailContestClearing(detail model.Tb_trade_detail) {
	// 持久化交易明细到数据到数据库
	saveTradeDetail(&detail)

	//  持久化用户持仓
	saveUserContestPosition(&detail)
	// 清算资产
	saveUserContestAssets(&detail)
}


// 订单逐笔明细 交易清算
func OrderDetailHandler(detail model.Tb_trade_detail) {
	if detail.Contest_id > 0 {
		orderDetailContestClearing(detail)
	}else if detail.Contest_id == 0{
		orderDetailGeneralClearing(detail)
	}else {
		log.Error("Clearing OrderDetailHandler err. detail:%#v", detail)
	}
}


// 初始化函數
func Init() {
	task_chan := manager.GetInstance().Clear_que

	go func() {
		for {
			task := <-task_chan
			log.Info("recv a new task: %T, %#v", task, task)

			handleCmd(task)
		}
	}()

	log.Info("Clearing Init ok.")
}