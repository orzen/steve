package server

import (
	"fmt"
	"net"

	"github.com/orzen/steve/srv/plugin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type Input struct {
	Backend      plugin.BackendPlugin
	Cli          *cli.Context
	Interceptors []plugin.InterceptorPlugin
	Srv          *grpc.Server
}

type Server struct {
	Addr         string
	Backend      plugin.BackendPlugin
	Interceptors []plugin.InterceptorPlugin
	Srv          *grpc.Server
}

func NewServer(in *Input) *Server {
	// Fetch values from the cli parser
	listenAddr := in.Cli.String("listen-addr")
	listenPort := in.Cli.Int("listen-port")

	// Start backend
	if err := in.Backend.Start(in.Cli); err != nil {
		log.Fatal().Err(err).
			Msgf("start backend '%s'", in.Backend.Cfg().Name)
	}

	// Configure interceptors (give them the possibility to fetch command
	// line values from the cli context).
	for _, i := range in.Interceptors {
		if err := i.Configure(in.Cli); err != nil {
			log.Fatal().Err(err).
				Msgf("configure interceptor '%s': %v", i.Cfg().Name, err)
		}
	}

	return &Server{
		Addr:         fmt.Sprintf("%s:%d", listenAddr, listenPort),
		Backend:      in.Backend,
		Interceptors: in.Interceptors,
		Srv:          in.Srv,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal().Err(err).Msgf("listen on '%s'", s.Addr)
	}

	return s.Srv.Serve(listener)
}

func (s *Server) Set(t string, r interface{}) error {
	if err := s.Backend.Set(t, r); err != nil {
		log.Error().Err(err).
			Str("op", "set").
			Str("type", t).
			Interface("resource", r).Send()
		return err
	}

	return nil
}

func (s *Server) Get(t string, m interface{}, retR interface{}) error {
	if err := s.Backend.Get(t, m, retR); err != nil {
		log.Error().Err(err).
			Str("op", "get").
			Str("type", t).
			Interface("filter", m).Send()
		return err
	}

	return nil
}

func (s *Server) List(t string, m interface{}, retArrM interface{}) error {
	if err := s.Backend.List(t, m, retArrM); err != nil {
		log.Error().Err(err).
			Str("op", "list").
			Str("type", t).
			Interface("filter", m).Send()
		return err
	}

	return nil
}

func (s *Server) Delete(t string, m interface{}, retR interface{}) error {
	if err := s.Backend.Delete(t, m, retR); err != nil {
		log.Error().Err(err).
			Str("op", "delete").
			Str("type", t).
			Interface("filter", m).Send()
		return err
	}

	return nil
}
