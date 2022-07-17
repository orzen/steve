package spb

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/orzen/steve/pkg/resource"
)

func LoadDir(resourceDir string) (map[string]*resource.Resource, error) {
	var f filepath.WalkFunc

	resources := make(map[string]*resource.Resource)

	f = func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".spb") && !info.IsDir() {
			if err := LoadFile(resources, path); err != nil {
				return fmt.Errorf("load file: %v", err)
			}
		}

		return nil
	}

	if err := filepath.Walk(resourceDir, f); err != nil {
		return resources, err
	}

	return resources, nil
}

func LoadFile(resources map[string]*resource.Resource, resourceFile string) error {
	data, err := os.ReadFile(resourceFile)
	if err != nil {
		return errors.New("read resource file")
	}

	return Parse(resources, data)
}
