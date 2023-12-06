package infra

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mikerybka/hetzner"
	"github.com/mikerybka/remote"
	"github.com/mikerybka/systemd"
	"github.com/mikerybka/util"
)

var DefaultMachine = &MachineConfig{}

type MachineConfig struct {
	Packages []*Package
	Binaries map[string]string
	Jobs     []*Job
	Services map[string]*systemd.Service
}

func (c *MachineConfig) Deploy(htoken, region string) (*Machine, error) {
	hmc := hetzner.DefaultMachines[region]
	machineName := util.UnixTimestamp()
	res, err := hetzner.CreateMachine(htoken, machineName, hmc)
	if err != nil {
		return nil, err
	}
	ip := res.Server.PublicNet.IPv4.IP.String()
	err = remote.Download("root", ip, "/bin/provision", "https://builds.mikerybka.com/main/linux/amd64/github.com/mikerybka/infra/cmd/provision")
	if err != nil {
		return nil, err
	}
	err = remote.Run("root", ip, "chmod +x /bin/provsion")
	if err != nil {
		return nil, err
	}
	err = remote.WriteJSONFile("root", ip, "/etc/infra/config.json", c)
	if err != nil {
		return nil, err
	}
	err = remote.Run("root", ip, "/bin/provison")
	if err != nil {
		return nil, err
	}
	return &Machine{
		IP:     ip,
		Region: region,
		Config: c,
	}, nil
}

func (m *MachineConfig) ProvsionLocal() error {
	err := m.InstallPackagesLocal()
	if err != nil {
		return err
	}
	err = m.InstallBinariesLocal()
	if err != nil {
		return err
	}
	err = m.WriteJobsFileLocal()
	if err != nil {
		return err
	}
	err = m.InstallJobRunnerLocal()
	if err != nil {
		return err
	}
	err = m.InstallServicesLocal()
	if err != nil {
		return err
	}
	return nil
}

func (m *MachineConfig) InstallBinariesLocal() error {
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

func (m *MachineConfig) WriteJobsFileLocal() error {
	return util.WriteFile("/etc/jobs/jobs.json", util.Serialize(m.Jobs))
}

func (m *MachineConfig) InstallJobRunnerLocal() error {
	url := fmt.Sprintf("https://builds.mikerybka.com/main/%s/%s/github.com/mikerybka/infra/cmd/job-runner", runtime.GOOS, runtime.GOARCH)
	err := util.Download(url, "/bin/job-runner")
	if err != nil {
		return err
	}
	err = os.Chmod("/bin/job-runner", 751)
	if err != nil {
		return err
	}
	return nil
}

func (m *MachineConfig) InstallPackagesLocal() error {
	for _, p := range m.Packages {
		err := p.InstallLocal()
		if err != nil {
			return fmt.Errorf("installing %s: %w", p.Name, err)
		}
	}
	return nil
}

func (m *MachineConfig) InstallServicesLocal() error {
	for id, s := range m.Services {
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
