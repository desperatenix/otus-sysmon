package repos

import (
	"github.com/desperatenix/otus-sysmon/internal/cpu"
	"github.com/desperatenix/otus-sysmon/internal/loadavg"
	"time"
)

type MetricsData struct {
	La   *loadavg.Stats
	Cpu       *cpu.Stats
}

type TimePoint struct {
	Time time.Time
	MP *MetricsData
}

type Snapshots map[time.Time]*MetricsData

func (tp *TimePoint) Save(t time.Time, i interface{}) TimePoint {
	switch i.(type) {
	case loadavg.Stats:
		tp.MP.La = i.(*loadavg.Stats)
	case cpu.Stats:
		tp.MP.Cpu = i.(*cpu.Stats)
	}
	tp.Time = t

	return *tp
}
