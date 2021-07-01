package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/schollz/progressbar/v3"
)

var extractAlbumsRegex = regexp.MustCompile("href=\"(/?(?:album|track)/.*)\"")

const inputURL string = "https://malcura.bandcamp.com/music/"

func main() {
	albumsToDownload := getAlbumsToDownload(inputURL)
	var albums []Album
	var output string
	output += "Found albums:\n"
	for _, albumURL := range albumsToDownload {
		output += " - " + albumURL.String() + "\n"
	}
	log.Println(strings.TrimSpace(output))

	log.Println("Getting album metadata")
	for _, albumURL := range albumsToDownload {
		album := getAlbumData(albumURL)
		albums = append(albums, album)
	}

	if len(albums) == 0 {
		log.Println("Could not find any tracks.")
		return
	}

	totalSize := 0
	log.Println("Computing total size")
	for _, album := range albums {
		totalSize += album.getSize()
	}

	log.Println("Total size", convertToHumanReadableSize(totalSize))

	log.Println("Starting download")
	bar := progressbar.DefaultBytes(int64(totalSize), "Downloading...")

	for _, album := range albums {
		log.Println("Downloading album", album.Current.Title)
		for _, track := range album.Trackinfo {

			response, err := http.Get(track.File.Mp3128)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			path := filepath.Join(cwd, album.Artist, album.Current.Title, track.Title+".mp3")

			err = os.MkdirAll(filepath.Dir(path), fs.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			io.Copy(io.MultiWriter(f, bar), response.Body)
		}
	}
	log.Println("Finished download")

}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func convertToHumanReadableSize(size int) string {
	suffixes := [5]string{"B", "KB", "MB", "GB", "TB"}
	if size == 0 {
		return "0" + suffixes[0]
	}
	base := math.Log(float64(size)) / math.Log(1024)

	value := Round(math.Pow(1024, base-math.Floor(base)), 2)
	unit := suffixes[int(math.Floor(base))]
	return fmt.Sprintf("%.f %s", value, unit)
}

func getAlbumsToDownload(inputURL string) []url.URL {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatal(err)
	}

	var albumsToDownload []url.URL
	u.Path = strings.TrimPrefix(u.Path, "/")
	u.Path = strings.TrimSuffix(u.Path, "/")
	u.Path = strings.TrimSuffix(u.Path, "music")

	log.Println("Determining URL type", u.String())

	if u.Path != "" {
		elements := strings.Split(u.Path, "/")
		fmt.Println(elements)
		switch elements[0] {
		case "album":
			log.Println("URL type is album")
			// Download album
			albumsToDownload = []url.URL{*u}
		case "track":
			log.Println("URL type is track")
			// Download track
			albumsToDownload = []url.URL{*u}
		}

	} else {
		log.Println("URL type is artist page")
		artistPage := url.URL{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
			Path:   "/music",
		}
		albumsToDownload = getDiscographyURLs(artistPage)
	}

	return albumsToDownload
}

// getDiscographyURLs fetches the artist page and extracts
// the album URLs.
func getDiscographyURLs(artistPage url.URL) []url.URL {
	log.Print("Getting artist page")
	resp, err := http.Get(artistPage.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	source := string(body)
	log.Print("Extracting album URLs")
	albumPaths := extractAlbumsRegex.FindAllStringSubmatch(source, -1)

	albumURLs := make([]url.URL, 0)
	for _, path := range albumPaths {
		url := url.URL{
			Scheme: artistPage.Scheme,
			Host:   artistPage.Host,
			Path:   path[1],
		}
		albumURLs = append(albumURLs, url)
	}
	return albumURLs
}

// getAlbumData fetches the album page and extracts
// album metadata
func getAlbumData(albumURL url.URL) Album {
	resp, err := http.Get(albumURL.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	startString := "data-tralbum=\"{"
	stopString := "}\""

	htmlCode := string(body)
	indexOfStartString := strings.Index(htmlCode, startString)

	if indexOfStartString == -1 {
		log.Fatal("Could not find string in html")
	}

	albumDataRaw := htmlCode[indexOfStartString+len(startString)-1:]
	albumDataRaw = albumDataRaw[0 : strings.Index(albumDataRaw, stopString)+1]
	albumDataString := html.UnescapeString(albumDataRaw)
	var album Album
	json.Unmarshal([]byte(albumDataString), &album)
	album.fixAlbum()
	return album
}

func download() {

}
