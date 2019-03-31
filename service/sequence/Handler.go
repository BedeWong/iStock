package sequence

import (
	"github.com/BedeWong/iStock/model"
	pq "github.com/BedeWong/iStock/utils/priority_queue"
	"github.com/gpmgo/gopm/modules/log"
	"sync"
	"github.com/pkg/errors"
	"encoding/json"
	"fmt"
	"github.com/BedeWong/iStock/db"
	"github.com/BedeWong/iStock/service/message"
	manager "github.com/BedeWong/iStock/service"
)

// 为每支股票维护两个队列，
// 入队：  添加了一个订单
// 出队：  订单匹配订单， 订单清算， 订单清算未全部完成的 继续入队
type Handler struct {

	sync.Mutex

	// 没支股票对应一个 优先队列[买，卖]
	PlateBuy map[string] pq.PriorityQueue
	PlateSale map[string] pq.NPriorityQueue

	//  登记 所有的在队列里等待的订单， 未成交的可以撤单， 成交的从此队列删除
	// key: order.ID
	// val： order recode
	orders map[uint] model.Tb_order_real
}

var OrderHandler *Handler

var sequence_que chan interface{}
////  返回这个 channel
//func GetChan() (chan interface{}) {
//	return sequence_que
//}

func Init() {
	// 初始化 chan
	sequence_que = manager.GetInstance().Sequence_que

	// 初始化 模块实例
	OrderHandler = &Handler{
		PlateBuy : make(map[string] pq.PriorityQueue),
		PlateSale : make(map[string] pq.NPriorityQueue),
		orders : make(map[uint]model.Tb_order_real),
	}

	 go func() {
	 	for {
	 		// 获取推送来的订单
	 		it := <- sequence_que

	 		log.Info("sequence recv a new task:%T, %#v", it, it)

			switch item := it.(type) {
			default:
				log.Error("recv a item type:%T, val:%#v, can not hander it.", it, it)
			case model.Tb_order_real:   		// 委托下单
				OrderHandler.AddOrder(item)
			case message.MsgRevokeOrder:		// 委托 撤单
				OrderHandler.DelOrder(item)
			case message.MsgTickData:			// 逐笔成交数据，
				OrderHandler.Match(item.Tick_code, item.Tick_price, item.Tick_count)
			}
		}
	 }()

	log.Info("Sequence Init ok.")
}

// 获取 实例
func GetInstance ()*Handler {
	return OrderHandler
}

// 添加訂單
func (this *Handler) AddOrder(order model.Tb_order_real) error{
	// 对操作加锁
	this.Lock()
	defer this.Unlock()

	id := order.ID
	if id == 0 {
		return errors.New("訂單ID錯誤:id=0")
	}

	old_order, ok := this.orders[id]
	if ok == true {
		res, _ := json.Marshal(old_order)
		log.Info("已存在旧的订单：%s", string(res))
	}

	// 添加[修改]订单
	this.orders[id] = order

	// 将订单添加到 买卖盘 队列
	item := pq.NewQueueNode(order, order.UpdatedAt.UnixNano(), order.Stock_price)
	if order.Trade_type == 0 {
		// 委托买单

		que, ok := this.PlateBuy[order.Stock_code];
		if ok == false {
			// 当前队列还未建立
			this.PlateBuy[order.Stock_code] = pq.NewPQ()
			que = this.PlateBuy[order.Stock_code]
		}
		que.Push(item)

		msg_tmp, err := json.Marshal(order)
		if err != nil {
			log.Error("json.Marshal(order) err:%v", err)
		}else {
			log.Info("添加order[%s]成功. len=%d", string(msg_tmp), len(que))
		}
	}else if order.Trade_type == 1 {
		// 委托卖单
		que, ok := this.PlateSale[order.Stock_code];
		if ok == false {
			// 当前队列还未建立
			this.PlateSale[order.Stock_code] = pq.NewNPQ()
			que = this.PlateSale[order.Stock_code];
		}
		que.Push(item)

		msg_tmp, err := json.Marshal(order)
		if err != nil {
			log.Error("json.Marshal(order) err:%v", err)
		}else {
			log.Info("添加order[%s]成功. len=%d", string(msg_tmp), len(que))
		}
	}

	return nil
}

// 撤單
//
func (this *Handler)DelOrder(order message.MsgRevokeOrder) error{
	// 对操作加锁
	this.Lock()
	defer this.Unlock()

	id := order.Order_id
	if id == 0 {
		return errors.New("訂單ID錯誤:id=0")
	}

	order_real, ok := this.orders[id]
	if ok == false {
		return errors.New(fmt.Sprintf("订单不存在（id=%d）", id))
	}


	// 从 orders 中删除
	// 从 match 队列中删除
	// 清算
	delete(this.orders, id)

	if order_real.Trade_type == model.TRADE_TYPE_SALE {
		que, ok := this.PlateSale[order_real.Stock_code]
		if ok == false {
			err := errors.New(fmt.Sprintf("%s stock queue not exist", order_real.Stock_code))
			log.Error(err.Error())
			return err
		}

		err := que.Remove(func(val interface{}) bool {
			item, ok := val.(model.Tb_order_real)
			if ok == false {
				e := errors.New(fmt.Sprintf("item is not a Tb_order_real object. %T", item))
				log.Error(e.Error())
				return false
			}

			if item.ID == id {
				log.Info("finded this item: %v", item)
				return true
			}

			return false
		})
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}else if order_real.Trade_type == model.TRADE_TYPE_BUY {
		que, ok := this.PlateBuy[order_real.Stock_code]
		if ok == false {
			err := errors.New(fmt.Sprintf("%s stock queue not exist", order_real.Stock_code))
			log.Error(err.Error())
			return err
		}

		err := que.Remove(func(val interface{}) bool {
			item, ok := val.(model.Tb_order_real)
			if ok == false {
				e := errors.New(fmt.Sprintf("item is not a Tb_order_real object. %T", item))
				log.Error(e.Error())
				return false
			}

			if item.ID == id {
				log.Info("finded this item: %v", item)
				return true
			}

			return false
		})
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	// 修改 状态   [ 撤单]
	order_real.Order_status = model.ORDER_STATUS_REVOKE

	// 清算
	manager.Send2Clearing(order_real, 2)

	// 保存订单的 更新
	db.DBSession.Save(&order_real)

	return nil
}

//// 本方法用于鉴定 订单是否是已经
//func (this *Handler)checkOrderIsRevoke(order_real model.Tb_order_real) bool{
//	return order_real.Order_status == model.ORDER_STATUS_REVOKE
//}
//
//// 撤单的 订单处理
//func (this *Handler)RevokeOrderHandler() error{
//
//	return nil
//}

// 处理 卖队列 订单
//
func (this *Handler)matchHandlerSaleQue(que pq.NPriorityQueue, price float64, count int) error {
	for {
		if len(que) == 0 {
			break
		}

		// 取出order item
		ele := que.Pop()
		it, b := ele.(*pq.Item)
		if b == false {
			log.Error("que.Pop() item not a *pq.Item object: %T, %v", ele, ele)
			continue
		}

		order_real, b := it.Value().(model.Tb_order_real)
		if b {
			if price >= order_real.Stock_price {
				// 成交价 >= 委托价： 模拟盘按 委托价 成交
				trade_detail := model.Tb_trade_detail{
					Order_id : order_real.Order_id,
					User_id : order_real.User_id,
					Stock_name : order_real.Stock_name,
					Stock_code : order_real.Stock_code,
					Stock_price : order_real.Stock_price,   // 委托价 成交
					Trade_type : order_real.Trade_type,
				}

				if count <= order_real.Stock_count {
					// 比 委托的单量小， 订单部分成交， 未完成的部分继续入队
					trade_detail.Stock_count = count
					order_real.Stock_count -= count
					count = 0
				}else {
					// 比委托的单量大， 全部成交， 修改订单记录状态为完成。
					trade_detail.Stock_count = order_real.Stock_count
					order_real.Order_status = model.ORDER_STATUS_FINISH
					// 剩餘的股 继续下个档位交易
					count -= order_real.Stock_count
					order_real.Stock_count = 0
				}

				// 保存订单
				db.DBSession.Save(&order_real)

				// 将交易明细 发送到 清算系统（clearing）
				manager.Send2Clearing(trade_detail, 2)

				//  本档 为全部完成， 继续入队。
				//  本档交易都未能全部完成， 不会再有
				if order_real.Stock_count > 0 {
					this.orders[order_real.ID] = order_real  //  更新 orders
					sequence_que <- order_real

					msg_tmp, err := json.Marshal(order_real)
					if err != nil {
						log.Error("json.Marshal(order_real) err:%v", err)
					}else {
						log.Info("sequence_que <- order_real: %s", string(msg_tmp))
					}
				}else {
					// 本次的订单需要从 orders 中移除
					delete(this.orders, order_real.ID)
				}

				// 外盘的本笔交易 已经完成
				if count == 0 {
					break
				}
			} else if price < order_real.Stock_price {
				// 成交价 比委托价低， 不成交
				log.Info("sale que当前成交价：%f, 盘口价：%f", price, order_real.Stock_price)

				// 后面也不会有 符合条件的订单了， 直接退出
				break
			}
		} else {
			log.Error("sale que item not handler: %T, %v", it.Value(), it.Value())
			return errors.New(fmt.Sprintf("sale que item not handler: %T, %v", it.Value(), it.Value()))
		}
	}
	return nil
}

// 处理 买盘队列
//
func (this *Handler)matchHandlerBuyQue(que pq.PriorityQueue, price float64, count int) error {
	for {
		if len(que) == 0 {
			break
		}

		// 取出order item
		ele := que.Pop()
		it, b := ele.(*pq.Item)
		if b == false {
			log.Error("que.Pop() item not a *pq.Item object: %T, %v", ele, ele)
			continue
		}

		order_real, b := it.Value().(model.Tb_order_real)
		if b {
			if price <= order_real.Stock_price {
				//  实盘成交价 要小于等于 模拟盘的 委托价
				trade_detail := model.Tb_trade_detail{
					Order_id : order_real.Order_id,
					User_id : order_real.User_id,
					Stock_name : order_real.Stock_name,
					Stock_code : order_real.Stock_code,
					Stock_price : order_real.Stock_price,   // 委托价 成交
					Trade_type : order_real.Trade_type,
				}

				if count <= order_real.Stock_count {
					// 比 委托的单量小， 订单部分成交， 未完成的部分继续入队
					trade_detail.Stock_count = count
					order_real.Stock_count -= count
					count = 0
				}else {
					// 比委托的单量大， 全部成交， 修改订单记录状态为完成。
					trade_detail.Stock_count = order_real.Stock_count
					order_real.Order_status = model.ORDER_STATUS_FINISH
					// 剩餘的股 继续下个档位交易
					count -= order_real.Stock_count
					order_real.Stock_count = 0
				}

				// 保存订单
				db.DBSession.Save(&order_real)

				// 将交易明细 发送到 清算系统（clearing）
				manager.Send2Clearing(trade_detail, 2)

				//  本档 为全部完成， 继续入队。
				//  本档交易都未能全部完成， 不会再有
				if order_real.Stock_count > 0 {
					this.orders[order_real.ID] = order_real  //  更新 orders
					sequence_que <- order_real

					msg_tmp, err := json.Marshal(order_real)
					if err != nil {
						log.Error("json.Marshal(order_real) err:%v", err)
					}else {
						log.Info("sequence_que <- order_real: %s", string(msg_tmp))
					}
				}else {
					// 本次的订单需要从 orders 中移除
					delete(this.orders, order_real.ID)
				}

				// 外盘的本笔交易 已经完成
				if count == 0 {
					break
				}
			}else {
				log.Info("BuyQue当前成交价：%f, 盘口价：%f", price, order_real.Stock_price)

				// 后面也不会有 符合条件的订单了， 直接退出
				break
			}
		}else {
			log.Error("buy que item not handler: %T, %v", it.Value(), it.Value())
			return errors.New(fmt.Sprintf("buy que item not handler: %T, %v", it.Value(), it.Value()))
		}
	}
	return nil
}

// 匹配 买卖 订单
// stock_code: 股票代码
func (this *Handler)Match(stock_code string, price float64, count int) (err error) {
	// 对操作加锁
	this.Lock()
	defer this.Unlock()

	que_sale, ok := this.PlateSale[stock_code]
	if ok {
		err = this.matchHandlerSaleQue(que_sale, price, count)
		if err != nil {
			log.Error("matchHandlerSaleQue err:%v", err)
		}
	}

	que_buy, ok := this.PlateBuy[stock_code]
	if ok {
		err = this.matchHandlerBuyQue(que_buy, price, count)
		if err != nil {
			log.Error("matchHandlerBuyQue err:%v", err)
		}
	}

	return nil
}