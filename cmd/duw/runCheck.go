package main

import (
	"DMarinuks/disk-usage-warner/hdd"
)

type runCheck struct {
	Verbose   bool     `help:"Show results of check" default:"false" env:"HDDM_SHOW_RESULTS"`
	Paths     []string `name:"path" help:"Disks to check, if empty, all will be checked." type:"path" env:"HDDM_PATHS"`
	Threshold int      `help:"Used percentage at witch a warning email should be send"`
}

func (rc *runCheck) Run() error {
	return hdd.Check(rc.Verbose, rc.Paths, rc.Threshold)
}
