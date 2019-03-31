package source

import (
	"github.com/BedeWong/iStock/service/message"
)

type ISourceWorker interface {
	// 协程 执行的主方法
	// 这个方法， 根据业务重写
	FetchWork(chan<- message.MsgTickData) (error)

	//
	FechOnce() (message.MsgTickData, error)

	// 控制协程的结束
	Cancel()
}
