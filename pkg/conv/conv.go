package conv

var (
	// https://developers.google.com/protocol-buffers/docs/proto3#scalar
	PbGoTypes = map[string]string{
		"double":   "float64",
		"float":    "float32",
		"int32":    "int32",
		"int64":    "int64",
		"uint32":   "uint32",
		"uint64":   "uint64",
		"sint32":   "int32",
		"sint64":   "int64",
		"fixed32":  "uint32",
		"fixed64":  "uint64",
		"sfixed32": "int32",
		"sfixed64": "int64",
		"bool":     "bool",
		"string":   "string",
		"bytes":    "[]byte",
	}
)

func PbToGoType(t string) string {
	v, ok := PbGoTypes[t]
	if !ok {
		return "Name"
	}
	return v
}
