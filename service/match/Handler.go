package match

import (
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/service/message"
	manager "github.com/BedeWong/iStock/service"
)

// 與外界通信的 channel
//var match_que chan interface{}
//
//func GetChan() (chan interface{}){
//	return match_que
//}

func Handler(task interface{}){
	switch item := task.(type) {
	default:
		log.Error("recv a item type:%T, val:%#v, can not hander it.", task, task)
	case message.MsgTickData:
		// 直接转发到 定序模块进行 匹配
		manager.Send2Senquence(item, 1)
	}
}

func Init() {
	manag := manager.GetInstance()
	match_que := manag.Match_que

	go func() {
		for {
			task := <-match_que
			log.Info("recv a new task: %T, %#v", task, task)

			Handler(task)
		}
	}()

	log.Info("Match Init ok.")
}