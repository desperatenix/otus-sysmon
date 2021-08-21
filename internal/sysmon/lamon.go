package sysmon2

import (
	"fmt"
	"github.com/desperatenix/otus-sysmon/internal/loadavg"
	"github.com/desperatenix/otus-sysmon/internal/repos"
	"log"
	"sync"
	"time"
)

//Todo

const timeToGetMetric = 950 * time.Millisecond


func GetLa(chStop chan bool, mu sync.Mutex, ch <-chan repos.TimePoint) {
	fmt.Println("debug")
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
				mu.Lock()
				defer mu.Unlock()

				tp.MP.La = la
				log.Printf("Load Avarage: %f", la)
			}()
		default:
			//nffmt.Println("time for sleep")
			//time.Sleep(100 * time.Millisecond)
		}
	}
}