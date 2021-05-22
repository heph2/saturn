package main

type Anime struct {
	// The name of the anime.
	Name string

	// URL to the MP4 file.
	URL string

	// The name of the anime that must be provided as referrer to
	// bypass cloudflare.
	An string
}
