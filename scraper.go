package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.animesaturn.it"

// From an user input return a list of all the Episodes
// from an anime searched
func FetchAnime(input string) (Episodes []string) {
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

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("a.btn.btn-dark.mb-1").Each(func(i int, s *goquery.Selection) {
		var episode string
		link, _ := s.Attr("href")
		episode = string(link)

		Episodes = append(Episodes, episode)
	})
	return Episodes
}

// From an user input ( episode choosen ) find the download link
// and stream it in the channel as an Anime Struct
func FetchEpisodes(episode string, out chan<- Anime) {
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
	}
}
