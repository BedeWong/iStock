package model

import (
	"github.com/jinzhu/gorm"
	"github.com/BedeWong/iStock/db"
)


// 订单状态
const (
	//  待完成
	ORDER_STATUS_TBD = iota
	//  完成
	ORDER_STATUS_FINISH
	//  撤销
	ORDER_STATUS_REVOKE
)


// 交易类型
const (
	TRADE_TYPE_BUY = iota
	TRADE_TYPE_SALE
)


// 用户的订单表，每次下单生成该记录，
// 当订单部分成交时，不更改表记录
// 完成时更改记录状态： 待成交， 已完成， 已撤单
type Tb_order struct {
	gorm.Model
	User_id 	int 		`gorm:"not null"`

	Stock_name  string  	`grom:"type:varchar(16); not null"`
	Stock_code  string  	`grom:"type:varchar(16); not null"`
	Stock_price float64 	`gorm:"type:decimal(12,2); not null"`
	Stock_count int			`grom:"default:0;"`

	Transfer_fee  float64	`grom:"type:decimal(12,2);default:0.0"`     //  交易 过户费
	Brokerage	float64		`gorm:"type:decimal(12,2);default:0.0"`	 // 佣金
	// 冻结金额： 买入：存放冻结金额，  卖出此值为0
	Freeze_amount float64   `gorm:"default:0"`

	// 交易类型：  0买入   1 卖出
	Trade_type int			`grom:not null`
	// 交易类型： 文字描述； eg 买入， 卖出
	Trade_type_desc string	`gorm:"type:varchar(32); not null"`

	// 订单状态
	Order_status int		`gorm:"default:0;"`
}


func init() {
	if db.DBSession.HasTable(&Tb_order{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_order{})
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_order{})
	}
}