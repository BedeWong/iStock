package main

import (
	_ "github.com/BedeWong/iStock/conf"
	_ "github.com/BedeWong/iStock/match_rpc"
	"github.com/BedeWong/iStock/db"

	_ "github.com/BedeWong/iStock/model"

	// 各个子模块
	"github.com/BedeWong/iStock/service/data_source"
	"github.com/BedeWong/iStock/service/match"
	"github.com/BedeWong/iStock/service/sequence"
	"github.com/BedeWong/iStock/service/clearing"
	"github.com/gpmgo/gopm/modules/log"
)

func main(){
	// 關閉數據庫鏈接
	defer db.CloseDB()

	// 這裏統一初始化 各個模塊
	data_source.Init()
	match.Init()
	sequence.Init()
	clearing.Init()

	log.Info("main init ok.")
	select{}
}
