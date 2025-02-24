package Web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func inrang(taem, from, to []string) bool {
	j := 0
	for i := 2; i >= 0; i-- {
		a, _ := strconv.Atoi(taem[j])
		if j < 3 {
			j++
		}

		b, _ := strconv.Atoi(from[i])
		c, err := strconv.Atoi(to[i])
		if err != nil {
			fmt.Println("err", err)
		}
		if a < c && a > b {
			return true
		} else if a > c && a < b {
			break
		}
	}
	return false
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var artists []Artist
	err := FetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var locations Location
	err = FetchData("https://groupietrackers.herokuapp.com/api/locations", &locations)
	if err != nil {
		log.Printf("Error fetching locations: %v", err)
		http.Error(w, "Internal Server Error11", http.StatusInternalServerError)
		return
	}

	creationFromStr, _ := strconv.Atoi(r.Form.Get("creation_from"))
	creationToStr, _ := strconv.Atoi(r.Form.Get("creation_to"))
	MembersStr, _ := strconv.Atoi(r.Form.Get("members"))
	firstAlbumFromStr := strings.Split(r.Form.Get("first_album_from"), "-")
	firstAlbumToStr := strings.Split(r.Form.Get("first_album_to"), "-")
	location := r.Form.Get("location")
	var filteredArtists []Artist


	for r, i := range artists {
		x:=true
		for _, a := range locations.Index {
			if a.ID == r {
				for _, j := range a.Locations {
					if x&& j == location {
						filteredArtists = append(filteredArtists, i)
						x=false
					}
				}
			}
		}

		firstAlbum := strings.Split(i.FirstAlbum, "-")
		if len(firstAlbumFromStr) == 3 && len(firstAlbumToStr) == 3 {
			if x&&len(firstAlbum) == 3 && inrang(firstAlbum, firstAlbumFromStr, firstAlbumToStr) {
				filteredArtists = append(filteredArtists, i)
				x=false
				continue
			}
		}

		if x&&creationFromStr <= i.CreationDate && creationToStr >= i.CreationDate {
			filteredArtists = append(filteredArtists, i)
			x=false
			continue
		}
		if x&&len(i.Members) == MembersStr {
			x=false
			filteredArtists = append(filteredArtists, i)
			continue
		}
	}
	if len(filteredArtists)==0 {
		filteredArtists=artists
	}
	var str []string

	for _, a := range locations.Index {
		str = append(str, a.Locations...)
	}

	type Data struct {
		Artist    []Artist
		Locations []string `json:"locations"`
	}
	data := Data{
		Artist:    filteredArtists,
		Locations: str,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", 500)
	}
}
