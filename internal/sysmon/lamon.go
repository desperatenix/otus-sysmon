package sysmon2

import (
	"log"
	"sync"

	"github.com/desperatenix/otus-sysmon/internal/loadavg"
	"github.com/desperatenix/otus-sysmon/internal/repos"
)

func GetLa(chStop chan bool, mu sync.RWMutex, ch <-chan repos.TimePoint, wg *sync.WaitGroup) { //nolint:govet
	defer wg.Done()
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
