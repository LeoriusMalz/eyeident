package workers

import (
	"context"
	"eyeident/internal/db"
	"log"
	"sync"
	"time"
)

type Worker struct {
	Name      string
	LastRun   time.Time
	IsRunning bool
	Stopped   bool
	mu        sync.Mutex
}

func NewWorker(name string) *Worker {
	return &Worker{
		Name:    name,
		Stopped: false,
	}
}

func (w *Worker) Run() {
	w.mu.Lock()
	if w.IsRunning || w.Stopped {
		w.mu.Unlock()
		return
	}
	w.IsRunning = true
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		w.IsRunning = false
		w.LastRun = time.Now()
		w.mu.Unlock()
	}()

	log.Printf("Worker %s started\n", w.Name)

	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return
	}

	sqlMigration, _ := db.LoadQuery("migration.sql")

	_, err = tx.Exec(context.Background(), sqlMigration)
	if err != nil {
		log.Println("Failed to insert data:", err)
		tx.Rollback(context.Background())
		return
	}

	_, err = tx.Exec(context.Background(), `TRUNCATE raw_data`)
	if err != nil {
		log.Println("Failed to truncate raw_data:", err)
		tx.Rollback(context.Background())
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("Failed to commit transaction:", err)
		return
	}

	log.Printf("Worker %s finished successfully\n", w.Name)
}

func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Stopped = true
}

func (w *Worker) Resume() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Stopped = false
}

func (w *Worker) Status() map[string]interface{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	return map[string]interface{}{
		"name":       w.Name,
		"last_run":   w.LastRun,
		"is_running": w.IsRunning,
		"stopped":    w.Stopped,
	}
}
