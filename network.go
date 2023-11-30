package infra

type Network struct {
	Machines    map[string]*Machine
	FloatingIPs map[string]*FloatingIP
}
