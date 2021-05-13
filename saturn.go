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
const workers = 5

type Anime struct {
	name  string
	epURL string
}

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

func downProgress(delay time.Duration) {
	for {
		for _, r := range `\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func Download(in <-chan Anime, done chan<- struct{}) error {
	for ep := range in {
		go downProgress(100 * time.Millisecond)
		nameFile := strings.ReplaceAll(ep.name, " ", "") + ".mp4"
		ffmpeg := newCmd(ep.epURL, nameFile)

		if err := ffmpeg.Run(); err != nil {
			fmt.Println(err)
		}
	}
	done <- struct{}{}
	return nil
}

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

func fetchEpisodes(episode string, out chan<- Anime) {
	// Find Watch Episode URL
	doc, _ := goquery.NewDocument(episode)
	epUrl, _ := doc.Find(".card-body a").Attr("href")

	// Find mp4 url
	d, _ := goquery.NewDocument(epUrl)
	mp4, _ := d.Find(".hero-unit source").Attr("src")
	name := d.Find(".text-white").Eq(0).First().Text()

	var a Anime
	a.epURL = string(mp4)
	a.name = string(name)
	out <- a
	// out <-  Anime{
	// 	epURL: string(mp4),
	// 	name:  string(name),
	// }
}

func pool(epToDownload string) {
	//	var wg sync.WaitGroup
	in := make(chan Anime)
	done := make(chan struct{})

	for i := 0; i < workers; i++ {
		//		wg.Add(1)
		go func() {
			Download(in, done)
			//			wg.Done()
		}()
	}

	fetchEpisodes(epToDownload, in)
	log.Println("Before Wait")

	//	wg.Wait()
	<-done
	close(in)

	log.Println("DONE")
}

func userInput() (res []int) {
	var input string
	fmt.Print("Choose which Episodes download, comma separated es: 1,2,3 :")
	fmt.Scan(&input)

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
	epToDown := userInput()

	// Start goroutines pool
	for _, n := range epToDown {
		pool(index[n])
	}
}
