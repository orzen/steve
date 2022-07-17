package spb

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/orzen/steve/pkg/conv"
	"github.com/orzen/steve/pkg/utils"
	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

func LoadDir(resourceDir string) (*parser.Proto, error) {
	var f filepath.WalkFunc
	var p *parser.Proto
	var err error

	f = func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".spb") && !info.IsDir() {
			if p, err = LoadFile(path); err != nil {
				return fmt.Errorf("load file: %v", err)
			}
		}

		return nil
	}

	if err = filepath.Walk(resourceDir, f); err != nil {
		return nil, err
	}

	return p, nil
}

func LoadFile(file string) (*parser.Proto, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.New("read resource file")
	}

	return Parse(data)
}

func Parse(data []byte) (*parser.Proto, error) {
	r := strings.NewReader(string(data))
	return protoparser.Parse(r)
}

type Sorted struct {
	Opts []*parser.Option
	Msgs []*parser.Message
}

func sortEntries(p *parser.Proto) (Sorted, error) {
	opts := []*parser.Option{}
	msgs := []*parser.Message{}

	for _, v := range p.ProtoBody {
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		switch t.Name() {
		case "Option":
			opts = append(opts, v.(*parser.Option))
		case "Message":
			msgs = append(msgs, v.(*parser.Message))
		default:
			return Sorted{}, fmt.Errorf("unsupported type: %s", t.Name())
		}
	}

	return Sorted{
		Opts: opts,
		Msgs: msgs,
	}, nil
}

func messagesToResources(msgs []*parser.Message) ([]*conv.Resource, error) {
	resources := []*conv.Resource{}

	for _, m := range msgs {
		res := &conv.Resource{
			Name: m.MessageName,
		}

		for _, e := range m.MessageBody {
			t := utils.Type(e)
			switch t {
			case "Field":
				res.Fields = append(res.Fields, e.(*parser.Field))
			case "Message":
				msg := e.(*parser.Message)
				if msg.MessageName == "Meta" {
					res.Meta = msg
				} else {
					res.Messages = append(res.Messages, msg)
				}
			case "MapField":
				res.Maps = append(res.Maps, e.(*parser.MapField))
			case "Enum":
				res.Enums = append(res.Enums, e.(*parser.Enum))
			case "Option":
				opt := e.(*parser.Option)
				if opt.OptionName == "operations" {
					res.SetOperations(opt)
				} else {
					res.Options = append(res.Options, opt)
				}
			case "Oneof":
				res.Oneofs = append(res.Oneofs, e.(*parser.Oneof))
			default:
				return nil, fmt.Errorf("unsupported resource content: %s", t)
			}
		}

		resources = append(resources, res)
	}

	return resources, nil
}

func ResourceLookup(name string, resources []*conv.Resource) (*conv.Resource, bool) {
	for _, r := range resources {
		if r.Name == name {
			return r, true
		}
	}
	return nil, false
}

func addOptions(resources []*conv.Resource, opts []*parser.Option) error {
	for _, o := range opts {
		split := strings.SplitN(o.OptionName, "_", 2)
		optionName, resourceName := split[0], split[1]

		res, ok := ResourceLookup(resourceName, resources)
		if !ok {
			continue
		}

		switch optionName {
		case "operations":
			constant := strings.Trim(o.Constant, "\"")
			res.Operations = strings.Split(constant, ",")
		default:
			return fmt.Errorf("invalid option: %s", optionName)
		}
	}

	return nil
}

func SpbToProto(p *parser.Proto) (*conv.Proto, error) {
	var res []*conv.Resource
	var err error
	var sorted Sorted

	proto := conv.NewProto()

	sorted, err = sortEntries(p)
	if err != nil {
		return proto, fmt.Errorf("convert spb to resource: %v", err)
	}

	res, err = messagesToResources(sorted.Msgs)
	if err != nil {
		return proto, fmt.Errorf("create resources: %v", err)
	}

	// Add default messages for each resource
	for _, r := range res {
		proto.Messages = append(proto.Messages, &parser.Message{
			MessageName: fmt.Sprintf("%sList", r.Name),
			MessageBody: []parser.Visitee{
				&parser.Field{
					IsRepeated:  true,
					Type:        fmt.Sprintf("%s.Meta", r.Name),
					FieldName:   "entries",
					FieldNumber: "1",
				},
			},
		})
	}

	if err = addOptions(res, sorted.Opts); err != nil {
		return proto, fmt.Errorf("add options: %v", err)
	}

	proto.Services = []*conv.Service{
		{
			Name:      "API",
			Resources: res,
		},
	}

	return proto, nil
}
