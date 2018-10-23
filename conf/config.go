package conf

import (
	"encoding/json"
	"io/ioutil"
	"github.com/gpmgo/gopm/modules/log"
	"os"
)

type ConfigRpc struct {
	Addr string      		// listen at
	Pattern string   		// http handler pattern
	// ...
}

type ConfigRedis struct {
	Host string
	DB	int
	MaxIdle int
	MaxActive int
	Auth string
}

type ConfigMysql struct {
	Host string
	Database string
}

type Config struct {
	Rpc ConfigRpc			// node rpc
	Rds ConfigRedis
	Mysql ConfigMysql
}

var Data Config

// 加载 配置文件，只加载一次，全局单例
func LoadConfig(filename string, conf interface{}) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("read config err:", err)
		os.Exit(-1)
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		log.Error("Unmarshal config err:", err)
		os.Exit(-1)
	}
}

//  获取配置对象
func GetConfig() Config {
	return Data
}

func init() {
	LoadConfig("H:/mygo/src/github.com/BideWong/iStock/conf/config.json", &Data)
	//LoadConfig("./config.json", &Data)
}

