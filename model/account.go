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
	User_money		float64     `gorm:"type:decimal(12,3)"`
	User_all_capital	float64 `gorm:"type:decimal(12,3)"`
}


type User struct {

}

func (this *User) AddUser(uid int, money float64, capital float64) error {
	return nil
}


func init() {
	if db.DBSession.HasTable(&Tb_user{}) == false {
		// will append "ENGINE=InnoDB" to the SQL statement when creating table `users`
		db.DBSession.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Tb_user{})
		//fmt.Println("创建表：Tb_user ok")
	}else {
		//fmt.Println("表：Tb_user ok已存在。")
	}
}