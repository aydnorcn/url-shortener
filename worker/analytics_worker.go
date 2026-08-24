package worker

import (
	"context"
	"log"
	"sync"
	"url-shortener/dto"
	"url-shortener/service"
)

type AnalyticsWorker interface {
	Start(ctx context.Context)
	Stop()
	Process(event dto.ClickEvent)
	Worker(id int)
}

type analyticsWorker struct {
	analyticsService service.AnalyticsService
	eventChan        chan dto.ClickEvent
	workerCount      int
	wg               sync.WaitGroup
	ctx              context.Context
	cancel           context.CancelFunc
	mu               sync.Mutex
	isStopped        bool
}

func NewAnalyticsWorker(analyticsService service.AnalyticsService, workerCount int, bufferSize int) AnalyticsWorker {
	if workerCount <= 0 {
		workerCount = 5
	}
	if bufferSize <= 0 {
		bufferSize = 1000
	}

	return &analyticsWorker{
		analyticsService: analyticsService,
		eventChan:        make(chan dto.ClickEvent, bufferSize),
		workerCount:      workerCount,
	}
}

func (w *analyticsWorker) Start(ctx context.Context) {
	w.ctx, w.cancel = context.WithCancel(ctx)
	for i := 1; i <= w.workerCount; i++ {
		w.wg.Add(1)
		go w.Worker(i)
	}
	log.Printf("Analytics worker pool started with %d workers", w.workerCount)
}

func (w *analyticsWorker) Worker(id int) {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			// Process remaining events in the channel upon shutdown
			for {
				select {
				case event, ok := <-w.eventChan:
					if !ok {
						return
					}
					if err := w.analyticsService.RecordClick(context.Background(), event); err != nil {
						log.Printf("[Worker %d] Error saving click event on shutdown for URL ID %d: %v", id, event.URLID, err)
					}
				default:
					return
				}
			}
		case event, ok := <-w.eventChan:
			if !ok {
				return
			}
			if err := w.analyticsService.RecordClick(context.Background(), event); err != nil {
				log.Printf("[Worker %d] Error recording click event for URL ID %d: %v", id, event.URLID, err)
			}
		}
	}
}

func (w *analyticsWorker) Process(event dto.ClickEvent) {
	w.mu.Lock()
	if w.isStopped {
		w.mu.Unlock()
		return
	}
	w.mu.Unlock()

	select {
	case w.eventChan <- event:
	default:
		log.Printf("Analytics worker queue full, dropping click event for URL ID %d", event.URLID)
	}
}

func (w *analyticsWorker) Stop() {
	w.mu.Lock()
	if w.isStopped {
		w.mu.Unlock()
		return
	}
	w.isStopped = true
	w.mu.Unlock()

	if w.cancel != nil {
		w.cancel()
	}
	close(w.eventChan)
	w.wg.Wait()
	log.Println("Analytics worker pool stopped")
}
