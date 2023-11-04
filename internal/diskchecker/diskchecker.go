package diskchecker

import (
	"fmt"
	"os"
	"strings"

	"github.com/DMarinuks/disk-usage-warner/internal/logger"
	"github.com/DMarinuks/disk-usage-warner/internal/messenger/types"

	humanize "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/disk"
	"go.uber.org/zap"
)

type permissionError struct {
	device string
	err    error
}

var defaultMessenger types.Messenger

func SetDefaultMessenger(messenger types.Messenger) {
	defaultMessenger = messenger
}

// Check - will check disk usage and send warning email
// if percentage was provided. Setting verbose as true
// will print disk usage.
func Check(verbose bool, paths []string, th int) error {
	log := logger.Named("disk")

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %w", err)
	}

	hostname = strings.ToLower(strings.TrimSpace(hostname))
	if len(hostname) == 0 {
		return fmt.Errorf("empty hostname is invalid")
	}

	formatter := "%-14s %7s %7s %7s %4s %s\n"
	if verbose {
		fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")
	}

	var multiError []permissionError
	var warningInfos []*types.WarningInfo

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
				humanize.Bytes(s.Total),
				humanize.Bytes(s.Used),
				humanize.Bytes(s.Free),
				percent,
				p.Mountpoint,
			)
		}

		if th != 0 && float64(th) <= s.UsedPercent {
			// collect threshold warnings
			warningInfos = append(warningInfos, &types.WarningInfo{
				Device:  device,
				Percent: percent,
			})
		}
	}

	if len(multiError) > 0 {
		for _, me := range multiError {
			log.Warn("error getting disk usage", zap.String("device", me.device), zap.Error(me.err))
		}
	}

	if len(warningInfos) > 0 {
		log.Debug("sending notification", zap.Any("warnings", warningInfos))

		if err := defaultMessenger.Send(hostname, warningInfos); err != nil {
			log.Error("error sending notification", zap.Error(err))
		}
	}

	return nil
}

func strInSlice(n string, s []string) bool {
	for _, ss := range s {
		if ss == n {
			return true
		}
	}
	return false
}
