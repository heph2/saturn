package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ReadSpy struct {
	r  io.Reader
	ch chan int
}

func (r *ReadSpy) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.ch <- n
	return
}

// func (r *ReadSpy) Close() error {
// 	close(r.ch)
// 	return r.r.Close()
// }

func Download(in <-chan Anime) {
	for ep := range in {
		nameFile := strings.TrimSpace(ep.Name) + ".mp4"
		out, _ := os.Create(nameFile)
		defer out.Close()

		res, err := http.Get(ep.URL)
		if err != nil {
			log.Println(err)
		}
		//		defer res.Body.Close()

		contentLenght := res.Header.Get("Content-Length")
		size, _ := strconv.ParseInt(contentLenght, 10, 64)

		src := &ReadSpy{r: res.Body, ch: make(chan int)}
		//		defer src.Close()

		go func() {
			var bar Bar
			bar.NewOption(0, int(size))
			for p := range src.ch {
				bar.Play(int(p))
			}
		}()

		_, err = io.Copy(out, src)
		if err != nil {
			log.Println(err)
		}
	}
}
