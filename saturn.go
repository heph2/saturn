package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.animesaturn.it"
const workers = 2 // Numbers of goroutines

// Get name and url of the episode to download
type Anime struct {
	name  string
	epURL string
}

// Wrapper to ffmpeg that download the file with low
// compression
func newCmd(inFile, outFile string) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-i", inFile,
		"-c:v",
		"libx264",
		"-preset",
		"fast",
		"-crf", "18",
		outFile,
	)
}

// Small random progress bar
func downProgress(delay time.Duration) {
	for {
		for _, r := range `\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// Range over the Anime Struct passed through the channel
// and call the ffmpeg wrapper to download
func Download(in <-chan Anime) {
	for ep := range in {
		go downProgress(100 * time.Millisecond)
		nameFile := strings.ReplaceAll(ep.name, " ", "") + ".mp4"
		ffmpeg := newCmd(ep.epURL, nameFile)

		if err := ffmpeg.Run(); err != nil {
			fmt.Println(err)
		}
	}
}

// From an user input return a list of all the Episodes
// from an anime searched
func fetchAnime(input string) (Episodes []string) {
	search := baseURL + "/anime/" + input
	doc, _ := goquery.NewDocument(search)

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
func fetchEpisodes(episode string, out chan<- Anime) {
	// Find Watch Episode URL
	doc, _ := goquery.NewDocument(episode)
	epUrl, _ := doc.Find(".card-body a").Attr("href")

	// Find .mp4 or .m3u8 url
	d, _ := goquery.NewDocument(epUrl)
	mp4, _ := d.Find(".hero-unit source").Attr("src")
	name := d.Find(".text-white").Eq(0).First().Text()

	out <- Anime{
		epURL: string(mp4),
		name:  string(name),
	}
}

// Spawn n goroutines that concurrently download the files
// streamed to the channel "in"
func pool(epToDownload string) {

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
	fetchEpisodes(epToDownload, in)
	// Close the channel when finish the stream
	close(in)

	// Wait until all the goroutines send a "done" signal
	for a := 0; a < workers; a++ {
		<-done
	}

	log.Println("DONE")
}

// Get and parse the user input and return a list of index
// associated ( by the map ) to a Episode URL
func userInput(index map[int]string) (res []int) {
	var input string
	fmt.Print("Choose which Episodes download, comma separated es: 1,2,3 :")
	fmt.Scan(&input)

	if input == "all" {
		res = append(res, len(index))
		return res
	}

	sanInput := strings.Split(input, ",")
	for i := 0; i < len(sanInput); i++ {
		v, _ := strconv.Atoi(sanInput[i])
		res = append(res, v)
	}
	return res
}

func main() {
	input := os.Args[1:]
	sanitezedInput := strings.Join(input, "-")

	index := make(map[int]string)
	episodes := fetchAnime(sanitezedInput)

	for i, ep := range episodes {
		fmt.Printf("Index: %d\t Episode: %s\n", i, ep)
		index[i] = ep
	}

	// Get user input
	epToDown := userInput(index)

	// Start goroutines pool
	for _, n := range epToDown {
		pool(index[n])
	}
}
