package main

import "fmt"

type Bar struct {
	percent int
	cur     int
	tot     int
	rate    string
	graph   string
}

func (bar *Bar) NewOption(start, total int) {
	bar.cur = start
	bar.tot = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph
	}
}

func (bar *Bar) getPercent() int {
	return int(float32(bar.cur) / float32(bar.tot) * 100)
}

func (bar *Bar) Play(cur int) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent, bar.cur, bar.tot)
}

func (bar *Bar) Finish() {
	fmt.Println()
}
