package model

import (
	"github.com/jinzhu/gorm"
	"github.com/BideWong/iStock/db"
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

	// 订单状态
	Order_status int		`gorm:"default:0;"`
}


func init() {
	if db.DBSession.HasTable(&Tb_order{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_order{})
	}
}