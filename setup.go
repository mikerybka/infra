package infra

// Setup sets up a new Hetzner Project with the config given.
// Each region is supplied with 2 VMs (1 for data, 1 for traffic), a public floating IP, and a private network connecting the two VMs.
func Setup()
