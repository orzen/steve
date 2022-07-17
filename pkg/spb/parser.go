package spb

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/orzen/steve/pkg/resource"
	"github.com/rs/zerolog/log"
)

var exprs map[string]*regexp.Regexp

func init() {
	exprs = make(map[string]*regexp.Regexp)
	exprCfg := []struct {
		Name string
		Expr string
	}{
		{
			Name: "id",
			Expr: `(?m)(?P<type>)\s+id\s*=\s*[0-9]+;$`,
		},
		{
			Name: "message",
			Expr: `(?m)^message\s+(?P<name>\w+)\s+(?P<blockStart>[{])`,
		},
		{
			Name: "meta",
			Expr: `(?m)message\s+(?P<name>Meta)\s+(?P<blockStart>[{])`,
		},
		{
			Name: "plural",
			Expr: `(?m)^plural\s+(?P<resource>[A-Z]?\w*)\s*=\s*"(?P<plural>[A-Z]?[a-z]+)"`,
		},
		{
			Name: "operations",
			Expr: `(?m)^operations\s+(?P<resource>[A-Z]?\w*)\s*=\s*"(?P<operations>[a-z,]+)"`,
		},
	}

	for _, e := range exprCfg {
		expr, err := regexp.Compile(e.Expr)
		if err != nil {
			log.Fatal().Err(err).Msgf("compile '%s' expr", e.Name)
		}
		exprs[e.Name] = expr
	}
}

// UNUSED
type Match struct {
	Groups map[string]Group
}

// UNUSED
type Group struct {
	Value string
	First int
	Last  int
}

// UNUSED
func MatchGroups(expr *regexp.Regexp, data string) ([]Match, error) {
	matches := expr.FindAllStringSubmatch(data, -1)
	indices := expr.FindAllStringSubmatchIndex(data, -1)
	res := make([]Match, 0)

	// Verify that there's actual groups before assuming that they exist
	if expr.NumSubexp() == 0 {
		return res, errors.New("no match")
	}

	// Remove first name since it's the entire match(name "") and we are
	// just interested in actual groups.
	groups := expr.SubexpNames()[1:]

	// Matches through out the entire data
	for i, m := range matches {
		match := &Match{
			Groups: make(map[string]Group),
		}
		matchIndices := indices[i]

		// Iterate groups found within an expression match
		for _, g := range groups {
			idx := expr.SubexpIndex(g)

			// Results in FindAllStringSubmatchIndex contain pairs, first
			// location is the match beginning and the second is the match
			// end + 1. In case of message regexp it will match n times at
			// writing times resEnd will have position 2 in the regexp
			// match that means that its locations in the
			// FindAllStringSubmatchIndex array will be index 3 and 4. We
			// multiple by 2 to skip the second value of the previous pair.
			// We need the index to tell the block parser where to start
			// parsing.
			firstIdx := idx * 2
			lastIdx := idx*2 + 1

			match.Groups[g] = Group{
				Value: m[idx],
				First: matchIndices[firstIdx],
				Last:  matchIndices[lastIdx],
			}
		}

		res = append(res, *match)
	}

	return res, nil
}

func ParseMessage(acc map[string]*resource.Resource, expr *regexp.Regexp, data string) error {
	res := expr.FindAllStringSubmatch(data, -1)
	pos := expr.FindAllStringSubmatchIndex(data, -1)

	n := expr.NumSubexp()
	if n == 0 {
		return errors.New("no message")
	}

	resName := expr.SubexpIndex("name")
	resBlockStart := expr.SubexpIndex("blockStart")
	for i, r := range res {
		p := pos[i]
		name := r[resName]
		if _, ok := acc[name]; ok {
			return fmt.Errorf("message '%s' exists", name)
		}

		// Results in FindAllStringSubmatchIndex contain pairs, first
		// location is the match beginning and the second is the match
		// end + 1. In case of message regexp it will match n times at
		// writing times resEnd will have position 2 in the regexp
		// match that means that its locations in the
		// FindAllStringSubmatchIndex array will be index 3 and 4. We
		// multiple by 2 to skip the second value of the previous pair.
		// We need the index to tell the block parser where to start
		// parsing.
		posIdx := resBlockStart * 2

		acc[name] = &resource.Resource{
			Message: &resource.Message{
				Name:       name,
				BlockStart: p[posIdx],
			},
		}
	}

	return nil
}

func ParsePlural(acc map[string]*resource.Resource, data string) error {
	expPlural := exprs["plural"]
	plurals := expPlural.FindAllStringSubmatch(data, -1)

	pluralName := expPlural.SubexpIndex("resource")
	pluralValue := expPlural.SubexpIndex("plural")
	for _, p := range plurals {
		name := p[pluralName]
		res, ok := acc[name]
		if !ok {
			return fmt.Errorf("plural for unspecified resource '%s'", name)
		}
		res.Plural = p[pluralValue]
	}

	return nil
}

func ParseOperations(acc map[string]*resource.Resource, data string) error {
	expOps := exprs["operations"]
	ops := expOps.FindAllStringSubmatch(data, -1)

	opsName := expOps.SubexpIndex("resource")
	opsValue := expOps.SubexpIndex("operations")
	for _, o := range ops {
		name := o[opsName]

		res, ok := acc[name]
		if !ok {
			return fmt.Errorf("operations using unspecified resource '%s'", name)
		}

		res.Operations = strings.Split(o[opsValue], ",")
	}

	return nil
}

func ParseBlock(acc map[string]*resource.Resource, data string) error {
	for k, v := range acc {
		nestDepth := 0
		buf := bytes.NewBufferString("")

		// Slice the data to get placed at the beginning of the
		// resource content
		for _, r := range data[v.BlockStart:] {
			buf.WriteRune(r)

			if r == '{' {
				nestDepth += 1
			}
			if r == '}' {
				nestDepth -= 1
				if nestDepth == 0 {
					break
				}
			}
		}

		if buf.Len() == 0 {
			return fmt.Errorf("no definition for resource '%s'", k)
		}

		v.Definition = buf.String()

	}

	return nil
}

func ParseMeta(acc map[string]*resource.Resource) error {
	expMeta := exprs["meta"]
	for k, v := range acc {
		nested := make(map[string]*resource.Resource)
		if err := ParseMessage(nested, expMeta, v.Definition); err != nil {
			return fmt.Errorf("parse meta for resource '%s'", k)
		}

		if err := ParseBlock(nested, v.Definition); err != nil {
			return fmt.Errorf("parse meta block: %v", err)
		}

		meta, ok := nested["Meta"]
		if !ok {
			return fmt.Errorf("meta is not definied for resource '%s'", k)
		}
		v.Meta = meta.Message
	}

	return nil
}

func Parse(acc map[string]*resource.Resource, data []byte) error {
	expMsg := exprs["message"]
	content := string(data)

	if err := ParseMessage(acc, expMsg, content); err != nil {
		return fmt.Errorf("parse message: %v", err)
	}

	if err := ParsePlural(acc, content); err != nil {
		return fmt.Errorf("parse plural: %v", err)
	}

	if err := ParseOperations(acc, content); err != nil {
		return fmt.Errorf("parse operations: %v", err)
	}

	if err := ParseBlock(acc, content); err != nil {
		return fmt.Errorf("parse block: %v", err)
	}

	if err := ParseMeta(acc); err != nil {
		return fmt.Errorf("parse meta: %v", err)
	}

	return nil
}
