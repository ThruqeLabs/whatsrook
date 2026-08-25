package cliutils

import (
	"whatsrook/cli/utils/system"
)

// CPUStats stores raw CPU usage counters from /proc/stat.
type CPUStats = system.CPUStats

// SystemMemory stores memory usage values from /proc/meminfo.
type SystemMemory = system.SystemMemory

var (
	// GetCPUModel retrieves the CPU model name from /proc/cpuinfo or returns a fallback.
	GetCPUModel = system.GetCPUModel

	// GetLoadAvg returns the 1, 5, and 15 minute load averages from /proc/loadavg.
	GetLoadAvg = system.GetLoadAvg

	// GetCPUStats reads and parses current CPU counters from /proc/stat.
	GetCPUStats = system.GetCPUStats

	// GetCPUUsage calculates active CPU utilization percentage.
	GetCPUUsage = system.GetCPUUsage

	// GetSystemMemory reads MemTotal, MemFree, and MemAvailable from /proc/meminfo.
	GetSystemMemory = system.GetSystemMemory

	// FormatBytes formats byte counts into human readable decimal string representation.
	FormatBytes = system.FormatBytes

	// FormatDuration formats duration into compact d/h/m/s components.
	FormatDuration = system.FormatDuration

	// MenuRuntime formats seconds into user-friendly runtime string.
	MenuRuntime = system.MenuRuntime

	// GetNumCPU returns the number of logical CPUs usable by the current process.
	GetNumCPU = system.GetNumCPU
)
