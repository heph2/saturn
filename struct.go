package main

// This struct contains the Name and the MP4 URL
// of the Episode. Also contains the Name of the anime (An)
// that must be provided as referrer in an http.Get request
// for bypassing 403 Forbidden Errors
type Anime struct {
	Name string
	URL  string
	An   string
}
