/*
   Concurrently download the url passed by the channel,
   also wrap the io.ReadCloser interfaces and implement
   a visual progress download update
*/
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ReadSpy implement an io.ReadCloser for reading
// the # of bytes downloaded
type ReadSpy struct {
	r io.ReadCloser
}

// Read implements io.Reader for ReadSpy.  It collects sends the
// number of bytes written to the channel.
func (r *ReadSpy) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	bar.Fetched(n)
	return
}

// Close implements io.Closer for ReadSpy.  It will also close the
// channel.
func (r *ReadSpy) Close() error {
	bar.DoneEpisode()
	return r.r.Close()
}

// byteConv converts the given number of bytes to megabytes, for a
// better human readability.
func byteConv(byte int) float64 {
	const unit = 1024
	if byte < unit {
		return float64(byte)
	}

	div, exp := int64(unit), 0
	for n := byte / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return float64(byte) / float64(div)
}

// Download fetches the url of the selected episode, creates a file
// with the name of the episode and then download it.  Also implement
// a progress bar.
func DownloadMP4(ep Anime) {
	// The name of the episode is full of whitespace
	// let's clean it a bit.
	nameFile := strings.TrimSpace(ep.Name) + ".mp4"
	out, err := os.Create(nameFile)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	req, err := http.NewRequest("GET", ep.URL, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Referer", ep.An)
	// fmt.Println(ep.URL)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// Retrive the Size of the file that will be downloaded,
	// and use it for the progress bar
	contentLenght := resp.Header.Get("Content-Length")
	size, _ := strconv.ParseInt(contentLenght, 10, 64)
	bar.AddEpisode(int(size))

	// Wrap the interfaces
	src := &ReadSpy{r: resp.Body}
	defer src.Close()

	// Finally write the content of src ( wrap of res.Body ) into
	// the file we create
	_, err = io.Copy(out, src)
	if err != nil {
		log.Println(err)
	}
}
