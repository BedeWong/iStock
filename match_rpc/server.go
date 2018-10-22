package match_rpc

import (
	"net/http"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"github.com/gpmgo/gopm/modules/log"
	"os"
	"github.com/BideWong/iStock/conf"
)

func rpc_server_start(addr string) error{
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
	err := rpc_server_start(conf.Data.Rpc.Addr)
	if err != nil {
		log.Error("rpc server start err:", err)
		os.Exit(-1)
	}
}