package main

import (
	disk "github.com/DMarinuks/disk-usage-warner/disk"
)

type runCheck struct {
	Verbose    bool     `help:"Show results of check" default:"false" env:"DUW_VERBOSE" short:"v"`
	Paths      []string `name:"path" help:"Disks to check, if empty, all will be checked." type:"path" env:"DUW_PATHS"`
	Percentage int      `help:"Used percentage at witch a warning email should be send" short:"p"`
}

func (rc *runCheck) Run() error {
	if !rc.Verbose && rc.Percentage == 0 {
		rc.Verbose = true
	}
	return disk.Check(rc.Verbose, rc.Paths, rc.Percentage)
}
