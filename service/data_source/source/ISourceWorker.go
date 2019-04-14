package source

import (
	"github.com/BedeWong/iStock/model"
)

type ISourceWorker interface {
	// 协程 执行的主方法
	// 这个方法， 根据业务重写
	FetchWork(chan<- model.Tb_tick_data) (error)

	//
	FechOnce() (model.Tb_tick_data, error)

	// 控制协程的结束
	Cancel()
}
