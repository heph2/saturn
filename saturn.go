package main

import (
	"flag"
	"fmt"
	"os"
)

func flagsFetch() *string {
	inputPtr := flag.String("fetch", "", "User input.")

	return inputPtr
}

func flagsDown() *int {
	idPtr := flag.Int("download", -1, "Episode to download.")

	return idPtr
}

func main() {

	inputPtr := flagsFetch()
	idPtr := flagsDown()

	flag.Parse()

	// Print on stdout the list of the episodes available
	// Then match the episode to an index
	episodes := FetchAnime(*inputPtr)

	index := make(map[int]string)
	for i, ep := range episodes {
		fmt.Printf("ID:%d \t %s\n", i, ep)
		index[i] = ep
	}

	var epDownload []string
	epDownload = append(epDownload, index[*idPtr])

	if *idPtr == -1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check if an episodes to download is provided
	//	if ok {
	// Start goroutines pool
	fmt.Println(epDownload)
	Pool(epDownload)
	//	}
}
