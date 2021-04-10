package main

import (
	"container/ring"

	lc "wesalvaro.com/nyandroidlogcat"
)

type Entring struct {
	start *ring.Ring
	end   *ring.Ring
}

func newEntring(logSize int) *Entring {
	r := ring.New(logSize)
	return &Entring{start: r, end: r}
}

func (r *Entring) ToList() []string {
	start := r.start
	end := r.end
	i := start
	var entries []string
	for {
		v, _ := i.Value.(*lc.Entry)
		entries = append(entries, v.TermUiString())
		i = i.Next()
		if i == end {
			break
		}
	}
	return entries
}

func (r *Entring) Append(e *lc.Entry) {
	if r.end == r.start && r.start.Value != nil {
		r.start = r.start.Next()
	}
	r.end.Value = e
	r.end = r.end.Next()
}
