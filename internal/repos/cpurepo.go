package repos

import (
	"github.com/desperatenix/otus-sysmon/internal/cpu"
	"time"
)

// CPUStats represents load average values
type CPUStats struct {
	Total  float64
	User   float64
	System float64
	Idle   float64
}

type CPURepo struct {
	SP []*cpu.Stats
}

func (sr *CPURepo) Save(timestamp time.Time, metric *cpu.Stats)  {
	sr.SP = append(sr.SP, metric)
}


func (sr *CPURepo) Get() []*cpu.Stats {
	return sr.SP
}