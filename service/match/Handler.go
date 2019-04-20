// 撮合模塊.
//
// 根據訂單tick_data進行匹配訂單.
// 由於是模擬的交易，本模塊沒有實質的作用，
// 訂單匹配部分代碼已經移動到定序模塊中處理了.
//
// @author: bedewong
// @create at:
// @update at: 2019年4月20日
// @ change log:

package match

import (
	"github.com/gpmgo/gopm/modules/log"
	manager "github.com/BedeWong/iStock/service"
	"github.com/BedeWong/iStock/model"
)


func Handler(task interface{}){
	switch item := task.(type) {
	default:
		log.Error("recv a item type:%T, val:%#v, can not hander it.", task, task)
	case []model.Tb_tick_data:
		log.Debug("recv a Tb_tick_data list, send to senquence module. data: %v",
			task)
		// 直接转发到 定序模块进行 匹配
		manager.Send2Senquence(item, 1)
	}
}

func Init() {
	task_chan := manager.GetInstance().Match_que

	go func() {
		for {
			task := <-task_chan
			log.Info("recv a new task: %T, %#v", task, task)

			Handler(task)
		}
	}()

	log.Info("Match Init ok.")
}