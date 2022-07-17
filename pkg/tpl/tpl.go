package tpl

import (
	"fmt"
	"html/template"
	"os"

	"github.com/orzen/steve/pkg/resource"
)

func TemplateFromFile(outputFile, tplFile string, resources map[string]*resource.Resource) error {
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return fmt.Errorf("read template: %v", err)
	}

	if err := RenderTemplate(outputFile, tpl, resources); err != nil {
		return fmt.Errorf("render template: %v", err)
	}

	return nil
}

// TODO remove if unused
//func TemplateFromString(outputDir, tpl string, resources map[string]*resource.Resource) (string, error) {
//	t, err := template.New("resources").Parse(tpl)
//	if err != nil {
//		return "", fmt.Errorf("parse template: %v", err)
//	}
//
//	dst, err := RenderTemplate(outputDir, t, resources)
//	if err != nil {
//		return "", fmt.Errorf("render template: %v", err)
//	}
//
//	return dst, nil
//}

func RenderTemplate(outputFile string, tpl *template.Template, resources map[string]*resource.Resource) error {
	// TODO write to tempfile and rename
	fd, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("proto file: %v", err)
	}
	defer fd.Close()

	res := []*resource.Resource{}
	for _, v := range resources {
		res = append(res, v)
	}

	in := struct {
		Resources []*resource.Resource
	}{
		Resources: res,
	}

	if err := tpl.Execute(fd, in); err != nil {
		os.Remove(fd.Name())
		return fmt.Errorf("execute template: %v", err)
	}

	return nil
}
