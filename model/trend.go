package model

import (
	"github.com/jinzhu/gorm"
	"github.com/BedeWong/iStock/db"
	"time"
)

// 用户资产走势图
//
type Tb_user_assets_trend struct {
	gorm.Model
	Day  time.Time     		   `gorm:"type:date; not null"`
	User_id 		int 		`gorm:"not null;"`
	//  用户的当前 市值
	User_mv	float64 			`gorm:"type:decimal(12,3);default 0"`
}


// 股价日走势图
//
type Tb_stock_trend struct {
	gorm.Model
	//  用户的当前 市值
	Day  time.Time     `gorm:"type:date; not null"`
	Price	float64 	`gorm:"type:decimal(12,3);default 0"`
	Stock_code string `gorm:"type:varchar(16);"`
}


func init() {
	if db.DBSession.HasTable(&Tb_user_assets_trend{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_user_assets_trend{})
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_user_assets_trend{})
	}

	//
	if db.DBSession.HasTable(&Tb_stock_trend{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_stock_trend{})
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_stock_trend{})
	}
}