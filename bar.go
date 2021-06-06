package main

import (
	"fmt"
	"sync"
)

type Bar struct {
	total      int
	downloaded int
	epsFetched int
	epsTotal   int
	m          sync.Mutex
}

func NewBar() *Bar {
	return &Bar{}
}

func (b *Bar) print() {
	var (
		fetched = byteConv(b.downloaded)
		total   = byteConv(b.total)
	)

	fmt.Printf(
		"\rDownloading %.2fMB of %.2fMB (%.2f%%) -- done %d of %d  \r",
		fetched,
		total,
		fetched*100.0/total,
		b.epsFetched,
		b.epsTotal,
	)
}

func (b *Bar) AddEpisode(size int) {
	b.m.Lock()
	defer b.m.Unlock()

	b.epsTotal++

	b.total += size

	b.print()
}

func (b *Bar) DoneEpisode() {
	b.m.Lock()
	defer b.m.Unlock()

	b.epsFetched++

	b.print()
}

func (b *Bar) Fetched(bytes int) {
	b.m.Lock()
	defer b.m.Unlock()

	b.downloaded += bytes

	b.print()
}
