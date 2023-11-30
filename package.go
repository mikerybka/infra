package infra

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mikerybka/util"
)

type Package struct {
	Name       string
	AptName    string
	DnfName    string
	PacmanName string
}

func (p *Package) Install(os string) error {
	if strings.HasPrefix(os, "debian") || strings.HasPrefix(os, "ubuntu") {
		cmd := exec.Command("apt", "install", "-y", p.AptName)
		return util.Run(cmd)
	} else if strings.HasPrefix(os, "fedora") || strings.HasPrefix(os, "rhel") {
		cmd := exec.Command("dnf", "install", "-y", p.DnfName)
		return util.Run(cmd)
	} else if strings.HasPrefix(os, "arch") || strings.HasPrefix(os, "manjaro") {
		cmd := exec.Command("pacman", "-Syu", p.PacmanName)
		return util.Run(cmd)
	} else {
		return fmt.Errorf("unknown os: %s", os)
	}
}
