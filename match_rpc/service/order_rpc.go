package service

import (
	"fmt"
	"github.com/BedeWong/iStock/service/account"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/service/order"
	manager "github.com/BedeWong/iStock/service"
	"github.com/BedeWong/iStock/service/message"
)

// rpc service
type OrderService struct {}

// 委托下单的 请求体
type AddOrderRequest struct {
	User_id		int		`json:"user_id"`
	// 代替用户的token, userid+reqtime 的hash值
	Sign		string	`json:"sign"`
	Stock_code	string	`json:"stock_code"`
	Stock_name	string	`json:"stock_name"`
	Stock_count	int		`json:"stock_count"`
	Stock_price	float64	`json:"stock_price"`
	Trade_type	int		`json:"trade_type"`
	Req_time	string  `json:"req_time"`
	// 比赛id： 默认非比赛情况
	Contest_id int          `json:"default:0"`
}

// 委托下单的响应
type AddOrderResponse struct {
	Err_msg 	string		`json:"err_msg"`
	Ret_code	int			`json:"ret_code"`
}


// 委托下单
func (this *OrderService)AddOrder(req AddOrderRequest, resp *AddOrderResponse) error {
	log.Debug("req: %#v", req)

	defer func(){
		err := recover()
			if err != nil {
				resp.Ret_code = -1
				resp.Err_msg = fmt.Sprintf("AddOrder err:%#v", err)
			}
	}()

	acc := account.Handler{}
	ok, err := acc.CheckIdentity(
		fmt.Sprintf("%d", req.User_id),
		req.Req_time,
		req.Sign)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("AddOrder err:%#v", err)
		return err
	}

	if ok != true {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("身份验证错误。")
		return nil
	}

	amount, stamp_tax, transfer_tax, brokerage, err :=
		acc.CalcTax(req.User_id,
					req.Trade_type,
					req.Stock_code,
					req.Stock_name,
					req.Stock_price,
					req.Stock_count,
					)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("服务器错误。")
		return nil
	}
	log.Info("amount:%f, stamp_tax:%f, transfer_tax:%f, brokerage:%f",
		amount, stamp_tax, transfer_tax, brokerage)

	// 扣算 金额， 税费
	if req.Contest_id == 0 {
		err = acc.DeductUserTax(req.User_id,
			amount,
			stamp_tax,
			transfer_tax,
			brokerage,
		)
	}else {
		// 比赛的情况下.
		err = acc.DeductUserContestTax(req.User_id,
			amount,
			stamp_tax,
			transfer_tax,
			brokerage,
		)
	}
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = err.Error()
		return nil
	}

	// 生成订单
	order_detail, err := order.NewOrder(req.User_id,
										req.Trade_type,
										req.Stock_code,
										req.Stock_name,
										req.Stock_price,
										req.Stock_count,
										req.Contest_id,
										amount,
										stamp_tax,
										transfer_tax,
										brokerage,
										)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = err.Error()
		return nil
	}

	// 冻结用户资产
	err = order.FreezeUserStock(order_detail)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = err.Error()
		return nil
	}

	// 将订单发送到 定序模块
	manager.Send2Senquence(order_detail, 2)
	// return ok
	resp.Ret_code = 0
	resp.Err_msg = "操作成功"
	return nil
}


// 委托撤单的 请求体
type RevokeOrderRequest struct {
	User_id		int		`json:"user_id"`
	// 代替用户的token, userid+reqtime 的hash值
	Sign		string	`json:"sign"`
	Order_id	uint	`json:"order_id"`
	Req_time	string  `json:"req_time"`

	// 比赛id： 默认非比赛情况
	Contest_id int          `json:"default:0"`
}

// 委托撤单的响应
type RevokeOrderResponse struct {
	Err_msg 	string		`json:"err_msg"`
	Ret_code	int			`json:"ret_code"`
}


func (this *OrderService)RevokeOrder(req RevokeOrderRequest,
	resp *RevokeOrderResponse) error {
	fmt.Println("req:", req)

	defer func(){
		err := recover()
		if err != nil {
			resp.Ret_code = -1
			resp.Err_msg = fmt.Sprintf("RevokeOrder err:%#v", err)
		}
	}()

	// 验证身份这段代码可以提取出来， 作个钩子函数（中间件）
	acc := account.Handler{}
	ok, err := acc.CheckIdentity(fmt.Sprintf("%d", req.User_id),
								req.Req_time,
								req.Sign)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("RevokeOrder err:%#v", err)
		return nil
	}
	if ok != true {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("身份验证错误。")
		return nil
	}

	msg := message.MsgRevokeOrder{
		Order_id: req.Order_id,
		User_id: req.User_id,
		Req_time: req.Req_time,
	}

	// 撤销订单 发送 到定序系统 从队列中删除委托订单
	manager.Send2Senquence(msg, 2)
	// return ok
	resp.Ret_code = 0
	resp.Err_msg = "操作成功"
	return nil
}
