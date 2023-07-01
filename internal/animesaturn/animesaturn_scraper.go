/*
   This file use goquery for scraping animesaturn.tv
*/

package animesaturn

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.animesaturn.tv"

type AnimeSaturnScraper struct{}

// SearchAnime returns a string of the availables anime that can be
// fetched and downloaded.
func (s AnimeSaturnScraper) SearchAnime(input string) (available []string) {
	search := baseURL + "/animelist?search=" + input
	doc, _ := goquery.NewDocument(search)

	doc.Find(".badge-archivio").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		available = append(available, string(link))
	})
	return available
}

// FetchAnime returns a list of all the episodes from the anime
// searched.
func (s AnimeSaturnScraper) FetchAnime(input string) (episodes []string) {
	search := baseURL + "/anime/" + input
	req, _ := http.NewRequest("GET", search, nil)

	// Add Referer to bypass cloudflare
	req.Header.Set("Referer", search)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	doc.Find("a.btn.btn-dark.mb-1").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		episodes = append(episodes, string(link))
	})
	return episodes
}

// PlotAnime returns the plot from the anime searched.
func (s AnimeSaturnScraper) PlotAnime(input string) {
	var thePlot string
	splot := baseURL + "/anime/" + input
	req, err := http.Get(splot)

	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		log.Fatalf(
			"got an error status %d %s",
			req.StatusCode,
			req.Status,
		)
	}
	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		log.Println(err)
	}

	doc.Find("#trama").Each(func(i int, s *goquery.Selection) {
		plot := s.Find("#full-trama").Text()
		thePlot = plot
	})

	fmt.Println(thePlot)

}

// FetchEpisodes finds the download link and stream it in the channel
// as an Anime Struct
func (s AnimeSaturnScraper) FetchEpisodes(episode string, out chan<- Anime, anime string) {
	// Find Watch Episode URL
	doc, _ := goquery.NewDocument(episode)
	epUrl, _ := doc.Find(".card-body a").Attr("href")

	// Find .mp4 or .m3u8 url
	d, _ := goquery.NewDocument(epUrl)
	mp4, exist := d.Find(".hero-unit source").Attr("src")
	name := d.Find(".text-white").Eq(0).First().Text()

	// Check if the mp4 url exist, if not let's search for an m3u8 url
	if !exist {
		m3u := d.Find("script").Text()
		out <- Anime{
			URL:  m3u,
			Name: name,
			An:   baseURL + "/anime/" + anime,
		}
	} else {
		out <- Anime{
			URL:  mp4,
			Name: name,
			An:   baseURL + "/anime/" + anime,
		}
	}
}
