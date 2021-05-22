/*
   This file use a const of workers ( i.e. goroutines )
   for concurrently download the episodes.
   Create two channels, scrape the episode URL (with FetchEpisodes),
   and pass to the goroutines the Anime struct, which contains the
   name and the URL of the mp4.
   Then the goroutines download the file and send a signal when end.
*/

package main

import (
	"log"
)

// Number of goroutines
const workers = 10

// Spawn n goroutines that concurrently download the files
// streamed to the channel "in"
func Pool(epToDownload []string, anime *string) {
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
		FetchEpisodes(ep, in, *anime)
	}
	// Close the channel when finish the stream
	close(in)

	// Wait until all the goroutines send a "done" signal
	for a := 0; a < workers; a++ {
		<-done
	}

	log.Println("DONE")
}
