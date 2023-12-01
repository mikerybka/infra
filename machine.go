package infra

type Machine struct {
	IP     string
	Region string // "na" or "eu"
	Config *MachineConfig
}
