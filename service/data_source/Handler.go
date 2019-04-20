package data_source

import (
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/service/message"
	manager "github.com/BedeWong/iStock/service"
	"github.com/BedeWong/iStock/model"
	"context"
	"github.com/BedeWong/iStock/service/data_source/source"
)


// 外部模块发送来的命处理方法
//
//
func HandlerCmd(workers SourceHandler, item interface{}) {
	switch obj := item.(type) {
	default:
		log.Error("recv a item type:%T, val:%#v, can not hander it.", item, item)
	case message.MsgSourceStockDeal:
		log.Debug("source:HandlerCmd 收到一个股票代码 stock: %#v", obj.Stock_code)
		if obj.Type == message.MsgSourceStockDealType_ADD {
			// add stock to source
			log.Info("recv add stock to source. %s", obj.Stock_code)
			workers.Append(obj.Stock_code)
		}else if obj.Type == message.MsgSourceStockDealType_Del {
			// del stock from source
			log.Info("recv del stock to source. %s", obj.Stock_code)
			workers.Remove(obj.Stock_code)
		}
	}
}


// 本協程處理Source worker 数据收集.
//
// worker會往ch裡面發送數據，獲取源數據發送到撮合模塊.
// @param ch:
func HanderTickData(ch <-chan []model.Tb_tick_data) {
	for {
		// 等待數據
		data, ok := <- ch
		if ok == false {
			log.Error("数据源被无情的关闭了,怎么肥四?")
			break
		}
		log.Debug("data_source:Handler:HandlerTickDara: recv a tick data: data: %#v", data)

		err := manager.Send2Match(data, 1)
		if err != nil {
			log.Error("HanderTickData:Send2Match err:%s", err.Error())
		}
	}
}


func Init() {
	// 與本模塊通信的channel
	task_chan := manager.GetInstance().Source_data_que

	// 数据源处理类
	var source_handler = NewSourceHandler()

	// 处理 数据源 worker 的tick数据
	go HanderTickData(source_handler.GetTickDataChan())

	// 本協程處理其他模塊發送過來的命令
	go func() {
		for {
			select {
				//获取过来的任务通常是一个 命令
				case task := <-task_chan:
					log.Info("data soure recv a new stock: %T, %#v", task, task)
					HandlerCmd(source_handler, task)
			}
		}
	}()

	log.Info("Source Init ok.")
}


// 逐筆成交記錄數據結構描述
//type Tick_data struct {
//	Tick_time 	time.Time	`json:"tick_time"`
//	Tick_count	int		`json:"tick_count"`
//	Tick_price  float64		`json:"tick_price"`
//	Tick_Type 	string 		`json:"tick_type"`			// 交易类型： 买盘， 卖盘， 中性盘
//}

// 每支要获取tick数据的股票（代码） 对应的是一个SourceWorker
type Stocks map[string] source.ISourceWorker

type SourceHandler struct {
	stocks Stocks

	// 接收數據
	// worker協程往這裡推送tick數據.
	// SourceHandler從這裡讀取數據做處理.
	ch chan []model.Tb_tick_data
}


// 創建source handler處理對象
func NewSourceHandler() SourceHandler {
	return SourceHandler{
		stocks: make(map[string] source.ISourceWorker),
		ch : make(chan []model.Tb_tick_data),
	}
}


// 添加一支股, 对应的起一个协程，用于获取 实时的逐笔交易数据。
func (this *SourceHandler)Append(code string) {
	log.Debug("添加一个股票代码：%s", code)

	_, ok := this.stocks[code]
	// 判斷該股票代碼的處理worker是否還未創建.
	if ok == false{
		log.Debug("股票代码: %s 未启动tick 数据获取.", code)

		ctx, cancel := context.WithCancel(context.Background())
		work := source.NewTestSourceWorker(cancel, ctx, code)
		this.stocks[code] = &work

		// 启动协程， 抓取数据
		go work.FetchWork(this.GetTickDataChan())
	}
}


// 刪除一支股票代碼.
//
// 刪除對應的股票數據獲取，關閉數據獲取協程.
// @parma: code: stock code
func (this *SourceHandler)Remove(code string) {
	worker, ok := this.stocks[code]
	if ok {
		delete(this.stocks, code)

		worker.Cancel()
	}
}


//返回 channel
func (this *SourceHandler)GetTickDataChan() (chan []model.Tb_tick_data) {
	return this.ch
}
