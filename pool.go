package main

import (
	"log"
)

const workers = 10 // Numbers of goroutines

// Spawn n goroutines that concurrently download the files
// streamed to the channel "in"
func Pool(epToDownload []string) {

	in := make(chan Anime)
	done := make(chan struct{})

	// Call the goroutines and let them in "listen" thought
	// the channel in, when the goroutines finish send
	// a message to "done" channel
	for i := 0; i < workers; i++ {
		go func() {
			Download(in)
			done <- struct{}{}
		}()

	}

	// Fetch the episode selected and stream the struct
	for _, ep := range epToDownload {
		FetchEpisodes(ep, in)
	}
	// Close the channel when finish the stream
	close(in)

	// Wait until all the goroutines send a "done" signal
	for a := 0; a < workers; a++ {
		<-done
	}

	log.Println("DONE")
}
