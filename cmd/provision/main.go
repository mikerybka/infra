package main

import (
	"github.com/mikerybka/infra"
	"github.com/mikerybka/util"
)

func main() {
	c := &infra.MachineConfig{}
	err := util.ReadJSONFile("/etc/infra/config.json", c)
	if err != nil {
		panic(err)
	}
	err = c.ProvsionLocal()
	if err != nil {
		panic(err)
	}
}
