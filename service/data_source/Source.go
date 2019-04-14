package data_source

import (
	"context"
	"github.com/BedeWong/iStock/service/data_source/source"
	"github.com/BedeWong/iStock/model"
	"github.com/gpmgo/gopm/modules/log"
)

// 逐筆成交記錄
//type TickData struct {
//	Tick_time 	string	`json:"tick_time"`
//	Tick_count	int		`json:"tick_count"`
//	Tick_price  float64		`json:"tick_price"`
//	Tick_Type 	string 		`json:"tick_type"`			// 交易类型： 买盘， 卖盘， 中性盘
//}

// 每支要获取tick数据的股票（代码） 对应的是一个协程， 这里存储的是 协程的context.cancel句柄
type Stocks map[string] source.ISourceWorker

type SourceHandler struct {
	stocks Stocks

	// 所有的 协程 往这里推送tick数据
	ch chan model.Tb_tick_data
}

// 暂时这样写好吧
func NewSourceHandler() SourceHandler {
	return SourceHandler{
		stocks: make(map[string] source.ISourceWorker),
		ch : make(chan model.Tb_tick_data),
	}
}

// 添加 一支股, 对应的起一个协程，用于获取 实时的逐笔交易数据。
func (this *SourceHandler)Append(code string) {
	log.Debug("添加一个股票代码：%s", code)
	_, ok := this.stocks[code]
	if ok == false{
		log.Debug("股票代码: %s 未启动tick 数据获取.", code)
		ctx, cancel := context.WithCancel(context.Background())
		work := source.NewBaseSourceWorker(cancel, ctx, code)
		this.stocks[code] = work

		// 启动协程， 抓取数据
		go work.FetchWork(this.GetTickDataChan())
	}
}

// 删除 一只股， 只有在订单队列（定序模块管理）中没有这只股的订单的时候就可以执行本段代码
// 同时删除 运行的协程
func (this *SourceHandler)Remove(code string) {
	worker, ok := this.stocks[code]
	if ok {
		delete(this.stocks, code)

		worker.Cancel()
	}
}

//返回 channel
func (this *SourceHandler)GetTickDataChan() (  chan model.Tb_tick_data) {
	return this.ch
}

