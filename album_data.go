package main

import (
	"log"
	"time"
)

type Album struct {
	Trackinfo   []Track
	ReleaseTime time.Time
	Current     struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
	} `json:"current"`
	Artist   string `json:"artist"`
	ItemType string `json:"item_type"`
	ArtID    int    `json:"art_id"`
}

// getSize returns the totalled size of all the tracks in the album
func (a Album) getSize() int {
	total := 0
	for _, track := range a.Trackinfo {
		size, err := track.getSize()
		if err != nil {
			log.Fatal(err)
		}
		total += size
	}
	return total
}

// fixAlbum will convert date format from string to time.Time
// and add a pointer to the album to each track
func (a Album) fixAlbum() {
	layout := "02 Jan 2006 15:04:05 GMT"
	var err error
	a.ReleaseTime, err = time.Parse(layout, a.Current.ReleaseDate)
	if err != nil {
		log.Fatal("Could not parse the release date.", err)
	}
	for _, track := range a.Trackinfo {
		track.Album = &a
	}
}

// func (a Album) download() {
// 	for _, track := range a.Trackinfo {
// 		cwd, err := os.Getwd()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		track.File.Path = filepath.Join(cwd, a.Artist, a.Current.Title, track.Title+".mp3")
// 		fmt.Println("Downloading", track.Title, "to", track.File.Path)
// 		data := track.download()

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		err = os.MkdirAll(filepath.Dir(track.File.Path), fs.ModePerm)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		err = ioutil.WriteFile(track.File.Path, data, fs.ModePerm)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }
