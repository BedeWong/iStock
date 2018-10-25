package message

type MsgRevokeOrder struct {
	Order_id	uint	`json:"order_id"`
	User_id		int		`json:"user_id"`
	Stock_code	string	`json:"stock_code"`
	Trade_type  int		`json:"trade_type"`
	Req_time	string		`json:"req_time"`
}
