package main

import (
	"os"
	"path/filepath"

	"github.com/orzen/steve/pkg/steve"
	"github.com/rs/zerolog/log"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msgf("get workdir: %v", err)
	}

	steve.Build(filepath.Join(wd, "_build"),
		filepath.Join(wd, "api"),
		filepath.Join(wd, "tpl"))
}
