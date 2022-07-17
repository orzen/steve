package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/orzen/steve/pkg/steve"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	Version = "dev"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msgf("get workdir: %v", err)
	}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "build-dir",
			Value:    filepath.Join(wd, "_build"),
			Usage:    "Build directory",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "spec-dir",
			Value:    filepath.Join(wd, "api"),
			Usage:    "steve specification directory",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "app-name",
			Value:    "app",
			Usage:    "Application name",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "app-version",
			Value:    "v1.0.0",
			Usage:    "Application version string",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "backend-mongodb",
			Usage:    "Backend plugin for MongoDB",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "with-client",
			Usage:    "Produce a client app",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "verbose",
			Usage:    "Increase verbosity",
			Required: false,
		},
	}

	app := &cli.App{
		Name:    "steve",
		Version: Version,
		Flags:   flags,
		Action: func(c *cli.Context) error {
			backend := ""

			if c.Bool("verbose") {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			// Select backend plugin
			if len(c.String("backend-mongodb")) != 0 {
				backend = "mongodb"
			}
			if backend == "" {
				return errors.New("one backend plugin is required")
			}

			buildCfg := steve.BuildCfg{
				AppName:      c.String("app-name"),
				AppVersion:   c.String("app-version"),
				SteveVersion: Version,
				Backend:      backend,
				BuildDir:     c.String("build-dir"),
				SpecDir:      c.String("spec-dir"),
				TplDir:       filepath.Join(wd, "tpl"),
				WithClient:   c.Bool("with-client"),
			}

			return steve.Build(buildCfg)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("run steve")
	}
}
