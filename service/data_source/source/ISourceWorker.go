package source

import "github.com/BedeWong/iStock/model"

/**
* 數據源woker接口類
*/
type ISourceWorker interface {
	// 协程 执行的主方法
	FetchWork(chan<- []model.Tb_tick_data) (error)

	FechOnce() ([]model.Tb_tick_data, error)

	// 控制协程的结束
	Cancel()
}
