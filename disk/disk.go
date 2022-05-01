package hdd

import (
	"fmt"

	"github.com/DMarinuks/disk-usage-warner/logger"
	"github.com/DMarinuks/disk-usage-warner/mailer"

	human "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/disk"
	"go.uber.org/zap"
)

func strInSlice(n string, s []string) bool {
	for _, ss := range s {
		if ss == n {
			return true
		}
	}
	return false
}

type permissionError struct {
	device string
	err    error
}

// Check - will check disk usage and send warning email
// if percentage was provided. Setting verbose as true
// will print disk usage.
func Check(verbose bool, paths []string, th int) error {
	log := logger.Named("disk")

	formatter := "%-14s %7s %7s %7s %4s %s\n"
	if verbose {
		fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")
	}

	var multiError []permissionError
	var warningInfos []*mailer.WarningInfo

	parts, _ := disk.Partitions(true)
	for _, p := range parts {
		device := p.Mountpoint
		// if paths were specified and path not in slice, skip
		if len(paths) != 0 && !strInSlice(device, paths) {
			continue
		}

		s, err := disk.Usage(device)
		if err != nil {
			// ignore permission errors
			if err.Error() == "permission denied" {
				multiError = append(multiError, permissionError{
					device: device,
					err:    err,
				})
				continue
			}
			log.Error("error getting disk usage", zap.String("device", device), zap.Error(err))
			return err
		}

		if s.Total == 0 {
			continue
		}

		// log.Debug("percentage", zap.String("device", device), zap.Float64("percent", s.UsedPercent))
		percent := fmt.Sprintf("%2.f%%", s.UsedPercent)
		if verbose {
			fmt.Printf(formatter,
				s.Fstype,
				human.Bytes(s.Total),
				human.Bytes(s.Used),
				human.Bytes(s.Free),
				percent,
				p.Mountpoint,
			)
		}

		if th != 0 && float64(th) <= s.UsedPercent {
			// collect threshold warnings
			warningInfos = append(warningInfos, &mailer.WarningInfo{
				Device:  device,
				Percent: percent,
			})
		}
	}

	if len(multiError) > 0 {
		// warn about errors
		for _, me := range multiError {
			log.Warn("error getting disk usage", zap.String("device", me.device), zap.Error(me.err))
		}
	}

	if len(warningInfos) > 0 {
		// send email
		if err := mailer.SendMail(warningInfos); err != nil {
			return err
		}
	}

	return nil
}
