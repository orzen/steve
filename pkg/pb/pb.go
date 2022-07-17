package pb

import (
	"github.com/orzen/steve/pkg/resource"
	"github.com/orzen/steve/pkg/tpl"
	"github.com/rs/zerolog/log"
)

func ToGo(workDir, pbFile, pbTpl string, resources map[string]*resource.Resource) {
	if err := tpl.TemplateFromFile(pbFile, pbTpl, resources); err != nil {
		log.Fatal().Err(err).Msgf("generate protobuf: %v", err)
	}

	if err := Protoc(workDir, pbFile); err != nil {
		log.Fatal().Err(err).Msgf("generate go from protobuf: %v", err)
	}
}
