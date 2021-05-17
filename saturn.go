package main

import (
	"flag"
	"fmt"
	"path"
	"strconv"
	"strings"
)

func runSearch(input *string) {
	available := SearchAnime(*input)

	for _, a := range available {
		fmt.Println(path.Base(a))
	}
}

// Print on stdout the list of the episodes available
// Then match the episode to an index
func runFetch(input *string) (index map[int]string) {
	episodes := FetchAnime(*input)

	index = make(map[int]string)
	for i, ep := range episodes {
		fmt.Printf("ID:%d \t %s\n", i, ep)
		index[i] = ep
	}
	return index
}

func runDown(str string) []int {
	var ids []int
	var s []string
	var t []string

	// check if we need to
	// download from an episode to another one; or
	// we need to download different episodes
	checkComma := strings.Index(str, ",")
	checkSep := strings.Index(str, "-")

	// Here for different episodes!
	if checkComma != -1 {
		s = strings.Split(str, ",")
		for _, value := range s {
			num, _ := strconv.Atoi(value)
			ids = append(ids, num)
		}
		return ids
	}

	// Here for "from - to"
	if checkSep != -1 {
		t = strings.Split(str, "-")
		from, _ := strconv.Atoi(t[0])
		to, _ := strconv.Atoi(t[len(t)-1])
		for i := from; i <= to; i++ {
			ids = append(ids, i)
		}
		return ids
	}

	// Now if the input is a single episode
	num, _ := strconv.Atoi(str)
	ids = append(ids, num)

	return ids
}

func main() {
	searchPtr := flag.String("search", "", "Search Anime")
	inputPtr := flag.String("fetch", "", "Fetch the available episodes for the anime selected.")
	idPtr := flag.String("down", "", "Episodes available")
	flag.Parse()

	if *searchPtr != "" {
		runSearch(searchPtr)
	}

	var index map[int]string
	if *inputPtr != "" {
		index = runFetch(inputPtr)
	}

	if *idPtr != "" {
		ids := runDown(*idPtr)
		fmt.Println(ids)
		var episodes []string
		for _, id := range ids {
			episodes = append(episodes, index[id])
		}
		Pool(episodes)
	}
}
