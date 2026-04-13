package scraper

import (
	"context"
	"log"
	"sync"

	"jobscraper/common"
)

type Scraper interface {
	Name() string
	Fetch(ctx context.Context) ([]*common.JobPayload, error)
}

type SyncResult struct {
	Company string `json:"company"`
	Status  string `json:"status"`
	Count   int    `json:"count,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Pool struct {
	scrapers    []Scraper
	ashbyLimit  int
	globalLimit int
}

func NewPool(scrapers []Scraper, ashbyLimit, globalLimit int) *Pool {
	return &Pool{
		scrapers:    scrapers,
		ashbyLimit:  ashbyLimit,
		globalLimit: globalLimit,
	}
}

func (p *Pool) Run(ctx context.Context) ([]SyncResult, error) {
	results := make([]SyncResult, 0, len(p.scrapers))
	var mu sync.Mutex
	var wg sync.WaitGroup
	globalSem := make(chan struct{}, p.globalLimit)
	ashbySem := make(chan struct{}, p.ashbyLimit)

	for _, s := range p.scrapers {
		wg.Add(1)
		go func(scraper Scraper) {
			defer wg.Done()
			globalSem <- struct{}{}
			defer func() { <-globalSem }()

			var sem chan struct{}
			if isAshby(scraper.Name()) {
				sem = ashbySem
			} else {
				sem = nil
			}

			if sem != nil {
				sem <- struct{}{}
				defer func() { <-sem }()
			}

			log.Printf("Starting sync - %s", scraper.Name())
			jobs, err := scraper.Fetch(ctx)

			mu.Lock()
			result := SyncResult{Company: scraper.Name()}
			if err != nil {
				result.Status = "failed"
				result.Error = "fetch failed"
				log.Printf("%s: FAILED - %v", scraper.Name(), err)
			} else {
				result.Status = "success"
				result.Count = len(jobs)
				log.Printf("%s: SUCCESS - %d jobs", scraper.Name(), len(jobs))
			}
			results = append(results, result)
			mu.Unlock()
		}(s)
	}

	wg.Wait()
	return results, nil
}

func isAshby(name string) bool {
	return name != "Amazon" && name != "Atlassian" && name != "Lever"
}
