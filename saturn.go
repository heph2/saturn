package main

import (
	"flag"
	"fmt"
	"os"
)

func runFetch(input *string) (index map[int]string) {
	// Print on stdout the list of the episodes available
	// Then match the episode to an index
	episodes := FetchAnime(*input)

	index = make(map[int]string)
	for i, ep := range episodes {
		fmt.Printf("ID:%d \t %s\n", i, ep)
		index[i] = ep
	}
	return index
}

func runDown(input *int, id map[int]string) bool {
	if *input == -1 {
		return false
	}
	var epDownload []string

	epDownload = append(epDownload, id[*input])

	Pool(epDownload)

	return true
}

func main() {
	inputPtr := flag.String("fetch", "", "User input.")
	idPtr := flag.Int("download", -1, "Episode to download.")
	flag.Parse()

	if *inputPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	id := runFetch(inputPtr)

	ok := runDown(idPtr, id)
	if ok {
		fmt.Print("Download Started")
	}
}
