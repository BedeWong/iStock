package model

import (
	"github.com/jinzhu/gorm"
	"github.com/BedeWong/iStock/db"
)


// 用户的 持股表
// 每个用户的每支股票对一个一条记录， 买卖时修改记录
type Tb_user_position struct {
	gorm.Model
	User_id 		int 		`gorm:"not null"`
	Stock_name 		string 		`grom:"type:varchar(16); not null"`
	Stock_code		string  	`grom:"type:varchar(16); not null"`
	// 持有？股
	Stock_count		int			`grom:"default:0;"`
	// 持倉價
	Stock_price 	float64 	`gorm:"type:decimal(12,2); not null"`

	// T+1 本股可卖部分， 当前交易日卖的股票本交易日不能出售
	Stock_count_can_sale int	`grom:"default 0"`
}

// 用户的 持股表
// 每个用户的每支股票对一个一条记录， 买卖时修改记录
type Tb_user_contest_position struct {
	gorm.Model
	User_id 		int 		`gorm:"not null"`
	Stock_name 		string 		`grom:"type:varchar(16); not null"`
	Stock_code		string  	`grom:"type:varchar(16); not null"`
	// 持有？股
	Stock_count		int			`grom:"default:0;"`
	// 持倉價
	Stock_price 	float64 	`gorm:"type:decimal(12,2); not null"`

	// T+1 本股可卖部分， 当前交易日卖的股票本交易日不能出售
	Stock_count_can_sale int	`grom:"default 0"`
	// 比赛id： 默认非比赛情况
	Contest_id int          `gorm:"default:0"`
}


func init() {
	if db.DBSession.HasTable(&Tb_user_position{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_user_position{})
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_user_position{})
	}
}