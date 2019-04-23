// 周期任务
//
// @author: BedeWong
// created_at: 2019年4月23日

package period

import (
	"github.com/robfig/cron"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/period/task"
)

// 使用介绍：
//    每隔5秒执行一次：*/5 * * * * ?
//    每隔1分钟执行一次：0 */1 * * * ?
//    每天23点执行一次：0 0 23 * * ?
//    每天凌晨1点执行一次：0 0 1 * * ?
//    每月1号凌晨1点执行一次：0 0 1 1 * ?
//    在26分、29分、33分执行一次：0 26,29,33 * * * ?
//    每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
//
func CornTask() {
	c := cron.New()

	spec := "0 0 1 * * ?"
	c.AddFunc(spec, task.UserRank)
	c.AddFunc(spec, task.UserAssets)

	spec = "0 20 8 * * ?"
	c.AddFunc(spec, task.ContestStatusCheckCanStart)

	spec = "0 50 15 * * ?"
	c.AddFunc(spec, task.ContestStatusCheckCanOver)

	c.Start()
	log.Info("period init ok.")
}
