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
	DnfRepo      string
	DnfName      string
	PacmanName   string
	PostInstall  string
}

func (p *Package) InstallLocal() error {
	pm := util.PackageManager()
	switch pm {
	case "homebrew":
		cmd := exec.Command("brew", "install", p.HomebrewName)
		err := util.Run(cmd)
		if err != nil {
			return err
		}
	case "apt":
		cmd := exec.Command("apt", "install", "-y", p.AptName)
		err := util.Run(cmd)
		if err != nil {
			return err
		}
	case "pacman":
		cmd := exec.Command("pacman", "-Syu", p.PacmanName)
		err := util.Run(cmd)
		if err != nil {
			return err
		}
	case "dnf":
		if p.DnfRepo != "" {
			cmd := exec.Command("dnf", "config-manager", "--add-repo", p.DnfRepo)
			err := util.Run(cmd)
			if err != nil {
				return err
			}
		}
		cmd := exec.Command("dnf", "install", "-y", p.DnfName)
		err := util.Run(cmd)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown package manager: %s", pm)
	}
	if p.PostInstall != "" {
		cmd := exec.Command("bash", "-c", p.PostInstall)
		err := util.Run(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
