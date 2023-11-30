package infra

import (
	"fmt"
	"os"

	"github.com/mikerybka/golang"
	"github.com/mikerybka/hetzner"
	"github.com/mikerybka/remote"
	"github.com/mikerybka/systemd"
	"github.com/mikerybka/util"
)

type Machine struct {
	IP     string
	Region string // "na" or "eu"
	Config *MachineConfig
}

type MachineConfig struct {
	Packages []*Package
	Binaries map[string]string
	Services map[string]*systemd.Service
	Jobs     []*Job
}

func (m *Machine) Deploy(htoken string) error {
	newMachine := m.IP == ""
	if newMachine {
		res, err := hetzner.CreateMachine(htoken, util.UnixTimestamp(), m.HetznerConfig())
		if err != nil {
			return err
		}
		m.IP = res.Server.PublicNet.IPv4.IP.String()

		err = remote.Run("root", m.IP, "wget https://builds.mikerybka.com/linux/x86/main/github.com/mikerybka/infra/cmd/provision /bin/provision")
		if err != nil {
			return err
		}
		err = remote.Run("root", m.IP, "chmod +x /bin/provsion")
		if err != nil {
			return err
		}
	}

	err := remote.WriteJSONFile("root", m.IP, "/etc/infra/config.json", m.Config)
	if err != nil {
		return err
	}
	err = remote.Run("root", m.IP, "/bin/provison")
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) HetznerConfig() *hetzner.MachineConfig {
	return hetzner.DefaultMachines[m.Region]
}

func (m *Machine) SetupLocally() error {
	err := m.InstallPackages()
	if err != nil {
		return err
	}
	err = m.InstallBinaries()
	if err != nil {
		return err
	}
	err = m.InstallServices()
	if err != nil {
		return err
	}
	err = m.WriteJobsFile()
	if err != nil {
		return err
	}
	err = m.InstallJobRunner()
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) InstallBinaries() error {
	for name, url := range m.Binaries {
		target := fmt.Sprintf("/bin/%s", name)
		err := util.Download(url, target)
		if err != nil {
			return err
		}
		err = os.Chmod(target, 751) // -rwxr-x--x
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Machine) WriteJobsFile() error {
	return util.WriteFile("/etc/jobs/jobs.json", util.Serialize(m.Jobs))
}

func (m *Machine) InstallJobRunner() error {
	return golang.Install("github.com/mikerybka/infra/cmd/job-runner")
}

func (m *Machine) InstallPackages() error {
	for _, p := range m.Packages {
		err := p.Install(m.HetznerConfig.OS)
		if err != nil {
			return fmt.Errorf("installing %s: %w", p.Name, err)
		}
	}
	return nil
}

func (m *Machine) InstallServices() error {
	for id, s := range m.Services {
		systemd.Stop(id)

		filename := fmt.Sprintf("/etc/systemd/system/%s.service", id)
		b := []byte(s.String())
		err := os.WriteFile(filename, b, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err := systemd.Reload()
	if err != nil {
		return err
	}
	for id := range m.Services {
		err := systemd.EnableNow(id)
		if err != nil {
			return err
		}
	}
	return nil
}
