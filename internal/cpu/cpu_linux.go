package cpu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/desperatenix/otus-sysmon/internal/common"
)

func get() (*Stats, error) {
	content, err := common.ReadProcFile("stat")
	if err != nil {
		return nil, fmt.Errorf("cannot read the stat file: %w", err)
	}

	return parseCPU(content)
}

func parseCPU(content []string) (*Stats, error) {
	values := strings.Fields(content[0])[1:]

	user, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse user field: %w", err)
	}
	nice, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse nice field: %w", err)
	}
	system, err := strconv.ParseInt(values[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse system field: %w", err)
	}
	idle, err := strconv.ParseInt(values[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse idle field: %w", err)
	}
	iowait, err := strconv.ParseInt(values[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse iowait field: %w", err)
	}
	irq, err := strconv.ParseInt(values[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse irq field: %w", err)
	}
	softirq, err := strconv.ParseInt(values[6], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse softirq field: %w", err)
	}
	steal, err := strconv.ParseInt(values[7], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse steal field: %w", err)
	}

	total := float64(user + nice + system + idle + iowait + irq + softirq + steal)

	return &Stats{
		Total:  total,
		User:   float64(user) / total * 100,
		System: float64(system) / total * 100,
		Idle:   float64(idle) / total * 100,
	}, nil
}
