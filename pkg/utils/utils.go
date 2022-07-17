package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func Cmd(args []string) error {
	exe := args[0]
	args = args[1:]

	fmt.Println("exe", exe)
	fmt.Println("args", args)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get wd: %v", err)
	}

	c, err := exec.LookPath(exe)
	if err != nil {
		log.Error().Err(err).Msgf("locate executable '%s'", exe)
		return fmt.Errorf("locate executable '%s': %v", exe, err)
	}

	// Setup command
	cmd := exec.Command(c, args...)
	cmd.Dir = wd

	// Configure command environment
	env := os.Environ()
	home := os.Getenv("HOME")

	if os.Getenv("GOPATH") == "" {
		gopath := fmt.Sprintf("GOPATH=%s", filepath.Join(home, "go"))
		env = append(env, gopath)
	}
	if os.Getenv("GOBIN") == "" {
		gobin := fmt.Sprintf("GOBIN=%s", filepath.Join(home, "go", "bin"))
		env = append(env, gobin)
	}

	cmd.Env = env

	// Output handling
	var sOut bytes.Buffer
	var sErr bytes.Buffer

	cmd.Stdout = &sOut
	cmd.Stderr = &sErr

	if err := cmd.Run(); err != nil {
		log.Error().Str("stderr", sErr.String()).Msgf("run executable '%s'", exe)
		return fmt.Errorf("run executable '%s': %v", exe, err)
	}

	log.Debug().Str("stdout", sOut.String()).Msgf("run executable '%s'", exe)

	return nil
}

func RelPath(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get wd: %v", err)
	}

	rel, err := filepath.Rel(wd, path)
	if err != nil {
		return "", fmt.Errorf("relative path: %v", err)
	}

	return rel, err
}
