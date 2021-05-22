/*
   This file use goquery for scraping animesaturn.it
*/

package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.animesaturn.it"

// SearchAnime returns a string of the availables anime that can be
// fetched and downloaded.
func SearchAnime(input string) (available []string) {
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
func FetchAnime(input string) (episodes []string) {
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

// FetchEpisodes finds the download link and stream it in the channel
// as an Anime Struct
func FetchEpisodes(episode string, out chan<- Anime, anime string) {
	// Find Watch Episode URL
	doc, _ := goquery.NewDocument(episode)
	epUrl, _ := doc.Find(".card-body a").Attr("href")

	// Find .mp4 or .m3u8 url
	d, _ := goquery.NewDocument(epUrl)
	mp4, _ := d.Find(".hero-unit source").Attr("src")
	name := d.Find(".text-white").Eq(0).First().Text()

	out <- Anime{
		URL:  string(mp4),
		Name: string(name),
		An:   baseURL + "/anime/" + anime,
	}
}
