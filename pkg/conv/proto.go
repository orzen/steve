package conv

import (
	"fmt"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type Proto struct {
	Options  []*parser.Option
	Package  *parser.Package
	Services []*Service
	Messages []*parser.Message
	Syntax   *parser.Syntax
}

func NewProto() *Proto {
	return &Proto{
		Options: []*parser.Option{
			{
				OptionName: "go_package",
				Constant:   "\"github.com/orzen/steve/_build/api\"",
			},
		},
		Package: &parser.Package{
			Name: "api",
		},
		Syntax: &parser.Syntax{
			ProtobufVersion:      "proto3",
			ProtobufVersionQuote: "\"proto3\"",
		},
	}
}

func (p *Proto) Resources() []*Resource {
	res := []*Resource{}

	for _, s := range p.Services {
		for _, r := range s.Resources {
			res = append(res, r)
		}
	}

	return res
}

func (p *Proto) String() string {
	b := &strings.Builder{}

	WriteString(b, fmt.Sprintf("syntax = %s;\n", p.Syntax.ProtobufVersionQuote))
	WriteString(b, fmt.Sprintf("package %s;\n", p.Package.Name))
	WriteRune(b, '\n')

	for _, o := range p.Options {
		WriteString(b, OptionToA(o, 0))
	}

	WriteRune(b, '\n')

	for _, s := range p.Services {
		WriteString(b, s.String())
	}

	for _, m := range p.Messages {
		WriteString(b, MessageToA(m, 0))
		WriteRune(b, '\n')
	}

	return b.String()
}
