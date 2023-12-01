package infra

import (
	"github.com/mikerybka/hetzner"
	"github.com/mikerybka/remote"
	"github.com/mikerybka/util"
)

type Machine struct {
	IP     string
	Region string // "na" or "eu"
	Config *MachineConfig
}

func (m *Machine) Deploy(htoken string) error {
	newMachine := m.IP == ""
	if newMachine {
		res, err := hetzner.CreateMachine(htoken, util.UnixTimestamp(), m.HetznerConfig())
		if err != nil {
			return err
		}
		m.IP = res.Server.PublicNet.IPv4.IP.String()

		err = remote.Run("root", m.IP, "wget https://builds.mikerybka.com/main/linux/x86/github.com/mikerybka/infra/cmd/provision /bin/provision")
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
