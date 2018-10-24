package model

import (
	"github.com/jinzhu/gorm"
	"github.com/BideWong/iStock/db"
)

// 用户账户表， 不存其他的业务数据， 关联 业务用户userid
// 本表只存储本系统中用到的数据
type Tb_user struct {
	gorm.Model
	User_id 		int 		`gorm:"not null; unique_index"`
	// 用户的可用资金
	User_money		float64     `gorm:"type:decimal(12,3)"`
	//  用户的总资产：包括当前的市值
	User_all_capital	float64 `gorm:"type:decimal(12,3)"`
}

func init() {
	if db.DBSession.HasTable(&Tb_user{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_user{})
		//fmt.Println("创建表：Tb_user ok")
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_user{})
	}
}