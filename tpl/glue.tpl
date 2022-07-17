// vim: ft=go

package api

type API struct {
	APIServer
}

type ResourceType interface {
	{{- if eq 1 (len .Resources) }}
		{{ printf "\n\t%s" (index .Resources 0).Name }}
	{{ else }}
		{{- range $i, $res := .Resources }}
			{{- if eq 0 $i -}}
				{{ printf "\n\t%s" $res.Name }}
			{{- else -}}
				{{ printf " | %s" $res.Name }}
			{{- end -}}
		{{ end }}
	{{- end }}
}

type MetaType interface {
	{{- if eq 1 (len .Resources) }}
		{{ printf "\n\t%s_Meta" (index .Resources 0).Name }}
	{{ else }}
		{{- range $i, $res := .Resources }}
			{{- if eq 0 $i -}}
				{{ printf "\n\t%s_Meta" $res.Name }}
			{{- else -}}
				{{ printf " | %s_Meta" $res.Name }}
			{{- end -}}
		{{ end }}
	{{- end }}
}

{{ range $res := .Resources }}
{{ printf "/* -- %s -- */\n" $res.Name }}

{{- range $op := $res.Operations -}}

{{ if eq "set" $op }}
func (a *API) Set{{$res.Name}}(ctx context.Context *{{$res.Name}}) (*{{$res.Name}}_Meta, error) {
	return nil, nil
}
{{ end }}

{{- if eq "get" $op }}
func (a *API) Get{{$res.Name}}(ctx context.Context *{{$res.Name}}_Meta) (*{{$res.Name}}, error) {
	return nil, nil
}
{{ end }}

{{- if eq "list" $op }}
func (a *API) List{{$res.Name}}(ctx context.Context, *Void) ([]*{{$res.Name}}_Meta, error) {
	return []*{{$res.Name}}_Meta, nil
}
{{ end }}

{{- if eq "delete" $op }}
func (a *API) Delete{{$res.Name}}(ctx context.Context *{{$res.Name}}_Meta) (*{{$res.Name}}, error) {
	return nil, nil
}
{{ end }}

{{- end -}}
{{- end }}
func (a *API) mustEmbedUnimplementedAPIServer(){}
