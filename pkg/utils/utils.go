package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"

	"github.com/rs/zerolog/log"
)

func Type(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}

// TODO remove if remain unused.
// It's not used becuse /tmp and /home are located on different
// devices(partitions) which causes the rename to be non atomic. This makes
// this function pointless.
func WriteFile(path string, content string) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat '%s': %v", path, err)
	}

	if os.IsExist(err) {
		return fmt.Errorf("file '%s' exists", path)
	}

	temp, err := os.CreateTemp("", "steve-tmp-******")
	if err != nil {
		return fmt.Errorf("create tempfile: %v", err)
	}
	defer temp.Close()

	if _, err = temp.WriteString(content); err != nil {
		return fmt.Errorf("write content: %v", err)
	}

	if err := os.Rename(temp.Name(), path); err != nil {
		return fmt.Errorf("move file(safe write): %v", err)
	}

	return nil
}

func Cmd(args []string) error {
	exe := args[0]
	args = args[1:]

	log.Debug().Str("executable", exe).Interface("args", args).Msg("command config")

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
		fmt.Println("stderr", sErr.String())
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
