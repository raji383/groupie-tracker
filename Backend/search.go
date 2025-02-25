package Web

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {

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

	name := strings.ToLower(r.Form.Get("searched"))

	var filteredArtists []Artist
	var str []string

	for _, a := range locations.Index {
		str = append(str, a.Locations...)
	}

	for r, i := range artists {

		x := true
		str1 := strings.ToLower(i.Name)
		if strings.Contains(str1, name) {

			filteredArtists = append(filteredArtists, i)
			x = false

		}
		if x && name == i.FirstAlbum {
			filteredArtists = append(filteredArtists, i)
			x = false
		}
		if x && name == strconv.Itoa(i.CreationDate) {
			filteredArtists = append(filteredArtists, i)
			x = false
		}
		for _, j := range i.Members {
			str2 := strings.ToLower(j)
			if x && strings.Contains(str2, name) {
				filteredArtists = append(filteredArtists, i)
				x = false
			}
		}
		for _, a := range locations.Index {
			if a.ID == r+1 {
				for _, j := range a.Locations {
					str3 := strings.ToLower(j)

					if x && strings.Contains(str3, name) {
						filteredArtists = append(filteredArtists, i)
						x = false
					}
				}
			}
		}
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
