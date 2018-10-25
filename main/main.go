package main

import (
	_ "github.com/BideWong/iStock/conf"
	_ "github.com/BideWong/iStock/match_rpc"
	"github.com/BideWong/iStock/db"

	_ "github.com/BideWong/iStock/model"

	// 各个子模块
	"github.com/BideWong/iStock/service/data_source"
	"github.com/BideWong/iStock/service/match"
	"github.com/BideWong/iStock/service/sequence"
	"github.com/BideWong/iStock/service/clearing"
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
