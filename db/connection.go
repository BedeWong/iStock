package db

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"github.com/gpmgo/gopm/modules/log"
	"os"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/BedeWong/iStock/conf"
)

var (
	DBSession *gorm.DB
	RedisClient *redis.Pool
)

func init() {
	var err error
	DBSession, err = gorm.Open("mysql", conf.GetConfig().Mysql.Host + "/" +
		conf.GetConfig().Mysql.Database + "?charset=utf8&parseTime=true")
	if err != nil {
		log.Error("gorm init err:", err)
		os.Exit(-1)
	}

	DBSession.DB().SetMaxIdleConns(3)
	DBSession.DB().SetMaxOpenConns(10)
	DBSession.DB().SetConnMaxLifetime(time.Hour)

	RedisClient = &redis.Pool{
		MaxIdle : conf.GetConfig().Rds.MaxIdle,
		MaxActive: conf.GetConfig().Rds.MaxActive,
		IdleTimeout : 300*time.Second,
		Dial : func() (redis.Conn, error){
			c, err := redis.Dial("tcp", conf.GetConfig().Rds.Host)
			if err != nil {
				return nil, err
			}

			c.Do("SELECT", conf.GetConfig().Rds.DB)
			c.Do("AUTH", conf.GetConfig().Rds.Auth)
			return c, nil
		},
	}
}

func CloseDB(){
	err := DBSession.Close()
	if err != nil {
		log.Error("close BD err:", err)
		return
	}
}
