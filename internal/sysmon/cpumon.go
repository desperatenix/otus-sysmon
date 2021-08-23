package sysmon2

import (
	"fmt"
	"github.com/desperatenix/otus-sysmon/internal/cpu"
	"github.com/desperatenix/otus-sysmon/internal/repos"
	"log"
	"sync"
)

func GetCpu(chStop chan bool, mu sync.RWMutex, ch <-chan repos.TimePoint) { //nolint:govet
	fmt.Println("debug CPU")
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