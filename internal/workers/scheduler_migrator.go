package workers

import (
	"log"
	"time"
)

type Scheduler struct {
	Worker   *Worker
	Interval time.Duration
	stopChan chan bool
}

func NewScheduler(worker *Worker, interval time.Duration) *Scheduler {
	return &Scheduler{
		Worker:   worker,
		Interval: interval,
		stopChan: make(chan bool),
	}
}

func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(s.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Worker.Run()
			case <-s.stopChan:
				log.Println("Scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.stopChan <- true
}
