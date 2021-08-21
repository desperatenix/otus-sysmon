package loadavg

import (
	"fmt"
	"github.com/desperatenix/otus-sysmon/internal/common"
	"strconv"
	"strings"
)

// Read loadavg from /proc.
func get() (*Stats, error) {
	data, err := common.ReadProcFile("loadavg")
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	loads, err := parseLoad(data)
	if err != nil {
		return nil, err
	}
	return &Stats{
		Load1: loads[0],
		Load5: loads[1],
		Load15: loads[2],
		}, nil
}

// Parse /proc loadavg and return 1m, 5m and 15m.
func parseLoad(data []string) (loads []float64, err error) {
	loads = make([]float64, 3)
	parts := strings.Fields(data[0])
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected content in %s", "/proc/loadavg")
	}
	for i, load := range parts[0:3] {
		loads[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse load '%s': %w", load, err)
		}
	}
	return loads, nil
}