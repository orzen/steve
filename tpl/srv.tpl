// vim: ft=go
package main

import (
	"os"

	"github.com/orzen/steve/_build/api"
	"github.com/orzen/steve/_build/glue"
	cli2 "github.com/orzen/steve/srv/cli"
	"github.com/orzen/steve/srv/plugin"
	"github.com/orzen/steve/srv/server"
	"github.com/orzen/steve/srv/backends/{{ .Backend }}"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	SteveVersion = "{{ .SteveVersion }}"
	AppName      = "{{ .AppName }}"
	AppVersion   = "{{ .AppVersion }}"
)


func main() {
	// Plugin list
	backend := {{ .Backend }}.New()
	iCli := cli2.NewCli()
	interceptors := []plugin.InterceptorPlugin{}
	var opts []grpc.ServerOption

	// Register backend
	c, err := backend.Register()
	if err != nil {
		log.Fatal().Err(err).Msgf("register backend '%s': %v", backend.Cfg().Name, err)
	}
	iCli.AppendFlags(c.Cli)

	// Register interceptors
	for _, i := range interceptors {
		c, err := i.Register()
		if err != nil {
			log.Fatal().Err(err).Msgf("register interceptor '%s': %v", i.Cfg().Name, err)
		}
		iCli.AppendFlags(c.Cli)
		opts = append(opts, grpc.UnaryInterceptor(i.Cfg().Interceptor))
	}

	var action cli.ActionFunc = func(c *cli.Context) error {
		if c.Bool("pretty") {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		switch c.String("log-level") {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		// gRPC server have to be created here and not in
		// server.Server, otherwise it will create an import cycle
		// since glue implements the API and redirects calls to server.
		grpcSrv := grpc.NewServer(opts...)

		srv := server.NewServer(&server.Input{
			Backend: backend,
			Cli: c,
			Interceptors: interceptors,
			Srv: grpcSrv,
		})

		glue := glue.New(srv)

		api.RegisterAPIServer(grpcSrv, glue)

		return srv.Run()
	}

	if err := iCli.Finalize(AppName, AppVersion, SteveVersion, action); err != nil {
		log.Fatal().Err(err).Msgf("%s error", AppName)
	}
}
