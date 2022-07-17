// vim: ft=proto

syntax = "proto3";

option go_package = "github.com/orzen/steve/_build/api";

package api;

service API {
{{- range $res := .Resources -}}
	{{ printf "\n\t/* -- %s -- */" $res.Name }}
	{{- range $op := $res.Operations -}}
		{{- if eq "set" $op -}}
			{{ printf "\n\trpc Set%s(%s) returns (%s.Meta);" $res.Name $res.Name $res.Name }}
		{{- end -}}
		{{- if eq "get" $op -}}
			{{ printf "\n\trpc Get%s(%s.Meta) returns (%s);" $res.Name $res.Name $res.Name }}
		{{- end -}}
		{{- if eq "list" $op -}}
			{{ printf "\n\trpc List%s(%s.Meta) returns (%sList);" $res.Name $res.Name $res.Name }}
		{{- end -}}
		{{- if eq "delete" $op -}}
			{{ printf "\n\trpc Delete%s(%s.Meta) returns (%s);" $res.Name $res.Name $res.Name }}
		{{- end -}}
	{{- end -}}
{{- end }}
}
message Void {}
{{ range $res := .Resources }}
message {{ $res.Name }}List {
	repeated {{ $res.Name }}.Meta metas = 1;
}

message {{ $res.Name }} {{ $res.Definition }}
{{ end -}}
