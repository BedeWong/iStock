package order

import (
	"github.com/BedeWong/iStock/model"
	"github.com/BedeWong/iStock/db"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/pkg/errors"
)

/***
	userid: 用户id
	trade_type: 交易类型
	stockc_code: 股票代码
	stock_name：股票名
	stock_price：委托价
	stock_count：委托量（股）
	amount： 总金额
	stamp_tax：印花税
	transfer_tax：过户费
	brokerage：交易佣金

添加一条新的订单记录。
在数据库中将产生两个recode， Tb_order 存储的记录不会被更改
	Tb_order_real 存储的记录会在每次的部分成交中修改其 数量，这个记录也是定序系统（sequence）中排序依据，存储依据。
 */
func NewOrder(userid int,
	trade_type int,  // 用户id
	stock_code,  // 股票代码
	stock_name string, // 股票名
	stock_price float64, // 委托价
	stock_count int,  // 委托量
	contest_id int,  // 比赛场次：0 非比赛
	amount,  // 总金额
	stamp_tax, // 印花税
	transfer_tax,  // 换手费
	brokerage float64, // 佣金
	) (model.Tb_order_real, error) {

	order := model.Tb_order{
		User_id : userid,
		Stock_name : stock_name,
		Stock_code : stock_code,
		Stock_price : stock_price,
		Stock_count : stock_count,
		Transfer_fee : transfer_tax,
		Brokerage : brokerage,
		Trade_type : trade_type,
		Trade_type_desc: "",   // 待填写
		Order_status : 0,
		Contest_id: contest_id,
	}
	if order.Trade_type == model.TRADE_TYPE_BUY {
		order.Trade_type_desc = "买入"
		order.Freeze_amount = amount
	}else if order.Trade_type == model.TRADE_TYPE_SALE {
		order.Trade_type_desc = "卖出"
	}
	// 持久化订单数据
	db.DBSession.Save(&order)

	// 保存订单实时信息
	order_real := model.Tb_order_real{
		Order_id : (int)(order.ID),
		User_id : userid,
		Stock_name : stock_name,
		Stock_code : stock_code,
		Stock_price : stock_price,
		Stock_count : stock_count,
		Order_status : 0,

		Trade_type : trade_type,
		Trade_type_desc: "",   // 待填写
		Contest_id: contest_id, // 比赛id， 0表示非比赛
	}
	if order_real.Trade_type == model.TRADE_TYPE_BUY {
		order_real.Trade_type_desc = "买入"
		order_real.Freeze_amount = amount
	}else if order_real.Trade_type == model.TRADE_TYPE_SALE {
		order_real.Trade_type_desc = "卖出"
	}

	log.Info("order_new order_detail: %#v", order_real)
	db.DBSession.Save(&order_real)

	return order_real, nil
}


// 设置订单状态完成
func SetOederStatusFinished(id int) {
	if id == 0 {
		log.Error("SetOederStatusFinished err: id: %d", id)
		return
	}

	var order = model.Tb_order{}
	var cnt = 0
	db.DBSession.Where("id=?", id).Find(&order).Count(&cnt)
	if cnt == 0 {
		log.Error("id: %d not exists.", id)
		return
	}

	order.Order_status = model.ORDER_STATUS_FINISH
	db.DBSession.Save(&order)
	return
}


// 设置订单状态撤销
func SetOederStatusRevoke(id int) {
	if id == 0 {
		log.Error("SetOederStatusFinished err: id: %d", id)
		return
	}

	var order = model.Tb_order{}
	var cnt = 0
	db.DBSession.Where("id=?", id).Find(&order).Count(&cnt)
	if cnt == 0 {
		log.Error("id: %d not exists.", id)
		return
	}

	order.Order_status = model.ORDER_STATUS_REVOKE
	db.DBSession.Save(&order)
	return
}

// 冻结用户持股
//
// order_real: 订单信息
func FreezeUserStock(order_real model.Tb_order_real) error {
	if order_real.Trade_type != model.TRADE_TYPE_BUY {
		log.Debug("order Handler freezeUserStock 非买单不以处理.")
		return nil
	}

	count := order_real.Stock_count
	user_id := order_real.User_id
	stock_code := order_real.Stock_code

	if order_real.Contest_id > 0 {
		pos := model.Tb_user_contest_position{}
		// select 数据
		err := db.DBSession.Where("user_id=? and stock_code=?",
			user_id, stock_code).First(&pos).Error
		if err != nil {
			log.Error("order Handler freezeUserStock err:%v", err)
			return errors.New("服务器开了个小差.")
		}

		if pos.Stock_count_can_sale < count {
			log.Info("order Handler freezeUserStock 用户持股数不足：" +
				"user_stock: %d, sale: %d", pos.Stock_count_can_sale, count)
			return errors.New("没有足够的股数可供卖出.")
		}
		pos.Stock_count_can_sale -= count
		db.DBSession.Save(&pos)
	}else {
		pos := model.Tb_user_position{}
		// select 数据
		err := db.DBSession.Where("user_id=? and stock_code=?",
			user_id, stock_code).First(&pos).Error
		if err != nil {
			log.Error("order Handler freezeUserStock err:%v", err)
			return errors.New("服务器开了个小差.")
		}
		if pos.Stock_count_can_sale < count {
			log.Info("order Handler freezeUserStock 用户持股数不足：" +
				"user_stock: %d, sale: %d", pos.Stock_count_can_sale, count)
			return errors.New("没有足够的股数可供卖出.")
		}
		pos.Stock_count_can_sale -= count
		db.DBSession.Save(&pos)
	}

	return nil
}
