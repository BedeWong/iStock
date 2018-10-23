package match_rpc

import (
	"net/http"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"github.com/gpmgo/gopm/modules/log"
	"os"
	"github.com/BideWong/iStock/conf"

	"github.com/BideWong/iStock/match_rpc/service"
	"fmt"
)

func register_service() {
	order := new(service.Order)

	rpc.RegisterName("order", order)

	fmt.Println("register service ok.")
}

func rpc_server_start(addr string) error{

	register_service()

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
		err := rpc_server_start(conf.Data.Rpc.Addr)
		if err != nil {
			log.Error("rpc server start err:", err)
			os.Exit(-1)
		}
	}()
}