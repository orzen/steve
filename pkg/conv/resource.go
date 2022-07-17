package conv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

var (
	allowedOperations = []string{"set", "get", "list", "delete"}
)

const (
	TypeEnum = iota
	TypeMessage
)

type Resource struct {
	Name       string
	Meta       *parser.Message
	Enums      []*parser.Enum
	Fields     []*parser.Field
	Maps       []*parser.MapField
	Messages   []*parser.Message
	Options    []*parser.Option
	Oneofs     []*parser.Oneof
	Operations []string
	Types      map[string]int
}

func (r *Resource) SetOperations(o *parser.Option) {
	ops := strings.Split(o.Constant, ",")

	for i, e := range ops {
		ops[i] = strings.Trim(e, " \t\n\r\"")
	}

	r.Operations = ops
}

func (r *Resource) OperationString() string {
	return OperationsToA(r.Name, r.Operations)
}

func (r *Resource) MessageString() string {
	b := &strings.Builder{}
	r.Types = map[string]int{}
	level := 0

	WriteString(b, fmt.Sprintf("message %s {\n", r.Name))

	for _, o := range r.Options {
		WriteString(b, OptionToA(o, level+1))
	}

	for _, e := range r.Enums {
		r.Types[e.EnumName] = TypeEnum
		WriteString(b, EnumToA(e, level+1))
	}

	WriteString(b, MessageToA(r.Meta, level+1))

	for _, m := range r.Messages {
		r.Types[m.MessageName] = TypeMessage
		WriteString(b, MessageToA(m, level+1))
	}

	l := len(r.Fields) + len(r.Maps)
	fields := make([]string, l)

	// Append fields to an array based on number to get them in indended
	// order.
	for _, f := range r.Fields {
		i, err := strconv.Atoi(f.FieldNumber)
		if err != nil {
			log.Fatal().Err(err).Msg("convert field number to int")
		}
		fields[i-1] = FieldToA(f, level+1)
	}

	for _, m := range r.Maps {
		i, err := strconv.Atoi(m.FieldNumber)
		if err != nil {
			log.Fatal().Err(err).Msg("convert field number to int")
		}
		fields[i-1] = MapToA(m, level+1)
	}

	for _, f := range fields {
		WriteString(b, f)
	}

	WriteRune(b, '}')

	return b.String()
}
