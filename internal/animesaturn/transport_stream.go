package animesaturn

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
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

// sanitizeURL search inside the <script> and find the .m3u8 url
// that end with playlist.m3u8, after that trim the url.
// Return an url of this type: https://.../
func sanitizeURL(epURL string) (playlistURL string) {
	re := regexp.MustCompile(`https:\/\/.*.m3u8`)
	playlistURL = re.FindString(epURL)
	playlistURL = strings.TrimSuffix(playlistURL, "playlist.m3u8")
	return playlistURL
}

// URL of kind playlist.m3u8 and return kind 720p.m3u8
func getResolution(playlistURL string) (episodeURL string, size int64) {
	res, _ := http.Get(playlistURL)
	contentLenght := res.Header.Get("Content-Length")
	size, _ = strconv.ParseInt(contentLenght, 10, 64)
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyStrings := strings.Split(string(bodyBytes), "\n")

	// range over bodyStrings and search for lines with .m3u8 links
	var resolutions []string
	for _, line := range bodyStrings {
		if strings.Contains(line, ".m3u8") {
			resolutions = append(resolutions, line)
		}
	}

	log.Println(resolutions)

	// Always find for 720p resolution
	var maxRes string
	maxRes = resolutions[len(resolutions)-1]
	for _, res := range resolutions {
		if strings.Contains(res, "720p") {
			maxRes = res
			break
		}
		if strings.Contains(res, "480p") {
			maxRes = res
			break
		}
		if strings.Contains(res, "240p") {
			maxRes = res
			break
		}
		if strings.Contains(res, "144p") {
			maxRes = res
			break
		}
	}

	log.Println("This is ", maxRes)

	// create base URL
	base, err := url.Parse(playlistURL)
	if err != nil {
		log.Fatal(err)
	}
	maxRes = strings.TrimSuffix(maxRes, "\r")
	u, err := url.Parse(maxRes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(base.ResolveReference(u))

	episodeURL = strings.Replace(playlistURL, "playlist.m3u8", maxRes, -1)

	return
}

// DownloadM3U retrive the .m3u8 URL, than start an http request with
// that URL. This will return a list of all the URLs for the fragments
// of the playlist. After that we retrive those URLs, download them, and
// create a single mp4 file
func DownloadM3U(ep Anime) {
	// retrive sanitez URL
	baseURL := sanitizeURL(ep.URL)
	playlistURL := baseURL + "playlist.m3u8"
	episodeURL, _ := getResolution(playlistURL) // this ends with 720p.m3u8

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
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error with ts url", err)
	}
	str := strings.Split(string(bytes), "\n")

	// Retrive all the .ts link
	var ts []string
	for _, s := range str {
		if strings.Contains(s, ".ts") {
			ts = append(ts, s)
		}
	}

	// Build the url appending the ts link to baseURL
	var downURL []string
	for _, t := range ts {
		//downURL = append(downURL, baseURL+t)
		u, err := url.Parse(t)
		if err != nil {
			log.Fatal(err)
		}
		base, err := url.Parse(episodeURL)
		abs := base.ResolveReference(u)
		log.Println(abs)
		downURL = append(downURL, abs.String())
	}

	// Start Downloading each video
	dataArray := downloadMultipleFiles(downURL)

	ep.Name = strings.TrimSpace(ep.Name) + ".mp4"
	file, err := os.Create(ep.Name)
	if err != nil {
		log.Println("Error creating file: ", err)
	}
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
