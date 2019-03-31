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
	DBSession, err = gorm.Open("mysql", conf.Data.Mysql.Host + "/" +
		conf.Data.Mysql.Database + "?charset=utf8")
	if err != nil {
		log.Error("gorm init err:", err)
		os.Exit(-1)
	}

	DBSession.DB().SetMaxIdleConns(3)
	DBSession.DB().SetMaxOpenConns(10)
	DBSession.DB().SetConnMaxLifetime(time.Hour)

	RedisClient = &redis.Pool{
		MaxIdle : conf.Data.Rds.MaxIdle,
		MaxActive: conf.Data.Rds.MaxActive,
		IdleTimeout : 300*time.Second,
		Dial : func() (redis.Conn, error){
			c, err := redis.Dial("tcp", conf.Data.Rds.Host)
			if err != nil {
				return nil, err
			}

			c.Do("SELECT", conf.Data.Rds.DB)
			c.Do("AUTH", conf.Data.Rds.Auth)
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
