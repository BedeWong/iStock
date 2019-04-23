package task

import (
	"github.com/BedeWong/iStock/db"
	"github.com/gpmgo/gopm/modules/log"
	"fmt"
)


// T+1
func tAdd1Process(){
	sql := "update tb_user_positions set " +
		"stock_count_can_sale=stock_count"
	db.DBSession.Exec(sql)
}


// T+1 比赛
func tAdd1ConteestProcess(){
	sql := "update tb_user_contest_positions set " +
		"stock_count_can_sale=stock_count"
	db.DBSession.Exec(sql)
}


// 资产走势处理
func assetsTrend()  {
	sql := "insert into tb_user_assets_trends(user_id, user_mv, `day`) " +
		"select user_id, user_mv, CURDATE() from tb_user_assets "
	db.DBSession.Exec(sql)
}


// 计算用户当前市值
// bedewong (NOTE): GORM 使用事务有坑, 暂时不使用事务
// todo...
func marketValueProcess() {
	sql_init := "update tb_user_assets set user_mv=user_money"
	//tx := db.DBSession.Begin()

	// init 操作
	err := db.DBSession.Exec(sql_init).Error
	if err != nil {
		log.Error("peroid task marketValueProcess update init err: %v", err)
		return
	}

	// 把用户持仓的数据市值计算到用户的市值中去
	//sql execute
	sql_user_pos := fmt.Sprintf("select user_id, sum(stock_count*stock_price) as mv " +
		" from tb_user_positions GROUP BY user_id")
	rows, err := db.DBSession.Raw(sql_user_pos).Rows()
	if err != nil {
		log.Error("peroid task marketValueProcess select sql_user_pos err: %v", err)
		log.Error(sql_user_pos)
		return
	}
	defer rows.Close()

	idx := 0
	for rows.Next() {
		var user_id int
		var mv float64
		rows.Scan(&user_id, &mv)

		idx++
		sql_user_mv := fmt.Sprintf("update tb_user_assets " +
			"set user_mv=user_mv+%f " +
			"where user_id=%d", mv, user_id)
		log.Debug("update set tb_user_assets: sql: %s",
			sql_user_mv)
		err := db.DBSession.Exec(sql_user_mv).Error
		if err != nil {
			log.Error("peroid task marketValueProcess update sql_user_mv " +
				"idx:%d, err: %v", idx, err)
			log.Error(sql_user_mv)
			continue
		}
	}// end for 为有股票资产的用户累加资产.

	// 提交事务
	//tx.Commit()
}


// 计算用户在比赛中的总市值数据
// todo...
func contestMarketValueProcess() {
	//sql_init := "update tb_user_assets set user_mv=user_money"
	//tx := db.DBSession.Begin()
	//
	//// init 操作
	//err := tx.Exec(sql_init).Error
	//if err != nil {
	//	log.Error("peroid task marketValueProcess update init err: %v", err)
	//	return
	//}
	//
	//page := 0
	//for {
	//	type user_pos_mv struct {
	//		user_id int
	//		mv string
	//	}
	//	var users []user_pos_mv
	//	var cnt int
	//	// 把用户持仓的数据市值计算到用户的市值中去
	//	page++
	//	sql_user_pos := fmt.Sprintf("select user_id, sum(stock_count*stock_price) as 'mv'" +
	//		" from tb_user_positions GROUP BY user_id limit %d, %d", (page-1)*1000, 1000)
	//	err := db.DBSession.Exec(sql_user_pos).Scan(&users).Count(&cnt).Error
	//	if err != nil {
	//		log.Error("peroid task marketValueProcess update sql_user_pos err: %v", err)
	//		log.Error(sql_user_pos)
	//		break
	//	}
	//
	//	// 没有数据可清算了
	//	if cnt == 0 {
	//		break
	//	}
	//
	//	for idx, user := range users {
	//		// 累加上用户的股票资产
	//		sql_user_mv := fmt.Sprintf("update tb_user_assets " +
	//			"set user_mv=user_mv+%s " +
	//			"where user_id=%d", user.mv, user.user_id)
	//		err := db.DBSession.Exec(sql_user_pos).Scan(&users).Count(&cnt).Error
	//		if err != nil {
	//			log.Error("peroid task marketValueProcess update sql_user_mv " +
	//				"idx:%d, err: %v", idx, err)
	//			log.Error(sql_user_mv)
	//			continue
	//		}
	//	}// end for 为有股票资产的用户累加资产.
	//}
	//// 提交事务
	//tx.Commit()
}


// 用户资产相关处理：
//
// T+1 股处理
// 用户资产走势
func UserAssets() {
	// T+1 处理
	tAdd1Process()
	tAdd1ConteestProcess()

	// user_mv
	marketValueProcess()
	// user contest mv
	contestMarketValueProcess()

	// 资产走势
	assetsTrend()
	// stock trend
	//todo...
}


// 排名处理
//
// 比赛排名
// 普通排名
func UserRank(){
	userRank()
	userContestRank()
}


func userContestRank(){
	// 获取正在比赛的场次
	//todo...
	// init..
	//sql := "update tb_user set u_rank=0"
	//err := db.DBSession.Exec(sql).Error
	//if err != nil {
	//	log.Error("peroid task userContestRank update init err: %v", err)
	//	return
	//}
	//
	//type user struct {
	//	user_id int
	//}
}


func userRank(){
	// init..
	sql := "update tb_user set u_rank=0"
	err := db.DBSession.Exec(sql).Error
	if err != nil {
		log.Error("peroid task UserRank update init err: %v", err)
		return
	}

	type user struct {
		User_id int
	}
	var users []user
	i := 1
	for i<=10 {
		err := db.DBSession.Raw(
			"select user_id " +
				"from tb_user_assets " +
				"order by user_mv " +
				"limit ?, ?", (i-1)*1000, 1000).Scan(&users).Error
		if err != nil {
			log.Error("peroid task UserRank i:%d, err:%v", i, err)
			break
		}
		log.Debug("peroid task UserRank users:%#v", users)

		// 写数据库
		for rank, user := range users {
			rank += 1   // rank 从0开始的
			sql := fmt.Sprintf("update tb_user set u_rank=%d where id=%d",
				(i-1)*1000+rank, user.User_id)
			log.Debug("peroid task UserRank update user rank: %d, user_id: %d",
				(i-1)*1000+rank, user.User_id)
			err := db.DBSession.Exec(sql).Error
			if err != nil {
				log.Error("peroid task UserRank update err: " +
					"user_id: %d, rank: %d, err:%v",
					user.User_id, (i-1)*1000+rank, err)
				continue
			}
		}
		//db.DBSession.Commit()
		i++
	}
}