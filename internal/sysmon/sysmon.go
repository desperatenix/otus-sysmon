package sysmon2

// Todo  постараться разобраться из-за чего возникает гонка nolint:govet

import (
	"errors"
	"fmt"
	//"github.com/desperatenix/otus-sysmon/internal/sysmon"
	"log"
	"sync"
	"time"

	"github.com/desperatenix/otus-sysmon/internal/config"
	"github.com/desperatenix/otus-sysmon/internal/repos"
)

type Sysmon2 struct {
	Conf        *config.Config
	MetricsData *repos.MetricsData
	StopChan    chan bool
	workerCh    []chan<- repos.TimePoint
	mu          sync.RWMutex
	Snapshots   repos.Snapshots
}

func (sm *Sysmon2) Start() {
	fmt.Println("daemon start")

	sm.Snapshots = make(repos.Snapshots)

	go func() {
		if err := sm.EnabledMetrics(); err != nil {
			log.Fatal(err)
		}
	}()
	select {
	case <-sm.StopChan:
		fmt.Println("daemon stop")
		return
	default:
		go sm.work()
	}
}

func (sm *Sysmon2) Stop() {
	close(sm.StopChan)
}

var wg sync.WaitGroup

func (sm *Sysmon2) EnabledMetrics() error {
	var i int
	for metric, isEnable := range sm.Conf.Metric {
		if isEnable {
			switch metric {
			case "la":
				i++
				wg.Add(1)
				go GetLa(sm.StopChan, sm.mu, sm.newWorkerCh(), &wg) //nolint:govet

			case "cpu":
				i++
				wg.Add(1)
				go GetCpu(sm.StopChan, sm.mu, sm.newWorkerCh(), &wg) //nolint:govet

			default:
				fmt.Printf("Unknown metrics type(%s) for collection", metric)
				continue
			}
		}
	}
	wg.Done()
	if i < 1 {
		err := errors.New("no interrogators for metric collections")
		return err
	}
	return nil
}

func (sm *Sysmon2) addPoint(now time.Time) *repos.MetricsData {
	point := &repos.MetricsData{}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.Snapshots[now] = point

	// Todo Debug log
	fmt.Println(sm.Snapshots)

	return point
}

func (sm *Sysmon2) newWorkerCh() chan repos.TimePoint {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	ch := make(chan repos.TimePoint, 1)
	sm.workerCh = append(sm.workerCh, ch)

	return ch
}

func (sm *Sysmon2) processTick(now time.Time) {
	// добавляется новая точка для статистики за эту секунду

	point := sm.addPoint(now)

	// точка отправляется всем горутинам, ответственным за получение части статистики для заполнения
	for _, ch := range sm.workerCh {
		tp := repos.TimePoint{
			Time: now,
			MP:   point,
		}
		select {
		case ch <- tp:
		default:
		}
	}
}

func (sm *Sysmon2) work() {
	defer close(sm.StopChan)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-sm.StopChan:
			return
		default:
		}

		select {
		case <-sm.StopChan:
			return
		case now := <-ticker.C:
			sm.processTick(now.Truncate(time.Second))
		}
	}
}
