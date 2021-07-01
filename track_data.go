package main

import (
	"net/http"
)

type Track struct {
	File struct {
		Mp3128 string `json:"mp3-128"`
		Path   string
	} `json:"file"`
	TrackNum int         `json:"track_num"`
	Lyrics   interface{} `json:"lyrics"`
	Title    string      `json:"title"`
	Duration float64     `json:"duration"`
	TrackID  int         `json:"track_id"`
	Album    *Album
}

// getSize fetches the headers for the track and
// returns the Content-Length header.
func (t Track) getSize() (size int, err error) {
	response, err := http.Head(t.File.Mp3128)
	if err != nil {
		return -1, err
	}

	defer response.Body.Close()
	//size, _ = strconv.Atoi(response.Header.Get("Content-Length"))
	return int(response.ContentLength), nil
}

// func (t Track) tag(reader io.Reader) bytes.Buffer {
// 	tags, err := id3v2.ParseReader(reader, id3v2.Options{Parse: true})
// 	if err != nil {
// 		log.Fatal()
// 	}
// 	tags.SetAlbum(t.Album.Current.Title)
// 	tags.SetArtist(t.Album.Artist)
// 	tags.SetYear(t.Album.Current.ReleaseDate)
// 	tags.SetTitle(t.Title)
// 	tags.SetYear(strconv.Itoa(t.Album.ReleaseTime.Year()))
// }
