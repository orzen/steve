package resource

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	fieldExpr         = regexp.MustCompile(`(?P<type>\w+)\s+(?P<name>\w+)\s*=\s*(?P<position>\d+)\s*;\s*$`)
	allowedOperations = []string{"set", "get", "list", "delete"}
)

type Field struct {
	Type     string
	Name     string
	Position int
}

type Message struct {
	BlockStart int
	Name       string
	Definition string
	fields     []Field
}

func (m *Message) Fields() ([]Field, error) {
	if len(m.fields) != 0 {
		return m.fields, nil
	}

	res := []Field{}

	fields := fieldExpr.FindAllStringSubmatch(m.Definition, -1)

	iType := fieldExpr.SubexpIndex("type")
	iName := fieldExpr.SubexpIndex("name")
	iPos := fieldExpr.SubexpIndex("position")

	for _, f := range fields {
		pos, err := strconv.Atoi(f[iPos])
		if err != nil {

		}

		res = append(res, Field{
			Type:     f[iType],
			Name:     f[iName],
			Position: pos,
		})
	}

	m.fields = res

	return m.fields, nil
}

type Resource struct {
	*Message

	Meta       *Message
	Plural     string
	Operations []string
}

func (r *Resource) String() string {
	return fmt.Sprintf(`resource '%s'
  meta '%v'
  operations '%s'
  plural '%s'
  definition '%s'`,
		r.Name, r.Meta, r.Operations, r.Plural, r.Definition)
}

func (r *Resource) Finalize() error {
	if r.Name == "Meta" {
		return fmt.Errorf("'%s' is reserved name", r.Name)
	}

	// Validate operations
	for _, o := range r.Operations {
		allowed := false

		for _, a := range allowedOperations {
			if o == a {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Errorf("invalid operation '%s' for resource '%s'", o, r.Name)
		}
	}

	if len(r.Plural) == 0 {
		r.Plural = fmt.Sprintf("%ss", r.Name)
	}

	return nil
}
