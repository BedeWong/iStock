package data_source

import (
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/service/message"
	manager "github.com/BedeWong/iStock/service"
)

// 與外界通信的 channel
//var source_que chan interface{}
//
//func GetChan() (chan interface{}){
//	return source_que
//}

//  股票列表：這個列表只存放當前定序隊列裏存在的股票代碼
//  以這個列表為依據，獲取實槃數據，推送到撮合引擎，進行訂單匹配

// 处理函数, 外部模块发送来的命处理
func HandlerCmd(workers SourceHandler, item interface{}) {
	switch obj := item.(type) {
	default:
		log.Error("recv a item type:%T, val:%#v, can not hander it.", item, item)
	case message.MsgSourceStockDeal:		// add stock to source
	if obj.Type == message.MsgSourceStockDealType_ADD {
		log.Info("recv add stock to source. %s", obj.Stock_code)
		workers.Append(obj.Stock_code)
	}else if obj.Type == message.MsgSourceStockDealType_Del {
		log.Info("recv del stock to source. %s", obj.Stock_code)
		workers.Remove(obj.Stock_code)
	}
	}
}

// Source worker 数据收集
// 收集的数据发送到 撮合 模块
func HanderTickData(ch <-chan message.MsgTickData) {
	for {
		data, ok := <- ch
		if ok == false {
			log.Error("数据源被无情的关闭了,怎么肥四?")
			break
		}

		err := manager.Send2Match(data, 1)
		if err != nil {
			log.Error("HanderTickData:Send2Match err:%s", err.Error())
		}
	}
}

func Init() {
	source_que := manager.GetInstance().Source_data_que

	// 数据源处理类
	var source_handler SourceHandler = NewSourceHandler()

	// 处理 数据源 worker 的tick数据
	go HanderTickData(source_handler.GetTickDataChan())

	go func() {
		for {
			select {
				//获取过来的任务通常是一个 命令
				case task := <-source_que:
					log.Info("recv a new task: %T, %#v", task, task)
					HandlerCmd(source_handler, task)
			}
		}
	}()

	log.Info("Source Init ok.")
}