package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Small random progress bar
func downProgress(delay time.Duration) {
	for {
		for _, r := range `\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func Download(in <-chan Anime) {
	for ep := range in {
		go downProgress(100 * time.Millisecond)
		nameFile := strings.ReplaceAll(ep.Name, " ", "") + ".mp4"
		out, _ := os.Create(nameFile)
		defer out.Close()

		res, err := http.Get(ep.URL)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()

		if _, err := io.Copy(out, res.Body); err != nil {
			log.Println(err)
		}
	}
}
