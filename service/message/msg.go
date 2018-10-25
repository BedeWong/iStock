package message

// 撤单消息体
type MsgRevokeOrder struct {
	ID   		uint	`json:"id"`
	Order_id	uint	`json:"order_id"`
	User_id		int		`json:"user_id"`
	Stock_code	string	`json:"stock_code"`
	Trade_type  int		`json:"trade_type"`
	Req_time	string		`json:"req_time"`
}

//  添加 stock code 到数据源 的处理
type MsgSourceStockDeal struct {
	ID   		uint			`json:"id"`
	Type 		int				`json:"type"`
	Stock_code  string			`json:"stock_code"`
	Stock_name	string			`json:"stock_name"`
	Req_time	string			`json:"req_time"`
}

const (
	MsgSourceStockDealType_ADD = iota
	MsgSourceStockDealType_Del
)

// 逐笔成交 数据
type MsgTickData struct {
	ID 		uint		`json:"id"`
	Tick_time 	string	`json:"tick_time"`
	Tick_code   string  `json:"tick_code"`
	Tick_count	int		`json:"tick_count"`
	Tick_price  float64		`json:"tick_price"`
	Tick_Type 	string 		`json:"tick_type"`			// 交易类型： 买盘， 卖盘， 中性盘
}