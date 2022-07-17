package server

import (
	"net"

	"github.com/orzen/steve/pkg/plugin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	Addr string
	Srv  *grpc.Server

	InitHooks        []plugin.InitHook
	PreRequestHooks  []plugin.PreReqHook
	PostRequestHooks []plugin.PostReqHook
}

func NewServer(addr string, srv *grpc.Server, plugins []plugin.Plugin) *Server {
	// RegisterPlugins

	// Run plugin InitHooks

	// RegisterPlugins as interceptors

	return &Server{
		Addr: addr,
		Srv:  srv,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal().Err(err).Msgf("listen on '%s'", s.Addr)
	}

	s.Srv.Serve(listener)
}
