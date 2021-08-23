package sysmon2

import (
	"github.com/desperatenix/otus-sysmon/internal/loadavg"
	"github.com/desperatenix/otus-sysmon/internal/repos"
	"log"
	"sync"
)

//Todo
//const timeToGetMetric = 950 * time.Millisecond


func GetLa(chStop chan bool, mu sync.RWMutex, ch <-chan repos.TimePoint) { //nolint:govet
	for {
		select {
		case <-chStop:
			return
		case tp := <-ch:
			func() {
				la, err := loadavg.Get()
				if err != nil {
					log.Print(err)
				}
				mu.RLock()
				defer mu.RUnlock()

				tp.MP.La = la
				log.Printf("Load Avarage: %f", la)
			}()
		}
	}
}