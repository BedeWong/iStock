package service

import (
	"fmt"
)

// rpc service
type Order struct {}

// 委托下单的 请求体
type AddOrderRequest struct {
	User_id		int		`json:"user_id"`
	Sign		string	`json:"sign"`					// 代替用户的token, userid+reqtime 的hash值
	Stock_code	string	`json:"stock_code"`
	Stock_name	string	`json:"stock_name"`
	Stock_count	string	`json:"stock_count"`
	Stock_price	string	`json:"stock_price"`
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



	return nil
}


func (this *Order)RevokeOrder(req string, resp *string) error {
	fmt.Println("req:", req)
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