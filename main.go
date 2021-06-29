package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var extractAlbumsRegex = regexp.MustCompile("href=\"(/?(?:album|track)/.*)\"")

const inputURL string = "https://malcura.bandcamp.com/album/malcura-ii"

func main() {
	fmt.Println("Fetching metadata...")
	artistPage := formatArtistPageURL(inputURL)
	album_urls := getAlbumURLs(artistPage)

	var albums []Album
	for _, albumURL := range album_urls {
		album := getAlbumData(albumURL)
		albums = append(albums, album)
	}
	fmt.Println("Finished.")
}

// formatArtistPageURL formats the input URL so it matches
// the expected format (https://malcura.bandcamp.com/music)
func formatArtistPageURL(inputURL string) url.URL {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatal(err)
	}

	artistPage := url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   "/music",
	}
	return artistPage
}

// getAlbumURLs fetches the artist page and extracts
// the discography URLs.
func getAlbumURLs(artistPage url.URL) []url.URL {
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
