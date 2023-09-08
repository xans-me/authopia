package configuration

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	// used for side effects
	_ "net/http/pprof"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// ListenGRPC is a func to start http server
func ListenGRPC(conf *AppConfig) net.Listener {
	addr := conf.App.Host + ":" + conf.App.Port
	log.Info("run GRPC & HTTP Transport on " + addr)

	lis, err := net.Listen(conf.App.Protocol, addr)
	if err != nil {
		log.Fatal(err)
	}
	return lis
}
