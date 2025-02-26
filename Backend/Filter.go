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

    var locData Location
    err = FetchData("https://groupietrackers.herokuapp.com/api/locations", &locData)
    if err != nil {
        log.Printf("Error fetching locations: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    creationFromStr, _ := strconv.Atoi(r.Form.Get("creation_from"))
    creationToStr, _ := strconv.Atoi(r.Form.Get("creation_to"))
    MembersStr, _ := strconv.Atoi(r.Form.Get("members"))

    firstAlbumFromStr := strings.Split(r.Form.Get("first_album_from"), "-")
    firstAlbumToStr := strings.Split(r.Form.Get("first_album_to"), "-")
    locationFilter := r.Form.Get("location")

    var filteredArtists []Artist

    for _, artist := range artists {
        passes := true

        
        if creationFromStr != 0 && creationToStr != 0 {
            if artist.CreationDate < creationFromStr || artist.CreationDate > creationToStr {
                passes = false
            }
        }

        if MembersStr != 0 {
            if MembersStr == 6 { 
                if len(artist.Members) < 6 {
                    passes = false
                }
            } else {
                if len(artist.Members) != MembersStr {
                    passes = false
                }
            }
        }

        if len(firstAlbumFromStr) == 3 && len(firstAlbumToStr) == 3 {
            firstAlbumParts := strings.Split(artist.FirstAlbum, "-")
            if len(firstAlbumParts) == 3 && !inrang(firstAlbumParts, firstAlbumFromStr, firstAlbumToStr) {
                passes = false
            }
        }

        if locationFilter != "" {
            var artistLocationFound bool
            
            for _, loc := range locData.Index {
                if loc.ID == artist.ID {
                    for _, locName := range loc.Locations {
                        if locName == locationFilter {
                            artistLocationFound = true
                            break
                        }
                    }
                    break
                }
            }
            if !artistLocationFound {
                passes = false
            }
        }

        if passes {
            filteredArtists = append(filteredArtists, artist)
        }
    }


    if len(filteredArtists) == 0 {
        filteredArtists = artists
    }

    var allLocations []string
    for _, loc := range locData.Index {
        allLocations = append(allLocations, loc.Locations...)
    }

    type Data struct {
        Artist    []Artist
        Locations []string `json:"locations"`
    }
    data := Data{
        Artist:    filteredArtists,
        Locations: allLocations,
    }

    err = templates.ExecuteTemplate(w, "index.html", data)
    if err != nil {
        log.Printf("Template error: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
