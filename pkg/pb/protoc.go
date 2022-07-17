package pb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orzen/steve/pkg/utils"
)

func Protoc(workDir, inputFile string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get wd: %v", err)
	}

	rel, err := filepath.Rel(wd, workDir)
	if err != nil {
		return fmt.Errorf("relative path: %v", err)
	}

	args := []string{"protoc",
		"--proto_path", rel,
		"--go_out", rel,
		"--go_opt=paths=source_relative",
		"--go-grpc_out", rel,
		"--go-grpc_opt=paths=source_relative",
		filepath.Base(inputFile)}

	return utils.Cmd(args)
}
