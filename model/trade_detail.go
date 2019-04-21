package model

import (
	"github.com/jinzhu/gorm"
	"github.com/BedeWong/iStock/db"
)


// 成交历史：在这里记录每笔成交的记录。
//  每个订单可能产生多个 成交记录
type Tb_trade_detail struct {
	gorm.Model
	Order_id	int			`gorm:"not null"`
	User_id 	int 		`gorm:"not null"`
	Stock_name  string  	`grom:"type:varchar(16); not null"`
	Stock_code  string  	`grom:"type:varchar(16); not null"`
	Stock_price float64 	`gorm:"type:decimal(12,2); not null"`
	Stock_count int			`grom:"default:0;"`
	Trade_type  int			`gorm:"default:0;"`

	Stamp_tax	float64		`gorm:"type:decimal(12,2);default:0.0"`	 // 印花税
	// 比赛id： 默认非比赛情况
	Contest_id int          `gorm:"default:0"`
}


func init() {
	if db.DBSession.HasTable(&Tb_trade_detail{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_trade_detail{})
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_trade_detail{})
	}
}