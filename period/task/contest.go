package task

import (
	"github.com/BedeWong/iStock/db"
	"github.com/gpmgo/gopm/modules/log"
	"fmt"
)


// 比赛状态设置：开始
func ContestStatusCheckCanStart(){
	type contestInfo struct {
		Id int
	}

	var contests []contestInfo
	// 查询今天该开始的比赛
	page := 1
	count := 0
	for {
		sql := fmt.Sprintf("select id from tb_contest where " +
			"c_status=0 and c_start_date=CURDATE() " +
			"limit %d, %d", (page-1)*100, 100)
		err := db.DBSession.Raw(sql).Scan(&contests).Count(&count).Error
		page++
		if count <= 0 {
			log.Info("period task ContestStatusCheckCanStart: " +
				"sql: %s, count: %d", sql, count)
			break
		}
		if err != nil {
			log.Error("period task ContestStatusCheckCanStart: "+
				"sql: %s, err: %v", sql, err)
			// 当前处理退出.
			break
		} else {
			// 有数据要处理.
			for idx, contest := range contests{
				sql := fmt.Sprintf("update tb_contest set c_status=1 " +
					"where id=%d", contest.Id)
				log.Debug("id: %d, staring. sql: %s", contest.Id, sql)
				err = db.DBSession.Exec(sql).Error
				if err != nil {
					log.Error("period task ContestStatusCheckCanStart: "+
						"idx: %d, sql: %s, err: %v", idx, sql, err)
					// 当前处理退出.
					continue
				}
			}
		}// end if
	}// end for -- 今天该开始的场次.
}

// 比赛状态设置：结束
func ContestStatusCheckCanOver(){
	type contestInfo struct {
		Id int
	}
	var contests []contestInfo

	// 更新今天该结束的场次
	page := 1
	count := 0
	for {
		sql := fmt.Sprintf("select id from tb_contest where " +
			"c_status=1 and  c_end_date=CURDATE() " +
			"limit %d, %d", (page-1)*100, 100)
		err := db.DBSession.Raw(sql).Scan(&contests).Count(&count).Error
		page++
		if count <= 0 {
			log.Info("period task ContestStatusCheckCanOver: " +
				"sql: %s, count: %d", sql, count)
			break
		}
		if err != nil {
			log.Error("period task ContestStatusCheckCanOver: "+
				"sql: %s, err: %v", sql, err)
			// 当前处理退出.
			break
		} else {
			// 有数据要处理.
			for idx, contest := range contests{
				sql := fmt.Sprintf("update tb_contest set c_status=2 " +
					"where id=%d", contest.Id)
				log.Debug("id: %d, over. sql: %s", contest.Id, sql)
				err = db.DBSession.Exec(sql).Error
				if err != nil {
					log.Error("period task ContestStatusCheckCanOver: "+
						"idx: %d, sql: %s, err: %v", idx, sql, err)
					// 当前处理退出.
					continue
				}
			}
		}// end if
	}
}