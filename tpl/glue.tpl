// vim: ft=go

{{ $Resources := .Proto.Resources }}

package glue

import (
	"context"
	"errors"

	"github.com/orzen/steve/_build/api"
	"github.com/orzen/steve/srv/server"
)

type Glue struct {
	api.APIServer

	Srv *server.Server
}

type ResourceType interface {
	{{- if eq 1 (len $Resources) }}
		{{ printf "\n\tapi.%s" (index .Resources 0).Name }}
	{{ else }}
		{{- range $i, $res := $Resources }}
			{{- if eq 0 $i -}}
				{{ printf "\n\tapi.%s" $res.Name }}
			{{- else -}}
				{{ printf " | api.%s" $res.Name }}
			{{- end -}}
		{{ end }}
	{{- end }}
}

type MetaType interface {
	{{- if eq 1 (len $Resources) }}
		{{ printf "\n\tapi.%s_Meta" (index .Resources 0).Name }}
	{{ else }}
		{{- range $i, $res := $Resources }}
			{{- if eq 0 $i -}}
				{{ printf "\n\tapi.%s_Meta" $res.Name }}
			{{- else -}}
				{{ printf " | api.%s_Meta" $res.Name }}
			{{- end -}}
		{{ end }}
	{{- end }}
}

type ListType interface {
	{{- if eq 1 (len $Resources) }}
		{{ printf "\n\tapi.%s_Meta" (index .Resources 0).Name }}
	{{ else }}
		{{- range $i, $res := $Resources }}
			{{- if eq 0 $i -}}
				{{ printf "\n\tapi.%s_Meta" $res.Name }}
			{{- else -}}
				{{ printf " | api.%s_Meta" $res.Name }}
			{{- end -}}
		{{ end }}
	{{- end }}
}

func New(srv *server.Server) *Glue {
	return &Glue{ Srv: srv }
}

{{ range $res := $Resources }}
{{ printf "/* -- %s -- */\n" $res.Name }}
{{- range $op := $res.Operations -}}

{{ if eq "set" $op }}
func (g *Glue) Set{{$res.Name}}(ctx context.Context, r *api.{{$res.Name}}) (*api.{{$res.Name}}_Meta, error) {
	if err := g.Srv.Set("{{$res.Name}}", r); err != nil {
		return nil, errors.New("failed to set {{$res.Name}}")
	}

	return r.Meta, nil
}
{{ end }}

{{- if eq "get" $op }}
func (g *Glue) Get{{$res.Name}}(ctx context.Context, m *api.{{$res.Name}}_Meta) (*api.{{$res.Name}}, error) {
	r := &api.{{$res.Name}}{}

	if err := g.Srv.Get("{{$res.Name}}", m, &r); err != nil {
		return nil, errors.New("failed to get {{$res.Name}}")
	}

	return r, nil
}
{{ end }}

{{- if eq "list" $op }}
func (g *Glue) List{{$res.Name}}(ctx context.Context, m *api.{{$res.Name}}_Meta) (*api.{{$res.Name}}List, error) {
	r := []*api.{{$res.Name}}_Meta{}

	if err := g.Srv.List("{{$res.Name}}", m, r); err != nil {
		return nil, errors.New("failed to list {{$res.Name}}")
	}

	return &api.{{$res.Name}}List {
		Entries: r,
	}, nil
}
{{ end }}

{{- if eq "delete" $op }}
func (g *Glue) Delete{{$res.Name}}(ctx context.Context, m *api.{{$res.Name}}_Meta) (*api.{{$res.Name}}, error) {
	r := &api.{{$res.Name}}{}

	if err := g.Srv.Delete("{{$res.Name}}", m, &r); err != nil {
		return nil, errors.New("failed to delete {{$res.Name}}")
	}

	return r, nil
}
{{ end }}

{{- end -}}
{{- end }}
func (g *Glue) mustEmbedUnimplementedAPIServer(){}
