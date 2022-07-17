package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Doc struct {
	Data    interface{}
	Message string
}

type Console struct {
	Err      error
	JSON     bool
	Pretty   bool
	Status   int
	JSONFunc func(interface{}) ([]byte, error)

	doc *Doc
}

func New() *Console {
	return &Console{
		doc: &Doc{},
	}
}

func (c *Console) WithData(in interface{}) {
	c.doc.Data = in
}

func (c *Console) WithPretty() {
	c.Pretty = true
}

func (c *Console) WithErr(err error) {
	c.Err = err
	c.Status = 1
}

func (c *Console) AsJSON() {
	c.JSON = true
}

func interfaceToStr(in interface{}, indent string, level int) string {
	v := reflect.ValueOf(in)
	var buf strings.Builder

	// Deref pointer
	if v.Type().Kind() == reflect.Pointer {
		v = v.Elem()
	}

	// Handle struct
	if v.Type().Kind() == reflect.Struct {
		n := v.NumField()
		for i := 0; i < n; i++ {
			f := v.Field(i)
			str := interfaceToStr(f.Interface(), indent, level+1)
			buf.WriteString(fmt.Sprintf("%0*s%s", level, indent, str))
		}
		return buf.String()
	}

	// Handle array
	if v.Type().Kind() == reflect.Slice ||
		v.Type().Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			str := interfaceToStr(f.Interface(), indent, level+1)
			buf.WriteString(fmt.Sprintf("%0*s%s", level, indent, str))
		}
		return buf.String()
	}

	// Handle map
	if v.Type().Kind() == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()

			kStr := interfaceToStr(k.Interface(), indent, level+1)
			vStr := interfaceToStr(v.Interface(), indent, level+1)

			str := fmt.Sprintf("%0*s%s: %s", level, indent, kStr, vStr)
			buf.WriteString(str)
		}
		return buf.String()
	}

	// Handle everything else
	return v.String()
}

func (c *Console) pretty() string {
	data := interfaceToStr(c.doc.Data, "\t", 0)
	return fmt.Sprintf("data:\n%s\nmessage: %s\n", data, c.doc.Message)
}

func (c *Console) json() string {
	data, err := json.Marshal(c.doc)
	if err != nil {
		c.Err = err
	}
	return string(data)
}

func (c *Console) jsonPretty() string {
	data, err := json.MarshalIndent(c.doc, "", "\t")
	if err != nil {
		c.Err = err
	}
	return string(data)
}

func (c *Console) error() string {
	return fmt.Sprintf("error: %s", c.Err.Error())
}

func (c *Console) errorJSON() string {
	return fmt.Sprintf(`{"error": "%s"}`, c.Err.Error())
}

func (c *Console) printer() func() string {
	if c.JSON {
		if c.Err != nil {
			return c.errorJSON
		}
		if c.Pretty {
			return c.jsonPretty
		}
		return c.json
	}

	if c.Err != nil {
		return c.error
	}

	return c.pretty
}

func (c *Console) Print(msg string) {
	c.doc.Message = msg
	c.printer()()

	os.Exit(c.Status)
}
