package model

import (
	"github.com/jinzhu/gorm"
	"github.com/BedeWong/iStock/db"
)


// 用户的订单表，每次下单生成该记录，
// 当订单部分成交时，更改记录，修改剩余股数量
// 完成时更改记录状态： 待成交， 已完成， 已撤单
type Tb_order_real struct {
	gorm.Model
	Order_id    int 		`gorm:"not null"`				// 外键引用： order 表
	User_id 	int 		`gorm:"not null"`
	Stock_name  string  	`grom:"type:varchar(16); not null"`
	Stock_code  string  	`grom:"type:varchar(16); not null"`
	Stock_price float64 	`gorm:"type:decimal(12,2); not null"`
	Stock_count int			`grom:"default:0;"`
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
	if db.DBSession.HasTable(&Tb_order_real{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_order_real{})
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_order_real{})
	}
}