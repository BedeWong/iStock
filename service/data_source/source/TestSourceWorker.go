// Test source worker. 調試使用類.
//
// 該類用於測試，FetchOnce從數據庫獲取tick數據
// 外部測試條件：往Model.Tick_data表中寫入測試數據即可
//
// @author:bedewong
// @create at:
// @update at: 2019年4月20日
// @change log:

package source

import (
	"context"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/model"
	"github.com/BedeWong/iStock/db"
)

type TestSourceWorker struct {
	// 組合BaseWoker
	BaseSourceWorker
}


// New TestSourceWorker 对象
func  NewTestSourceWorker(fn_ context.CancelFunc, ctx_ context.Context, code_ string) TestSourceWorker {
	return TestSourceWorker{
		BaseSourceWorker: NewBaseSourceWorker(fn_, ctx_, code_),
	}
}


// 實現數據獲取方法.
//
// 本方法從數據獲取測試數據。
func (this *TestSourceWorker)FechOnce() ([]model.Tb_tick_data, error){
	var tick_datas []model.Tb_tick_data // datas
	var cnt = 0

	// 數據獲取
	db.DBSession.Where("tick_code = ?", this.code).
		Order("tick_price, tick_time").  // 價格，時間升序
		Find(&tick_datas).Count(&cnt)
	log.Debug("BaseSource:FechOnce: fetch data: %v", tick_datas)

	// 数据时效性：本次获取到的数据不会在下次的查询中出现.
	//
	// 获取所有查询出来的数据的id
	ids := make([]uint, len(tick_datas))
	for _, item := range tick_datas {
		ids = append(ids, item.ID)
	}
	// 删除数据
	db.DBSession.Where("id in (?)", ids).Delete(&model.Tb_tick_data{})

	return tick_datas, nil
}