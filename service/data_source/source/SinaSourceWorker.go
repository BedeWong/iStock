package source

import (
	"context"
	"github.com/BedeWong/iStock/model"
)

type SinaSourceWorker struct {
	BaseSourceWorker
	BaseUrl      string
}


// 创建对象
func NewSinaSourceWorker(fn_ context.CancelFunc, ctx_ context.Context, code_ string, url string) SinaSourceWorker{
	worker := SinaSourceWorker {
		BaseSourceWorker: BaseSourceWorker{
			cancel: fn_,
			ctx: ctx_,
			code: code_,
		},
		BaseUrl : url,
	}

	worker.specific = &worker
	return worker
}


func(this *SinaSourceWorker)FechOnce() (model.Tb_tick_data, error){


	return model.Tb_tick_data{}, nil
}
