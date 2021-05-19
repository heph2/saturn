/*
   This file fetch the episodes selected for streaming,
   create a tmp file in /tmp, and add each url (mp4) in that file
   separated by a newline. Then exec mpv with --playlist argument
   with that tmp file.
*/

package main

import (
	"io/ioutil"
	"log"
	"os/exec"
)

// This function create a tmpFile for storing the urls
func tmpFile(urls string) string {
	var playlistFile string

	tmp, err := ioutil.TempFile("", "playlist-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	text := []byte(urls)
	_, err = tmp.Write(text)
	if err != nil {
		log.Fatal(err)
	}

	playlistFile = "--playlist=" + tmp.Name()
	return playlistFile
}

// This function make use of goroutines for concurrently
// get the Urls and stream them via MPV
func Stream(epToStream []string, anime *string) {

	in := make(chan Anime)
	done := make(chan struct{})

	go func(in chan Anime) {

		// Now range over the selected episodes to Stream
		// and add them to a string, each one separated
		// by a newline
		// Also store the referrer option to call mpv
		var referrerOption string
		var urls string
		for ep := range in {
			urls += ep.URL + "\n"
			referrerOption = "--referrer=" + ep.An
		}

		// Create a tmp file where the playlist will be stored
		playlist := tmpFile(urls)

		// Exec mpv with the referrer and a playlist
		cmd := exec.Command("mpv", referrerOption, playlist)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}

		// When a signal from the user is send ( quit from mpv )
		// send the end to the channel
		done <- struct{}{}
	}(in)

	// Fetch the episode selected and stream the struct
	for _, ep := range epToStream {
		FetchEpisodes(ep, in, *anime)
	}

	// Close the channel when finish the stream
	close(in)

	// Wait until the goroutine send a "done" signal
	<-done

	log.Println("STREAM END")

}
