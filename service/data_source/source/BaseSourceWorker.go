package source

import (
	"context"
	"time"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/model"
	"github.com/BedeWong/iStock/db"
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
func (this *BaseSourceWorker) FetchWork(ch chan<- model.Tb_tick_data) (error) {
	log.Debug("BaseSourceWorker:FetchWork starting. code: [%s]", this.code)
	for {
		tick := time.Tick(time.Second)
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

// 获取一条数据
func (this *BaseSourceWorker)FechOnce() (model.Tb_tick_data, error){
	//tm := time.Now()
	//return message.MsgTickData{
	//	Tick_code : this.code,
	//	Tick_time : tm.String(),
	//	Tick_count : 100,
	//	Tick_price : 12.66,
	//	Tick_Type : "UP",
	//}, nil

	var tick_data model.Tb_tick_data
	db.DBSession.Where("tick_code = ?", this.code).Scan(&tick_data)

	log.Debug("BaseSource:FechOnce: fetch data: %v", tick_data)

	return tick_data, nil
}

//  控制协程结束执行
func (this *BaseSourceWorker)Cancel() {
	this.cancel()
}