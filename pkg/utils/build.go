package utils

func Compile(mainFile, binFile string) error {
	args := []string{"go", "build", "-o", binFile, mainFile}

	return Cmd(args)
}
