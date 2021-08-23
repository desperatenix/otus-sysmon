package sysmon2

import (
	"log"
	"sync"

	"github.com/desperatenix/otus-sysmon/internal/cpu"
	"github.com/desperatenix/otus-sysmon/internal/repos"
)

func GetCPU(chStop chan bool, mu sync.RWMutex, ch <-chan repos.TimePoint, wg *sync.WaitGroup) { //nolint:govet
	defer wg.Done()
	for {
		select {
		case <-chStop:
			return
		case tp := <-ch:
			func() {
				cpu, err := cpu.Get()
				if err != nil {
					log.Print(err)
				}
				mu.RLock()
				defer mu.RUnlock()

				tp.MP.Cpu = cpu
				log.Printf("CPU USage: %f", cpu)
			}()
		}
	}
}
