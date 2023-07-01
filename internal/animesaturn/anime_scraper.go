package animesaturn

type AnimeScraper interface {
	SearchAnime(input string) []string
	FetchAnime(input string) []string
	PlotAnime(input string)
	FetchEpisodes(episode string, out chan<- Anime, anime string)
}

type Anime struct {
	Name string
	URL  string
	// The name of the anime that must be provided as referrer to
	// bypass cloudflare.
	An string
}
