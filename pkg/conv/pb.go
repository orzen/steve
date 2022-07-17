package conv

import (
	"fmt"
	"strings"

	"github.com/orzen/steve/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

func BlockToA(ident, name string, level int, content func(int) string) string {
	b := &strings.Builder{}
	indent := Indentation(level)

	WriteString(b, fmt.Sprintf("%s%s %s {\n", indent, ident, name))
	WriteString(b, content(level))
	WriteString(b, fmt.Sprintf("%s}\n", indent))

	return b.String()
}

func LineOptionsToA(opts map[string]string) string {
	if len(opts) == 0 {
		return ""
	}

	b := &strings.Builder{}
	WriteString(b, " [")

	acc := []string{}
	for k, v := range opts {
		acc = append(acc, fmt.Sprintf("%s=%s", k, v))
	}

	WriteString(b, strings.Join(acc, ","))
	WriteRune(b, ']')

	return b.String()
}

func LineToA(mod, typ, name, number string, opts map[string]string, level int) string {
	b := &strings.Builder{}

	WriteString(b, Indentation(level))

	if mod != "" {
		WriteString(b, fmt.Sprintf("%s ", mod))
	}
	if typ != "" {
		WriteString(b, fmt.Sprintf("%s ", typ))
	}

	WriteString(b, fmt.Sprintf("%s = %s", name, number))
	WriteString(b, LineOptionsToA(opts))
	WriteString(b, ";\n")

	return b.String()
}

// === LINE TYPES ===

func FieldToA(f *parser.Field, level int) string {
	mod := ""
	if f.IsRepeated {
		mod = "repeated"
	}

	opts := map[string]string{}
	for _, o := range f.FieldOptions {
		opts[o.OptionName] = o.Constant
	}

	return LineToA(mod, f.Type, f.FieldName, f.FieldNumber, opts, level)
}

func MapToA(m *parser.MapField, level int) string {
	typ := fmt.Sprintf("map<%s, %s>", m.KeyType, m.Type)

	opts := map[string]string{}
	for _, o := range m.FieldOptions {
		opts[o.OptionName] = o.Constant
	}

	return LineToA("", typ, m.MapName, m.FieldNumber, opts, level)
}

// === BLOCK TYPES ===

func EnumToA(e *parser.Enum, level int) string {
	enumContent := func(level int) string {
		b := &strings.Builder{}

		for _, v := range e.EnumBody {
			if utils.Type(v) != "EnumField" {
				continue
			}

			val := v.(*parser.EnumField)

			opts := map[string]string{}
			for _, o := range val.EnumValueOptions {
				opts[o.OptionName] = o.Constant
			}

			WriteString(b, LineToA("", "", val.Ident, val.Number, opts, level+1))
		}

		return b.String()
	}

	return BlockToA("enum", e.EnumName, level, enumContent)
}

func MessageToA(m *parser.Message, level int) string {
	messageContent := func(level int) string {
		b := &strings.Builder{}
		for _, e := range m.MessageBody {
			t := utils.Type(e)
			switch t {
			case "Field":
				WriteString(b, FieldToA(e.(*parser.Field), level+1))
			case "Message":
				WriteString(b, MessageToA(e.(*parser.Message), level+1))
			default:
				log.Warn().Str("func", "MessageToA").Msgf("unhandled type: %s", t)
			}
		}
		return b.String()
	}

	return BlockToA("message", m.MessageName, level, messageContent)
}

func OneofToA(o *parser.Oneof, level int) string {
	oneofContent := func(level int) string {
		b := &strings.Builder{}
		for _, f := range o.OneofFields {
			opts := map[string]string{}
			for _, o := range f.FieldOptions {
				opts[o.OptionName] = o.Constant
			}

			WriteString(b, LineToA("", f.Type, f.FieldName, f.FieldNumber, opts, level+1))
		}
		return b.String()
	}

	return BlockToA("oneof", o.OneofName, level, oneofContent)
}

func OperationsToA(name string, ops []string) string {
	if len(ops) == 0 {
		return ""
	}

	b := &strings.Builder{}

	for _, o := range ops {
		switch o {
		case "set":
			WriteString(b, fmt.Sprintf("\trpc Set%s(%s) returns (%s.Meta);\n", name, name, name))
		case "get":
			WriteString(b, fmt.Sprintf("\trpc Get%s(%s.Meta) returns (%s);\n", name, name, name))
		case "list":
			WriteString(b, fmt.Sprintf("\trpc List%s(%s.Meta) returns (%sList);\n", name, name, name))
		case "delete":
			WriteString(b, fmt.Sprintf("\trpc Delete%s(%s.Meta) returns (%s);\n", name, name, name))
		default:
			log.Warn().Str("func", "ServiceToA").Msgf("invalid operation: %s", o)
		}
	}

	return b.String()
}

func OptionToA(o *parser.Option, level int) string {
	return LineToA("", "option", o.OptionName, o.Constant, map[string]string{}, level)
}
