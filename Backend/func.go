package Web

import (
	"encoding/json"
	"strings"
	// "fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}
type Index []struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Location struct {
	Index Index `json:"index"`
}
type Index_Date struct{
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Date struct {
	Index_Date []Index_Date `json:"index"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var templates, err1 = template.ParseGlob("./Frontend/*.html")

func FetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/"{
		http.Error(w, "Page not found", 404)
		return
	}
	if r.Method != http.MethodGet{
		http.Error(w, "Method not allowed", 405)
		return
	}
	var artists []Artist
	err := FetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err1 != nil{
		log.Printf("Internal server error: %v", err1)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", artists)
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method !=  http.MethodGet {
		http.Error(w, "error 405: Method not allowed", 405)
		return
	}
	artistID := r.URL.Query().Get("id")
	splitted := strings.Split(artistID, "/")
	id, err := strconv.Atoi(splitted[0])
	if err != nil||id<0 {
		// log.Printf("Page not found: %v", err)
		http.Error(w, "Bad request", 400)
		return
	}
	

	var artists []Artist
	err = FetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
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
	var dates Date

	err = FetchData("https://groupietrackers.herokuapp.com/api/dates", &dates)
	if err != nil {
		log.Printf("Error fetching dates: %v", err)
		http.Error(w, "Internal Server Error2", http.StatusInternalServerError)
		return
	}

	var relations Relation
	err = FetchData("https://groupietrackers.herokuapp.com/api/relation", &relations)
	if err != nil {
		log.Printf("Error fetching relations: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var artist Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a

			break
		}
	}

	if artist.ID == 0 {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}
	

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
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	type Data struct {
    Artist    Artist
		Locations []string `json:"locations"`
		Dates []string `json:"dates"`
	}


	data := Data{
		Artist:    artist,
		Locations: location.Locations,
		Dates: wakt.Dates,
	}
	if err1 != nil{
		log.Printf("Internal server error: %v", err1)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if len(splitted) >= 2{
		http.Error(w, "Page not found", 404)
		return
	}

	err = templates.ExecuteTemplate(w, "artist.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path == "/css/") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Serve the CSS file
	cssFile := "Frontend/css/" + r.URL.Path[len("/css/"):]
	http.ServeFile(w, r, cssFile)
}
