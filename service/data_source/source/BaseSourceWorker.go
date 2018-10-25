package source

import (
	"context"
	"time"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BideWong/iStock/service/message"
)

type BaseSourceWorker struct {
	// cancel 句柄
	cancel context.CancelFunc

	// ctx对象
	ctx context.Context

	// 股票代码  sh000001
	code string
}


// New 对象
func  NewBaseSourceWorker(fn_ context.CancelFunc, ctx_ context.Context, code_ string) *BaseSourceWorker {
	return &BaseSourceWorker{
		cancel:fn_,
		ctx:ctx_,
		code:code_,
	}
}

// 协程的执行主方法体，
func (this *BaseSourceWorker) FetchWork(ch chan<- message.MsgTickData) (error) {
	for {
		tick := time.Tick(time.Second)
		select {
		case <- this.ctx.Done():
			// this coroutine over
			log.Info("[%s] worker over.", this.code)
			return nil

		case <- tick:
			dat, err := this.FechOnce()
			if err != nil {
				log.Error("FechOnce err:%v", err)
			}

			//  发送数据到
			ch <- dat
		}
	}
	return nil
}

// 获取一条数据
func (this *BaseSourceWorker)FechOnce() (message.MsgTickData, error){
	tm := time.Now()
	return message.MsgTickData{
		Tick_code : this.code,
		Tick_time : tm.String(),
		Tick_count : 100,
		Tick_price : 12.66,
		Tick_Type : "UP",
	}, nil
}

//  控制协程结束执行
func (this *BaseSourceWorker)Cancel() {
	this.cancel()
}