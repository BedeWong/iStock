package main

import (
	_ "github.com/BideWong/iStock/conf"
	_ "github.com/BideWong/iStock/match_rpc"
	"github.com/BideWong/iStock/db"

	_ "github.com/BideWong/iStock/model"
)

func main(){
	// 關閉數據庫鏈接
	defer db.CloseDB()

	select{}
}
