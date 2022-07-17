package tpl

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/orzen/steve/pkg/conv"
)

type Cfg struct {
	AppName      string
	AppVersion   string
	SteveVersion string
	Backend      string
	Proto        *conv.Proto
	Funcs        template.FuncMap
}

func TemplateFromFile(outputFile, tplFile string, cfg *Cfg) error {
	tplName := filepath.Base(tplFile)
	tpl, err := template.New(tplName).Funcs(cfg.Funcs).ParseFiles(tplFile)
	if err != nil {
		return fmt.Errorf("read template: %v", err)
	}

	if err := RenderTemplate(outputFile, tpl, cfg); err != nil {
		return fmt.Errorf("render template: %v", err)
	}

	return nil
}

func RenderTemplate(outputFile string, tpl *template.Template, cfg *Cfg) error {
	fd, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("proto file: %v", err)
	}
	defer fd.Close()

	if err := tpl.Execute(fd, cfg); err != nil {
		os.Remove(fd.Name())
		return fmt.Errorf("execute template: %v", err)
	}

	return nil
}
