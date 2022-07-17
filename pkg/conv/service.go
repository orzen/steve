package conv

import (
	"fmt"
	"strings"
)

type Service struct {
	Name      string
	Resources []*Resource
}

func (s *Service) String() string {
	b := &strings.Builder{}

	// Build the service
	WriteString(b, fmt.Sprintf("service %s {\n", s.Name))
	for _, r := range s.Resources {
		WriteString(b, r.OperationString())
	}
	WriteString(b, "}\n\n")

	// Append messages
	for _, r := range s.Resources {
		WriteString(b, r.MessageString())
		WriteString(b, "\n\n")
	}

	return b.String()
}
