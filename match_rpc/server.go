package match_rpc

import (
	"os"
	"net/http"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gpmgo/gopm/modules/log"
	"github.com/BedeWong/iStock/conf"
	"github.com/BedeWong/iStock/match_rpc/service"
)

// 服务注册.
func registerService() {
	order := new(service.OrderService)

	rpc.RegisterName("orderService", order)

	log.Info("register service ok.")
}


func rpcServerStart(addr string) error{
	registerService()

	http.HandleFunc(conf.Data.Rpc.Pattern, func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}

		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	return http.ListenAndServe(addr, nil)
}


func init(){
	go func() {
		err := rpcServerStart(conf.Data.Rpc.Addr)
		if err != nil {
			log.Error("rpc server start err:", err)
			os.Exit(-1)
		}
	}()

	log.Info("rpc server started ok.")
}