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
type Order struct {}

// 委托下单的 请求体
type AddOrderRequest struct {
	User_id		int		`json:"user_id"`
	Sign		string	`json:"sign"`					// 代替用户的token, userid+reqtime 的hash值
	Stock_code	string	`json:"stock_code"`
	Stock_name	string	`json:"stock_name"`
	Stock_count	int		`json:"stock_count"`
	Stock_price	float64	`json:"stock_price"`
	Trade_type	int		`json:"trade_type"`
	Req_time	string  `json:"req_time"`
}

// 委托下单的响应
type AddOrderResponse struct {
	Err_msg 	string		`json:"err_msg"`
	Ret_code	int			`json:"ret_code"`
}

// 委托下单
func (this *Order)AddOrder(req AddOrderRequest, resp *AddOrderResponse) error {
	fmt.Println("req:", req)

	defer func(){
		err := recover()
			if err == nil {
				resp.Ret_code = -1
				resp.Err_msg = fmt.Sprintf("AddOrder err:%v", err)
			}
	}()

	acc := account.Handler{}
	ok, err := acc.CheckIdentity(fmt.Sprintf("%d", req.User_id), req.Req_time, req.Sign)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("AddOrder err:%v", err)
		return err
	}

	if ok != true {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("身份验证错误。")
		return nil
	}

	amount, stamp_tax, transfer_tax, brokerage, err := acc.CalcTax(req.User_id, req.Trade_type, req.Stock_code, req.Stock_name, req.Stock_price, req.Stock_count)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("服务器错误。")
		return nil
	}
	log.Info("amount:%f, stamp_tax:%f, transfer_tax:%f, brokerage:%f", amount, stamp_tax, transfer_tax, brokerage)

	// 扣算 金额， 税费
	err = acc.CheckAccountMoney(req.User_id, amount, stamp_tax, transfer_tax, brokerage)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = err.Error()
		return nil
	}

	// 生成订单
	order_detail, err := order.NewOrder(req.User_id, req.Trade_type, req.Stock_code, req.Stock_name, req.Stock_price, req.Stock_count, amount, stamp_tax, transfer_tax, brokerage)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = err.Error()
		return nil
	}

	// 将订单发送到 定序系统
	manager.Send2Senquence(order_detail, 2)

	return nil
}


// 委托下单的 请求体
type RevokeOrderRequest struct {
	User_id		int		`json:"user_id"`
	Sign		string	`json:"sign"`					// 代替用户的token, userid+reqtime 的hash值
	Order_id	uint	`json:"order_id"`
	Stock_code	string	`json:"stock_code"`
	Stock_name	string	`json:"stock_name"`
	Trade_type	int		`json:"trade_type"`
	Req_time	string  `json:"req_time"`
}

// 委托下单的响应
type RevokeOrderResponse struct {
	Err_msg 	string		`json:"err_msg"`
	Ret_code	int			`json:"ret_code"`
}

func (this *Order)RevokeOrder(req RevokeOrderRequest, resp *RevokeOrderResponse) error {
	fmt.Println("req:", req)

	defer func(){
		err := recover()
		if err == nil {
			resp.Ret_code = -1
			resp.Err_msg = fmt.Sprintf("RevokeOrder err:%v", err)
		}
	}()

	// 验证身份这段代码可以提取出来， 作个钩子函数（中间件）
	acc := account.Handler{}
	ok, err := acc.CheckIdentity(fmt.Sprintf("%d", req.User_id), req.Req_time, req.Sign)
	if err != nil {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("RevokeOrder err:%v", err)
		return nil
	}
	if ok != true {
		resp.Ret_code = -1
		resp.Err_msg = fmt.Sprintf("身份验证错误。")
		return nil
	}

	msg := message.MsgRevokeOrder{
		Order_id:req.Order_id,
		User_id:req.User_id,
		Req_time : req.Req_time,
	}

	// 撤销订单 发送 到定序系统 从队列中删除委托订单
	manager.Send2Senquence(msg, 2)

	return nil
}

//func (this *Order_service)Test_1(req interface{}, resp *string) error {
//
//	fmt.Println(req)
//
//	m := make(map[string]interface{})
//
//	m["name"] = "wong"
//	m["age"] = 21
//	m["fv"] = 3.14
//
//	res, err := json.Marshal(m)
//	if err != nil {
//		fmt.Println("json.Marshal err:", err)
//	}
//
//	*resp = (string)(res)
//
//	return nil
//}
//
//type Req1 struct {
//	 Name string 	`json:"name"`
//	 Age 	int		`json:"age"`
//}
//
//type Resp1 struct {
//	Id      int 		`json:"id"`
//	Desc    string 		`json:"desc"`
//}
//
//func (this *Order_service)Test_2(req Req1, resp *Resp1) error{
//
//	fmt.Println(req)
//
//	*resp = Resp1{123, "安排一下"}
//
//	return nil
//}

//func init() {
//	rpc.RegisterName("order", new (Order_service))
//	fmt.Println("register rpc sevice 'order' ok.")
//}