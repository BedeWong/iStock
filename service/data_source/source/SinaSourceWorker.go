package source

import (
	"time"
	"context"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/model"
)

type SinaSourceWorker struct {
	BaseSourceWorker
	BaseUrl      string
}

// 创建对象
func NewSinaSourceHandler(fn_ context.CancelFunc, ctx_ context.Context, code_ string, url string) *SinaSourceWorker{
	return &SinaSourceWorker {
		BaseSourceWorker : *NewBaseSourceWorker(fn_, ctx_, code_),
		BaseUrl : url,
	}
}

func(this *SinaSourceWorker) FetchWork(ch chan<- model.Tb_tick_data) (error){
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


func(this *SinaSourceWorker)FechOnce() (model.Tb_tick_data, error){


	return model.Tb_tick_data{}, nil
}
