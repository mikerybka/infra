package infra

import (
	"fmt"
	"os/exec"
	"strings"
)

type Package struct {
	Name       string
	AptName    string
	DnfName    string
	PacmanName string
}

func (p *Package) InstallCmd(os string) *exec.Cmd {
	if strings.HasPrefix(os, "debian") || strings.HasPrefix(os, "ubuntu") {
		return exec.Command("apt", "install", "-y", p.AptName)
	} else if strings.HasPrefix(os, "fedora") || strings.HasPrefix(os, "rhel") {
		return exec.Command("dnf", "install", "-y", p.DnfName)
	} else if strings.HasPrefix(os, "arch") || strings.HasPrefix(os, "manjaro") {
		return exec.Command("pacman", "-Syu", p.PacmanName)
	} else {
		panic(fmt.Errorf("unknown os: %s", os))
	}
}
