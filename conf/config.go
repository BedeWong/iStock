package conf

import (
	"io/ioutil"
	"github.com/gpmgo/gopm/modules/log"
	"os"
	"encoding/json"
)

/***
	rpc 配置
 */
type ConfigRpc struct {
	Addr string      		// listen at
	Pattern string   		// http handler pattern
	// ...
}

/****
	reids配置
 */
type ConfigRedis struct {
	Host string
	DB	int
	MaxIdle int
	MaxActive int
	Auth string
}

/****
	数据库配置
 */

type ConfigMysql struct {
	Host string
	Database string
}

/***
	交易规则相关配置
*/
type ConfigTrade struct{
	TransferFeeSZ		float64
	TransferFeeSH		float64
	StampTax		float64
	Brokerage		float64
}

type Config struct {
	Rpc ConfigRpc			// node rpc
	Rds ConfigRedis
	Mysql ConfigMysql
	Trade ConfigTrade
}

var confData Config

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
	return confData
}

func init() {
	//LoadConfig("H:/mygo/src/github.com/BedeWong/iStock/conf/config.json", &confData)
	LoadConfig("./conf/config.json", &confData)

	log.Info("config module init ok.")
}

