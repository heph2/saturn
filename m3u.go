package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// downloadFile take an URL as input, start an http request on
// that URL and return a slice of bytes downloaded.
func downloadFile(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, response.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

// downloadMultipleFiles simply take a slice or urls and
// download them with the  function downloadFile, then return them
func downloadMultipleFiles(urls []string) [][]byte {
	bytesArray := make([][]byte, 0)
	for _, URL := range urls {
		b, _ := downloadFile(URL)
		bytesArray = append(bytesArray, b)
	}
	return bytesArray
}

func sanitizeURL(epURL string) (playlistURL string) {
	lines := strings.Split(epURL, "\n")

	var url string
	for _, line := range lines {
		if strings.Contains(line, ".m3u8") {
			r := strings.ReplaceAll(line, "file:", "")
			trim := strings.TrimSpace(r)

			// Clean
			url = trim[1 : len(trim)-2]
		}
	}
	playlistURL = url[:len(url)-13]
	return playlistURL
}

// DownloadM3U retrive the .m3u8 URL, than start an http request with
// that URL. This will return a list of all the URLs for the fragments
// of the playlist. After that we retrive those URLs, download them, and
// create a single mp4 file
func DownloadM3U(ep Anime) {

	// retrive sanitez URL
	playlistURL := sanitizeURL(ep.URL)
	episodeURL := playlistURL + "720p.m3u8"

	// At this point need to read the body of an http request
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		log.Println("Error with http req", err)
	}
	req.Header.Set("Referer", ep.An)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error with http client", err)
	}
	defer resp.Body.Close()

	// Read the playlist body, split them and retrive
	// the URLs of the videos
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := strings.Split(string(bytes), "\n")

	var ts []string
	for _, s := range str {
		if strings.Contains(s, ".ts") {
			ts = append(ts, s)
		}
	}

	var downUrl []string
	for _, t := range ts {
		downUrl = append(downUrl, playlistURL+t)
	}

	dataArray := downloadMultipleFiles(downUrl)

	ep.Name = strings.TrimSpace(ep.Name) + ".mp4"
	file, _ := os.Create(ep.Name)
	defer file.Close()

	// This loop range over all the bytes and
	// write them to the file created above
	for _, data := range dataArray {
		_, err := file.Write(data)
		if err != nil {
			log.Println(err)
		}
	}
}
