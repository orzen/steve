package steve

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/orzen/steve/pkg/pb"
	"github.com/orzen/steve/pkg/spb"
	"github.com/orzen/steve/pkg/tpl"
	"github.com/orzen/steve/pkg/utils"
	"github.com/rs/zerolog/log"
)

const (
	GlueFilename    = "steve.glue.go"
	GlueTplFilename = "glue.tpl"

	PbFilename    = "steve.proto"
	PbTplFilename = "proto.tpl"

	MainFilename    = "main.go"
	MainTplFilename = "main.tpl"

	BinFilename = "api"
)

func Build(buildDir, apiDir, tplDir string) error {
	// Ensure prerequisites
	if err := os.Mkdir(buildDir, 0755); err != nil {
		if !strings.HasSuffix(err.Error(), "file exists") {
			log.Fatal().Err(err).Msg("create api dir")
		}
	}

	// Convert steve protobuf to regular protobuf
	resources, err := spb.LoadDir(apiDir)
	if err != nil {
		log.Fatal().Err(err).Msg("load resource dir")
	}

	for k, r := range resources {
		if err := r.Finalize(); err != nil {
			log.Fatal().Err(err).Msgf("finalize resource '%s'", k)
		}
		log.Debug().Interface("resource", r).Msg("resources")
	}

	// Generate Go gRPC stubs from protobuf
	pbFile := filepath.Join(buildDir, PbFilename)
	pbTpl := filepath.Join(tplDir, PbTplFilename)
	pb.ToGo(buildDir, pbFile, pbTpl, resources)

	// Generate Go glue between Go gRPC stubs and steve
	glueFile := filepath.Join(buildDir, GlueFilename)
	glueTpl := filepath.Join(tplDir, GlueTplFilename)
	if err := tpl.TemplateFromFile(glueFile, glueTpl, resources); err != nil {
		log.Fatal().Err(err).Msg("generate glue")
	}

	// Generate main file
	mainFile := filepath.Join(buildDir, MainFilename)
	mainTpl := filepath.Join(tplDir, MainTplFilename)
	if err := tpl.TemplateFromFile(mainFile, mainTpl, resources); err != nil {
		log.Fatal().Err(err).Msg("generate main")
	}

	// Compile
	binFile := filepath.Join(buildDir, BinFilename)
	if err := utils.Compile(mainFile, binFile); err != nil {
		log.Fatal().Err(err).Msg("compile api")
	}

	return nil
}
