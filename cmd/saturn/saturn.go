package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"git.mrkeebs.eu/saturn/internal/animesaturn"
	"nullprogram.com/x/optparse"
)

var (
	plot               string
	search             string
	fetch              string
	animeSaturnScraper animesaturn.AnimeScraper
)

func init() {
	animeSaturnScraper = animesaturn.AnimeSaturnScraper{}
}

func runSearch(input string) {
	available := animeSaturnScraper.SearchAnime(input)

	for _, a := range available {
		fmt.Println(path.Base(a))
	}
}

func runFetch(input string) (map[int]string, int) {
	episodes := animeSaturnScraper.FetchAnime(input)

	index := make(map[int]string)

	i := 1
	for _, ep := range episodes {
		if *idPtr == "" {
			fmt.Printf("ID:%d \t %s\n", i, ep)
		}
		index[i] = ep
		i++
	}
	return index, len(episodes) - 1
}

func main() {
	options := []optparse.Option{
		{"plot", 'p', optparse.KindRequired},
		{"search", 's', optparse.KindRequired},
		{"fetch", 'f', optparse.KindRequired},
	}
	results, _, err := optparse.Parse(options, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		switch result.Long {
		case "plot":
			plot = result.Optarg
			fmt.Println(plot)
		case "search":
			search = result.Optarg
			runSearch(search)
		}
	}
}

// package main

// import (
// 	"flag"
// 	"fmt"
// 	"path"
// 	"strconv"
// 	"strings"
// )

// var (
// 	// bar = NewBar()

// 	plot   string
// 	search string

// 	// plotPtr   = flag.String("plot", "", "Read the plot of the given anime")
// 	// searchPtr = flag.String("search", "", "Search for an anime")
// 	// inputPtr  = flag.String("fetch", "", "Fetch the available episodes for the selected anime.")
// 	// idPtr     = flag.String("down", "", "Episodes to fetch")
// 	// streamPtr = flag.String("stream", "", "Episodes to stream")
// )

// // runSearch searches all the available anime that can be
// // fetched. Then print them to stdout
// func runSearch(input *string) {
// 	available := SearchAnime(*input)

// 	for _, a := range available {
// 		fmt.Println(path.Base(a))
// 	}
// }

// // runFetch prints on stdout the list of the episodes available, each
// // one with an index.
// func runFetch(input *string) (map[int]string, int) {
// 	episodes := FetchAnime(*input)

// 	index := make(map[int]string)

// 	i := 1
// 	for _, ep := range episodes {
// 		if *idPtr == "" {
// 			fmt.Printf("ID:%d \t %s\n", i, ep)
// 		}
// 		index[i] = ep
// 		i++
// 	}
// 	return index, len(episodes) - 1
// }

// // getEp gets a string as input and return a slice of ints of all the
// // episodes to download or stream.
// func getEp(str string, eCount int) []int {
// 	var ids []int

// 	// check if we need to
// 	// download from an episode to another one; or
// 	// we need to download different episodes
// 	checkComma := strings.Index(str, ",")
// 	checkSep := strings.Index(str, "-")

// 	// Here for different episodes!
// 	if checkComma != -1 {
// 		s := strings.Split(str, ",")
// 		for _, value := range s {
// 			num, _ := strconv.Atoi(value)
// 			ids = append(ids, num)
// 		}
// 		return ids
// 	}

// 	// Here for "from - to"
// 	if checkSep != -1 {
// 		t := strings.Split(str, "-")
// 		from, _ := strconv.Atoi(t[0])
// 		//		to, _ := strconv.Atoi(t[len(t)-1])
// 		var to int
// 		if t[1] == "" {
// 			to = eCount
// 		} else {
// 			to, _ = strconv.Atoi(t[1])
// 		}
// 		for i := from; i <= to; i++ {
// 			ids = append(ids, i)
// 		}
// 		return ids
// 	}

// 	// Now if the input is a single episode
// 	num, _ := strconv.Atoi(str)
// 	ids = append(ids, num)

// 	return ids
// }

// func main() {
// 	flag.Parse()

// 	// if -search is passed
// 	if *searchPtr != "" {
// 		runSearch(searchPtr)
// 	}

// 	// if -fetch is passed
// 	var maxEp int
// 	var index map[int]string
// 	if *inputPtr != "" {
// 		//		index = runFetch(inputPtr)
// 		index, maxEp = runFetch(inputPtr)
// 	}

// 	// if -down is passed range over the
// 	// numbers of episodes selected, then
// 	// append the string(url) associated to the id
// 	if *idPtr != "" {
// 		ids := getEp(*idPtr, maxEp)
// 		// fmt.Println(ids)
// 		var episodes []string
// 		for _, id := range ids {
// 			episodes = append(episodes, index[id])
// 		}
// 		// Call the pool of goroutines with the selected episodes
// 		Pool(episodes, inputPtr)
// 	}

// 	// if -stream is passed get the episodes
// 	// selected and then start streaming them via mpv
// 	if *streamPtr != "" {
// 		ids := getEp(*streamPtr, maxEp)
// 		fmt.Println(ids)
// 		var episodes []string
// 		for _, id := range ids {
// 			episodes = append(episodes, index[id])
// 		}
// 		Stream(episodes, inputPtr)
// 	}

// 	if *plotPtr != "" {
// 		PlotAnime(*plotPtr)
// 	}
// }
