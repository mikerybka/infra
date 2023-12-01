package infra

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mikerybka/util"
)

func Build(path string, buildDir string) error {
	envs := []struct {
		OS   string
		Arch string
	}{
		{
			OS:   "darwin",
			Arch: "arm64",
		},
		{
			OS:   "js",
			Arch: "wasm",
		},
		{
			OS:   "linux",
			Arch: "386",
		},
		{
			OS:   "linux",
			Arch: "amd64",
		},
		{
			OS:   "linux",
			Arch: "arm",
		},
		{
			OS:   "linux",
			Arch: "arm64",
		},
		{
			OS:   "linux",
			Arch: "riscv64",
		},
		// {
		// 	OS:   "wasip1",
		// 	Arch: "wasm",
		// },
	}

	for _, env := range envs {
		o := filepath.Join(buildDir, env.OS, env.Arch, path)
		cmd := exec.Command("go", "build", "-o", o, path)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", env.OS), fmt.Sprintf("GOARCH=%s", env.Arch))
		cmd.Dir = filepath.Join(util.HomeDir(), "src")
		err := util.Run(cmd)
		if err != nil {
			return fmt.Errorf("building %s/%s: %w", env.OS, env.Arch, err)
		}
	}
	return nil
}
