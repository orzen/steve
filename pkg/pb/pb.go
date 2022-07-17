package pb

import (
	"os"

	"github.com/orzen/steve/pkg/conv"
	"github.com/orzen/steve/pkg/tpl"
	"github.com/rs/zerolog/log"
)

func ProtoToFile(path string, p *conv.Proto) error {
	return os.WriteFile(path, []byte(p.String()), 0644)
}

func ToGo(workDir, pbFile string, tplCfg *tpl.Cfg) {
	if err := ProtoToFile(pbFile, tplCfg.Proto); err != nil {
		log.Fatal().Err(err).Msgf("generate protobuf: %v", err)
	}
	// TODO remove
	//if err := tpl.TemplateFromFile(pbFile, pbTpl, tplCfg); err != nil {
	//	log.Fatal().Err(err).Msgf("generate protobuf: %v", err)
	//}

	if err := Protoc(workDir, pbFile); err != nil {
		log.Fatal().Err(err).Msgf("generate go from protobuf: %v", err)
	}
}
