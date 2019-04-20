// 基於ISourceWorker 實現的數據源獲取類.
//
// 基本實現了基礎的功能，該類不能直接使用，
// 具體的類需要實現FechOnce方法用於獲取數據，本類相當於一個模板類
//
// @author：bedewong
// @create at:
// @update at: 2019年4月20日
// @change log:
//		@author:
package source

import (
"context"
"time"
"github.com/gpmgo/gopm/modules/log"
"github.com/BedeWong/iStock/model"
	"github.com/pkg/errors"
)

type BaseSourceWorker struct {
	specific ISourceWorker
	// cancel 句柄
	cancel context.CancelFunc
	// ctx对象
	ctx context.Context
	// 股票代码  same as: sh000001
	code string
}


// New BaseSourceWorker 对象
//
// @param: 見 BaseSourceWorker 結構的聲明
// @return:
func NewBaseSourceWorker(fn_ context.CancelFunc, ctx_ context.Context, code_ string) BaseSourceWorker {
	return BaseSourceWorker{
		cancel:fn_,
		ctx:ctx_,  // 協程控制context
		code:code_,  // stock code
	}
}

// 協程工作，獲取tick數據.
//
// 循環每秒執行一次數據更新， 協程結束條件：外部控制
// @param ch: 數據輸出通道.
// @return:
func (this *BaseSourceWorker) FetchWork(ch chan<- []model.Tb_tick_data) (error) {
	log.Debug("BaseSourceWorker:FetchWork starting. code: [%s]", this.code)
	for {
		tick := time.Tick(time.Second * 5)
		select {
		case <- this.ctx.Done():
			// this coroutine over
			log.Info("[%s] worker over.", this.code)
			return nil

		case <- tick:
			tick_data, err := this.FechOnce()
			log.Debug("获取到一个tick数据包: %v", tick_data)
			if err != nil {
				log.Error("FechOnce err:%v", err)
			}

			//  发送数据到
			ch <- tick_data
		}
	}
	return nil
}

// 获取一条数据接口
//
// 留給具體類實現.
func (this *BaseSourceWorker)FechOnce() ([]model.Tb_tick_data, error){
	if this.specific != nil {
		return this.specific.FechOnce()
	}

	// error
	panic(errors.New("BaseSourceWorker::FechOnce NotImplement."))
}

// 控制协程结束执行
//
func (this *BaseSourceWorker)Cancel() {
	this.cancel()
}
