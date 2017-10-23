package stats

import (
	"sync"
	"time"
)

// Stats holds stats on bulk processing
type Stats struct {
	mu                    sync.RWMutex
	StartTime             time.Time
	TotalExtensions       int
	TotalExtensionsClosed int
	TotalExtensionsFailed int
	TotalFiles            int
	TotalSize             uint64
}

// New constructs and returns a Stats pointer
func New() *Stats {

	stats := &Stats{
		StartTime:             time.Now(),
		TotalExtensions:       0,
		TotalExtensionsClosed: 0,
		TotalFiles:            0,
		TotalSize:             0,
	}

	return stats

}

func (s *Stats) GetTimeTaken() time.Duration {

	end := time.Now()

	s.mu.RLock()
	totalTime := end.Sub(s.StartTime).Round(time.Second)
	s.mu.RUnlock()

	return totalTime

}

func (s *Stats) IncrementTotalExtensions() {

	s.mu.Lock()
	s.TotalExtensions++
	s.mu.Unlock()

}

func (s *Stats) GetTotalExtensions() int {

	s.mu.RLock()
	totalExtensions := s.TotalExtensions
	s.mu.RUnlock()

	return totalExtensions

}

func (s *Stats) IncrementTotalExtensionsClosed() {

	s.mu.Lock()
	s.TotalExtensionsClosed++
	s.mu.Unlock()

}

func (s *Stats) GetTotalExtensionsClosed() int {

	s.mu.RLock()
	totalExtensionsClosed := s.TotalExtensionsClosed
	s.mu.RUnlock()

	return totalExtensionsClosed

}

func (s *Stats) IncrementTotalExtensionsFailed() {

	s.mu.Lock()
	s.TotalExtensionsFailed++
	s.mu.Unlock()

}

func (s *Stats) GetTotalExtensionsFailed() int {

	s.mu.RLock()
	totalExtensionsFailed := s.TotalExtensionsFailed
	s.mu.RUnlock()

	return totalExtensionsFailed

}

func (s *Stats) IncrementTotalFiles() {

	s.mu.Lock()
	s.TotalFiles++
	s.mu.Unlock()

}

func (s *Stats) GetTotalFiles() int {

	s.mu.RLock()
	totalFiles := s.TotalFiles
	s.mu.RUnlock()

	return totalFiles

}

func (s *Stats) IncreaseTotalSize(size uint64) {

	s.mu.Lock()
	s.TotalSize += size
	s.mu.Unlock()

}

func (s *Stats) GetTotalSize() uint64 {

	s.mu.RLock()
	totalSize := s.TotalSize
	s.mu.RUnlock()

	return totalSize

}
