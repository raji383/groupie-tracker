package Web

import (
	"log"
	"net/http"
	"strconv"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
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

	creationFromStr := r.Form.Get("creation_from")
	creationToStr := r.Form.Get("creation_to")
	minMembersStr := r.Form.Get("min_members")
	maxMembersStr := r.Form.Get("max_members")
	firstAlbumFromStr := r.Form.Get("first_album_from")
	firstAlbumToStr := r.Form.Get("first_album_to")

	var creationFrom, creationTo int
	if creationFromStr != "" {
		if val, err := strconv.Atoi(creationFromStr); err == nil {
			creationFrom = val
		}
	}
	if creationToStr != "" {
		if val, err := strconv.Atoi(creationToStr); err == nil {
			creationTo = val
		}
	}

	minMembers, maxMembers := 0, 1000
	if minMembersStr != "" {
		if val, err := strconv.Atoi(minMembersStr); err == nil {
			minMembers = val
		}
	}
	if maxMembersStr != "" {
		if val, err := strconv.Atoi(maxMembersStr); err == nil {
			maxMembers = val
		}
	}

	var firstAlbumFrom, firstAlbumTo int
	if firstAlbumFromStr != "" {
		if val, err := strconv.Atoi(firstAlbumFromStr); err == nil {
			firstAlbumFrom = val
		}
	}
	if firstAlbumToStr != "" {
		if val, err := strconv.Atoi(firstAlbumToStr); err == nil {
			firstAlbumTo = val
		}
	}

	filteredArtists := []Artist{}
	for _, artist := range artists {

		if creationFrom != 0 && artist.CreationDate < creationFrom {
			continue
		}
		if creationTo != 0 && artist.CreationDate > creationTo {
			continue
		}

		numMembers := len(artist.Members)
		if numMembers < minMembers || numMembers > maxMembers {
			continue
		}

		if firstAlbumFrom != 0 || firstAlbumTo != 0 {
			if len(artist.FirstAlbum) < 4 {
				continue
			}
			albumYear, err := strconv.Atoi(artist.FirstAlbum[:4])
			if err != nil {
				continue
			}
			if firstAlbumFrom != 0 && albumYear < firstAlbumFrom {
				continue
			}
			if firstAlbumTo != 0 && albumYear > firstAlbumTo {
				continue
			}
		}

		filteredArtists = append(filteredArtists, artist)
	}

	err = templates.ExecuteTemplate(w, "index.html", filteredArtists)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", 500)
	}
}
