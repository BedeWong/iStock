package conf_test

import (
	"testing"
	"codehome/match/conf"
)

func TestfnInit(t *testing.T) {
	conf.Init()

	if conf.GetConfig() != conf.GetConfig() {
		t.Logf("两次获取的对象不唯一:\n")
		t.Logf("obj1:%v\n", conf.GetConfig())
		t.Logf("obj2:%v\n", conf.GetConfig())
	}
}


//func Test2(t *testing.T) {
//	t.Error("就是通不过。")
//}