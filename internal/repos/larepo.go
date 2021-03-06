package repos

import (
	"time"

	"github.com/desperatenix/otus-sysmon/internal/loadavg"
)

type LARepo struct {
	SP []*loadavg.Stats
}

func (sr *LARepo) Save(timestamp time.Time, metric *loadavg.Stats) {
	// todo удаление старых данных
	sr.SP = append(sr.SP, metric)
}

func (sr *LARepo) Get() []*loadavg.Stats {
	return sr.SP
}
