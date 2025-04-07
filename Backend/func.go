package Web

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)
type ErrorData struct {
	Errr    int
	Kalma string
}

func renderErrorPage(w http.ResponseWriter, statusCode int, message string) {
	data := ErrorData{
		Errr:    statusCode,
		Kalma: message,
	}

	w.WriteHeader(statusCode)
	if err := templates.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, message, statusCode)
	}
}

// this function Serves the home page by fetching artists data for the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		
		data := struct {
			Errr    int
			Kalma string
		}{
			Errr:    404,
			Kalma: "Page not found",
		}
	
		
		w.WriteHeader(http.StatusNotFound)
		

		if err := templates.ExecuteTemplate(w, "error.html", data); err != nil {
		
			http.Error(w, "Page not found", http.StatusNotFound)
		}
		return
	}
	
	if r.Method != http.MethodGet {
		data := struct {
			Errr    int
			Kalma string
		}{
			Errr:    405,
			Kalma: "Method Not Allowed",
		}
	
		
		w.WriteHeader(http.StatusMethodNotAllowed)
		

		if err := templates.ExecuteTemplate(w, "error.html", data); err != nil {
		
			http.Error(w, "Page not found", http.StatusNotFound)
		}
		return
	}
	var artists []Artist
	err := FetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		renderErrorPage(w,500,"itrnal sirval error")
		return
	}
	var locations Location
	err = FetchData("https://groupietrackers.herokuapp.com/api/locations", &locations)
	if err != nil {
		renderErrorPage(w,500,"itrnal sirval error")
		return
	}
	if err1 != nil {
		renderErrorPage(w,500,"itrnal sirval error")
		return
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
		Artist:    artists,
		Locations: str,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		renderErrorPage(w,500,"itrnal sirval error")
	}
}

// this function prepare the data for the artist page
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderErrorPage(w,405,"Method Not Allowed")
		return
	}
	artistID := r.URL.Query().Get("id")
	splitted := strings.Split(artistID, "/")
	id, err := strconv.Atoi(splitted[0])
	if err != nil || id < 0 {
		// log.Printf("Page not found: %v", err)
		renderErrorPage(w,400,"bad requst")
		return
	}
	// fetching artist's data
	var artists []Artist
	err = FetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		renderErrorPage(w,500,"itrnal sirval error")
		return
	}
	// fetching Locations' data
	var locations Location
	err = FetchData("https://groupietrackers.herokuapp.com/api/locations", &locations)
	if err != nil {
		renderErrorPage(w,500,"itrnal sirval error")
		return
	}
	// fetching Dates' data
	var dates Date
	err = FetchData("https://groupietrackers.herokuapp.com/api/dates", &dates)
	if err != nil {
		renderErrorPage(w,500,"itrnal sirval error")
		return
	}
	// fetching Relations' data
	var relations Relation
	err = FetchData("https://groupietrackers.herokuapp.com/api/relation", &relations)
	if err != nil {
		renderErrorPage(w,500,"itrnal sirval error")
		return
	}
	// Preparing artist's data to be returned
	var artist Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a
			break
		}
	}
	// check the artist's id
	if artist.ID == 0 {		
		renderErrorPage(w,404,"not found")
		return
	}
	// Preparing artist's data to be returned
	var location struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	}
	for _, a := range locations.Index {
		if a.ID == id {
			location = a
			break
		}
	}
	if location.ID == 0 {
		renderErrorPage(w,404,"not found")
		return
	}
	// Preparing Dates' data to be returned
	var wakt struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}
	for _, a := range dates.Index_Date {
		if a.ID == id {
			wakt = a
			break
		}
	}
	if wakt.ID == 0 {
		renderErrorPage(w,404,"not found")
		return
	}
	// Preparing Relations' data to be returned
	var DL struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}
	for _, a := range relations.Index {
		if a.ID == id {
			DL = a
			break
		}
	}
	// Data to be printed
	type Data struct {
		Artist        Artist
		Locations     []string `json:"locations"`
		Dates         []string `json:"dates"`
		Relations_map map[string][]string
	}
	data := Data{
		Artist:        artist,
		Locations:     location.Locations,
		Dates:         wakt.Dates,
		Relations_map: DL.DatesLocations,
	}
	// checking errors
	if err1 != nil {
		
		renderErrorPage(w,500,"itrnal sirval error")
		return
	}
	// check if the user's input has the artist is + another link
	if len(splitted) >= 2 {
		
		renderErrorPage(w,404,"not fund")
		return
	}
	err = templates.ExecuteTemplate(w, "artist.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", 500)
	}
}
