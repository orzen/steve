// vim: ft=go
package main

import (
	"fmt"

	"github.com/orzen/steve/pkg/server"
	"google.golang.org/grpc"
)

func main() {
	addr := fmt.Sprintf("%s:%d", {{.Addr}}, {{.Port}})
	var opts []grpc.ServerOption

	grpcSrv := grpc.NewServer(opts...)

	srv := server.NewServer(addr, grpcSrv, []plugin.Plugin{})
}
