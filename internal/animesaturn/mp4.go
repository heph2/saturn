/*
Concurrently download the url passed by the channel,
also wrap the io.ReadCloser interfaces and implement
a visual progress download update
*/
package animesaturn

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// Download fetches the url of the selected episode, creates a file
// with the name of the episode and then download it.  Also implement
// a progress bar.
func DownloadMP4(ep Anime) {
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// Retrive the Size of the file that will be downloaded,
	// and use it for the progress bar
	// contentLenght := resp.Header.Get("Content-Length")

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		log.Println(err)
	}
}
