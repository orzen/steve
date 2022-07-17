package conv

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func WriteString(b *strings.Builder, s string) {
	i, err := b.WriteString(s)
	if err != nil {
		log.Fatal().Int("size", i).Str("string", s).Err(err).Msg("write string to builder")
	}
}

func WriteRune(b *strings.Builder, r rune) {
	i, err := b.WriteRune(r)
	if err != nil {
		log.Fatal().Int("size", i).Str("rune", string(r)).Err(err).Msg("write rune to builder")
	}
}

func Indentation(level int) string {
	return strings.Repeat("\t", level)
}
