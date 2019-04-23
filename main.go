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
	"github.com/BedeWong/iStock/period"
)

func main(){
	// 開啟調試日誌.
	log.Verbose = true
	//  關閉數據庫鏈接
	defer db.CloseDB()

	//  這裏統一初始化 各個模塊
	data_source.Init()
	match.Init()
	sequence.Init()
	clearing.Init()

	// 周期任务
	period.CornTask()

	log.Info("main init ok.")
	select{}
}
