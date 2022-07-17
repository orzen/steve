package steve

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/orzen/steve/pkg/conv"
	"github.com/orzen/steve/pkg/pb"
	"github.com/orzen/steve/pkg/spb"
	"github.com/orzen/steve/pkg/tpl"
	"github.com/orzen/steve/pkg/utils"
	"github.com/rs/zerolog/log"
)

const (
	GlueDirname     = "glue"
	GlueFilename    = "steve.glue.go"
	GlueTplFilename = "glue.tpl"

	PbFilename    = "steve.proto"
	PbTplFilename = "proto.tpl"

	SrvDirname     = "srv"
	SrvFilename    = "main.go"
	SrvTplFilename = "srv.tpl"

	CliDirname     = "cli"
	CliFilename    = "main.go"
	CliTplFilename = "cli.tpl"

	APIDirname = "api"
)

type BuildCfg struct {
	AppName      string
	AppVersion   string
	SteveVersion string
	Backend      string
	BuildDir     string
	SpecDir      string
	TplDir       string
	WithClient   bool
}

func Build(c BuildCfg) error {
	// Ensure prerequisites
	if err := os.Mkdir(c.BuildDir, 0755); err != nil {
		if !strings.HasSuffix(err.Error(), "file exists") {
			log.Fatal().Err(err).Msg("create api dir")
		}
	}

	// Sub-dirs
	apiDir := filepath.Join(c.BuildDir, APIDirname)
	glueDir := filepath.Join(c.BuildDir, GlueDirname)
	srvDir := filepath.Join(c.BuildDir, SrvDirname)
	cliDir := filepath.Join(c.BuildDir, CliDirname)

	dirs := []string{apiDir, glueDir, srvDir}

	if c.WithClient {
		dirs = append(dirs, cliDir)
	}

	for _, d := range dirs {
		if err := os.Mkdir(d, 0755); err != nil {
			if !strings.HasSuffix(err.Error(), "file exists") {
				log.Fatal().Err(err).Msgf("create 's' dir", d)
			}
		}
	}

	// Convert steve protobuf to regular protobuf
	parsed, err := spb.LoadDir(c.SpecDir)
	if err != nil {
		log.Fatal().Err(err).Msg("load resource dir")
	}

	proto, err := spb.SpbToProto(parsed)
	if err != nil {
		log.Fatal().Err(err).Msg("convert parser result to internal representation")
	}

	tplCfg := &tpl.Cfg{
		AppName:      c.AppName,
		AppVersion:   c.AppVersion,
		SteveVersion: c.SteveVersion,
		Backend:      c.Backend,
		Proto:        proto,
	}

	// Generate Go gRPC stubs from protobuf
	pbFile := filepath.Join(apiDir, PbFilename)
	pb.ToGo(apiDir, pbFile, tplCfg)

	// Generate Go glue between Go gRPC stubs and steve
	glueFile := filepath.Join(glueDir, GlueFilename)
	glueTpl := filepath.Join(c.TplDir, GlueTplFilename)
	if err := tpl.TemplateFromFile(glueFile, glueTpl, tplCfg); err != nil {
		log.Fatal().Err(err).Msg("generate glue")
	}

	// Generate server's main file
	srvFile := filepath.Join(srvDir, SrvFilename)
	srvTpl := filepath.Join(c.TplDir, SrvTplFilename)
	if err := tpl.TemplateFromFile(srvFile, srvTpl, tplCfg); err != nil {
		log.Fatal().Err(err).Msg("generate server main")
	}

	// Compile
	srvBin := filepath.Join(c.BuildDir, fmt.Sprintf("%s", c.AppName))
	if err := utils.Compile(srvFile, srvBin); err != nil {
		log.Fatal().Err(err).Msg("compile api srv")
	}

	if c.WithClient {
		tplCfg.Funcs = template.FuncMap{
			"MetaToFlag":       conv.MetaToFlag,
			"ResourceToFlag":   conv.ResourceToFlag,
			"MetaToAction":     conv.MetaToAction,
			"ResourceToAction": conv.ResourceToAction,
		}

		// Generate client's main file
		cliFile := filepath.Join(cliDir, CliFilename)
		cliTpl := filepath.Join(c.TplDir, CliTplFilename)
		if err := tpl.TemplateFromFile(cliFile, cliTpl, tplCfg); err != nil {
			log.Fatal().Err(err).Msg("generate cli main")
		}

		// Compile
		cliBin := filepath.Join(c.BuildDir, fmt.Sprintf("%sctl", c.AppName))
		if err := utils.Compile(cliFile, cliBin); err != nil {
			log.Fatal().Err(err).Msg("compile api cli")
		}
	}

	return nil
}
