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

type Config struct {
	Rpc ConfigRpc			// node rpc
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

func Init() {
	LoadConfig("./config.json", Data)
}

