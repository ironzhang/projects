package collect

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

type entry struct {
	name  string
	count int
	sum   time.Duration
	min   time.Duration
	max   time.Duration
}

type stat struct {
	entries map[string]*entry
}

func (s *stat) Add(name string, d time.Duration) {
	if s.entries == nil {
		s.entries = make(map[string]*entry)
	}

	e, ok := s.entries[name]
	if !ok {
		e = &entry{
			name:  name,
			count: 1,
			sum:   d,
			min:   d,
			max:   d,
		}
		s.entries[name] = e
		return
	}

	e.count++
	e.sum += d
	if d < e.min {
		e.min = d
	}
	if d > e.max {
		e.max = d
	}
	return
}

func (s *stat) Clear() {
	s.entries = nil
}

func (s *stat) Print(w io.Writer) {
	if len(s.entries) <= 0 {
		return
	}

	entries := make([]entry, 0, len(s.entries))
	for _, e := range s.entries {
		entries = append(entries, *e)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].sum > entries[j].sum })

	cw := csv.NewWriter(w)
	cw.Write([]string{"time", "name", "count", "sum", "ave", "min", "max"})
	now := time.Now().Format("2006-01-02T15:04:05")
	for _, e := range entries {
		ave := e.sum / time.Duration(e.count)
		record := []string{
			now,
			e.name,
			fmt.Sprint(e.count),
			e.sum.String(),
			ave.String(),
			e.min.String(),
			e.max.String(),
		}
		cw.Write(record)
	}
	cw.Flush()
}

var Default Collect

type Collect struct {
	mu sync.Mutex
	st stat
}

func (p *Collect) Collect(name string) func() {
	start := time.Now()
	return func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.st.Add(name, time.Since(start))
	}
}

func (p *Collect) Add(name string, d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.st.Add(name, d)
}

func (p *Collect) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.st.Clear()
}

func (p *Collect) Print(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.st.Print(w)
}

func PrintDefaultCollect(w io.Writer, tick time.Duration, clear bool) {
	t := time.NewTicker(tick)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			Default.Print(w)
			if clear {
				Default.Clear()
			}
		}
	}
}
