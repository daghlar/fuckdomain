package progress

import (
	"fmt"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type Progress struct {
	bar        *pb.ProgressBar
	stats      *Stats
	mu         sync.RWMutex
	startTime  time.Time
	showStats  bool
}

type Stats struct {
	Total       int
	Completed   int
	Found       int
	Errors      int
	Rate        float64
	ETA         time.Duration
	Elapsed     time.Duration
}

func NewProgress(total int, showStats bool) *Progress {
	bar := pb.New(total)
	bar.SetTemplateString(`{{counters . }} {{bar . }} {{percent . }} {{speed . }} {{rtime . "ETA %s"}}`)
	bar.SetWidth(80)
	bar.SetMaxWidth(80)

	return &Progress{
		bar:       bar,
		stats:     &Stats{Total: total},
		startTime: time.Now(),
		showStats: showStats,
	}
}

func (p *Progress) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bar.Start()
}

func (p *Progress) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bar.Finish()
}

func (p *Progress) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bar.Increment()
	p.stats.Completed++
	p.updateStats()
}

func (p *Progress) AddFound() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stats.Found++
}

func (p *Progress) AddError() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stats.Errors++
}

func (p *Progress) SetTotal(total int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stats.Total = total
	p.bar.SetTotal(total)
}

func (p *Progress) updateStats() {
	now := time.Now()
	elapsed := now.Sub(p.startTime)
	p.stats.Elapsed = elapsed

	if p.stats.Completed > 0 {
		p.stats.Rate = float64(p.stats.Completed) / elapsed.Seconds()
		
		if p.stats.Rate > 0 {
			remaining := p.stats.Total - p.stats.Completed
			p.stats.ETA = time.Duration(float64(remaining)/p.stats.Rate) * time.Second
		}
	}
}

func (p *Progress) GetStats() Stats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return *p.stats
}

func (p *Progress) PrintStats() {
	if !p.showStats {
		return
	}

	p.mu.RLock()
	stats := *p.stats
	p.mu.RUnlock()

	fmt.Printf("\n")
	fmt.Printf("Statistics:\n")
	fmt.Printf("===========\n")
	fmt.Printf("Total:     %d\n", stats.Total)
	fmt.Printf("Completed: %d\n", stats.Completed)
	fmt.Printf("Found:     %d\n", stats.Found)
	fmt.Printf("Errors:    %d\n", stats.Errors)
	fmt.Printf("Rate:      %.2f req/s\n", stats.Rate)
	fmt.Printf("Elapsed:   %s\n", stats.Elapsed.Round(time.Second))
	fmt.Printf("ETA:       %s\n", stats.ETA.Round(time.Second))
	fmt.Printf("\n")
}

type MultiProgress struct {
	bars  map[string]*Progress
	mu    sync.RWMutex
	stats *MultiStats
}

type MultiStats struct {
	TotalBars   int
	ActiveBars  int
	TotalItems  int
	Completed   int
	Found       int
	Errors      int
	StartTime   time.Time
	Elapsed     time.Duration
}

func NewMultiProgress() *MultiProgress {
	return &MultiProgress{
		bars:  make(map[string]*Progress),
		stats: &MultiStats{StartTime: time.Now()},
	}
}

func (mp *MultiProgress) AddBar(name string, total int, showStats bool) *Progress {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	progress := NewProgress(total, showStats)
	mp.bars[name] = progress
	mp.stats.TotalBars++
	mp.stats.ActiveBars++
	mp.stats.TotalItems += total

	return progress
}

func (mp *MultiProgress) RemoveBar(name string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if bar, exists := mp.bars[name]; exists {
		bar.Stop()
		delete(mp.bars, name)
		mp.stats.ActiveBars--
	}
}

func (mp *MultiProgress) GetBar(name string) (*Progress, bool) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	bar, exists := mp.bars[name]
	return bar, exists
}

func (mp *MultiProgress) UpdateStats() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.stats.Completed = 0
	mp.stats.Found = 0
	mp.stats.Errors = 0

	for _, bar := range mp.bars {
		stats := bar.GetStats()
		mp.stats.Completed += stats.Completed
		mp.stats.Found += stats.Found
		mp.stats.Errors += stats.Errors
	}

	mp.stats.Elapsed = time.Since(mp.stats.StartTime)
}

func (mp *MultiProgress) PrintStats() {
	mp.UpdateStats()

	fmt.Printf("\n")
	fmt.Printf("Multi-Progress Statistics:\n")
	fmt.Printf("=========================\n")
	fmt.Printf("Total Bars:  %d\n", mp.stats.TotalBars)
	fmt.Printf("Active Bars: %d\n", mp.stats.ActiveBars)
	fmt.Printf("Total Items: %d\n", mp.stats.TotalItems)
	fmt.Printf("Completed:   %d\n", mp.stats.Completed)
	fmt.Printf("Found:       %d\n", mp.stats.Found)
	fmt.Printf("Errors:      %d\n", mp.stats.Errors)
	fmt.Printf("Elapsed:     %s\n", mp.stats.Elapsed.Round(time.Second))
	fmt.Printf("\n")
}

func (mp *MultiProgress) StopAll() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	for _, bar := range mp.bars {
		bar.Stop()
	}
	mp.bars = make(map[string]*Progress)
	mp.stats.ActiveBars = 0
}
