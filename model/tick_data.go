package model

import (
	"time"
	"github.com/BedeWong/iStock/db"
	"github.com/jinzhu/gorm"
)

type Tb_tick_data struct {
	gorm.Model
	Tick_time 	time.Time	`json:"tick_time"`
	Tick_code   string  	`json:"tick_code"`
	Tick_count	int			`json:"tick_count"`
	Tick_price  float64		`json:"tick_price"`
	Tick_Type 	string 		`json:"tick_type"`			// 交易类型： 买盘， 卖盘， 中性盘
}


func init(){
	if db.DBSession.HasTable(&Tb_tick_data{}) == false {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_tick_data{})
	}else {
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Tb_tick_data{})
	}
}
