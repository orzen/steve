package conv

// This file contains utility functions to convert parser.Field and
// parser.Message to cli.Flag strings that is used for templating the client.

import (
	"fmt"
	"strings"

	"github.com/orzen/steve/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type ResourceEntries struct {
	MetaEntries []*Entry
	Entries     []*Entry
}

type Entry struct {
	Name     string
	FlagName string
	Type     string
	Default  string
	Required bool
}

func (e *Entry) ToFlag(level int) string {
	b := &strings.Builder{}

	indent := Indentation(level)

	WriteString(b, fmt.Sprintf("%s&cli.%sFlag{\n", indent, e.Type))
	WriteString(b, fmt.Sprintf("%s\tName:     \"%s\",\n", indent, e.FlagName))
	WriteString(b, fmt.Sprintf("%s\tValue:    %s,\n", indent, e.Default))
	WriteString(b, fmt.Sprintf("%s\tRequired: %t,\n", indent, e.Required))
	WriteString(b, fmt.Sprintf("%s},\n", indent))

	return b.String()
}

// TODO add string to int lookup for Enums in the action

func (e *Entry) ToField(level int) string {
	return fmt.Sprintf("%s: c.%s(\"%s\"),\n",
		strings.Title(e.Name), e.Type, e.FlagName)
}

func PrefixName(prefix, name string) string {
	if prefix != "" {
		return fmt.Sprintf("%s-%s", strings.ToLower(prefix), strings.ToLower(name))
	}
	return strings.ToLower(name)
}

// === FLAG UTILS ===

func FieldToEntry(f *parser.Field, level int, prefix string, required bool) *Entry {
	t := PbToGoType(f.Type)
	typ := ""
	def := ""

	switch t {
	case "bool":
		typ = "Bool"
		def = "false"
	case "[]byte":
		typ = "Path"
		def = "\"\""
	case "float32", "float64":
		typ = "Float64"
		def = "0"
	case "int32":
		typ = "Int"
		def = "0"
	case "int64":
		typ = "Int64"
		def = "0"
	case "uint32":
		typ = "Uint"
		def = "0"
	case "uint64":
		typ = "Uint64"
		def = "0"
	case "string":
		typ = "String"
		def = "\"\""
	default:
		log.Fatal().Msgf("unhandled cli type: %s", t)
	}

	return &Entry{
		Name:     f.FieldName,
		FlagName: PrefixName(prefix, f.FieldName),
		Type:     typ,
		Default:  def,
		Required: required,
	}
}

func EnumToEntry(e *parser.Enum, level int, prefix string, required bool) *Entry {
	def := ""
	for _, v := range e.EnumBody {
		if utils.Type(v) != "EnumField" {
			continue
		}

		f := v.(*parser.EnumField)
		fmt.Println("ENUM NUMBER", f.Number)
		if f.Number == "1" {
			def = f.Ident
			break
		}
	}

	return &Entry{
		Name:     e.EnumName,
		FlagName: PrefixName(prefix, e.EnumName),
		Type:     "String",
		Default:  fmt.Sprintf("\"%s\"", def),
		Required: required,
	}
}

// TODO change MessageToFlag > MessageToEntry
func MessageToFlag(m *parser.Message, level int, prefix string, required bool) (string, error) {
	b := &strings.Builder{}
	pre := PrefixName(prefix, m.MessageName)

	for _, e := range m.MessageBody {
		t := utils.Type(e)
		switch t {
		case "Field":
			f := FieldToFlag(e.(*parser.Field), level, pre, required)

			WriteString(b, f)
		case "Message":
			f, err := MessageToFlag(e.(*parser.Message), level, pre, required)
			if err != nil {
				return "", fmt.Errorf("message flag: %v", err)
			}

			WriteString(b, f)
		case "Enum":
			WriteString(b, EnumToFlag(e.(*parser.Enum), level, pre, required))
		default:
			return "", fmt.Errorf("no flag for type(message): '%s'", t)
		}
	}

	return b.String(), nil
}

// TODO change MetaToFlag > MetaToEntry
func MetaToFlag(r *Resource, level int) string {
	str, err := MessageToFlag(r.Meta, level, "", false)
	if err != nil {
		log.Fatal().Err(err).Msg("generate meta flags")
	}
	return str
}

// TODO change ResourceToFlag > ResourceToEntry
func ResourceToFlag(r *Resource, level int) string {
	b := &strings.Builder{}

	// Append meta
	metaFlags, err := MessageToFlag(r.Meta, level, "", true)
	if err != nil {
		log.Fatal().Err(err).Msg("generate meta flags")
	}
	WriteString(b, metaFlags)

	for _, m := range r.Messages {
		msgFlags, err := MessageToFlag(m, level, "", true)
		if err != nil {
			log.Fatal().Err(err).Msg("generate message flags")
		}
		WriteString(b, msgFlags)
	}

	for _, f := range r.Fields {
		fieldFlags := FieldToEntry(f, level, "", true)

		WriteString(b, fieldFlags)
	}

	return b.String()
}

// === ACTION UTILS ===

func MessageToAction(m *parser.Message, level int) string {
	return ""
}

func MetaToAction(r *Resource, level int) string {
	return ""
}

func ResourceToAction(r *Resource, level int) string {
	b := &strings.Builder{}

	return b.String()
}
