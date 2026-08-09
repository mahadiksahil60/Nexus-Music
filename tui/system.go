package tui

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type systemStats struct {
	cpuUsage float64
	ramUsage float64
}

func getSystemStats() systemStats {
	cpuUsage := 0.0

	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	ramUsage := 0.0

	memory, err := mem.VirtualMemory()
	if err == nil {
		ramUsage = memory.UsedPercent
	}

	return systemStats{
		cpuUsage: cpuUsage,
		ramUsage: ramUsage,
	}
}

func renderSystemStats() string {
	stats := getSystemStats()

	return fmt.Sprintf(
		"CPU %5.1f%%    RAM %5.1f%%",
		stats.cpuUsage,
		stats.ramUsage,
	)
}
