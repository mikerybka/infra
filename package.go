package infra

import (
	"fmt"
	"os/exec"

	"github.com/mikerybka/util"
)

type Package struct {
	Name         string
	HomebrewName string
	AptName      string
	DnfName      string
	PacmanName   string
}

func (p *Package) InstallLocal() error {
	pm := util.PackageManager()
	switch pm {
	case "homebrew":
		cmd := exec.Command("brew", "install", p.HomebrewName)
		return util.Run(cmd)
	case "apt":
		cmd := exec.Command("apt", "install", "-y", p.AptName)
		return util.Run(cmd)
	case "pacman":
		cmd := exec.Command("pacman", "-Syu", p.PacmanName)
		return util.Run(cmd)
	case "dnf":
		cmd := exec.Command("dnf", "install", "-y", p.DnfName)
		return util.Run(cmd)
	default:
		return fmt.Errorf("unknown package manager: %s", pm)
	}
}
