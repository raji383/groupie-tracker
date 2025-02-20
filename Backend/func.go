package Web
import (
	"encoding/json"
	"strings"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
)
//artist struct
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}
//Location's variables/struct
type Index []struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}
type Location struct {
	Index Index `json:"index"`
}
//Date's variables/struct
type Index_Date struct{
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Date struct {
	Index_Date []Index_Date `json:"index"`
}
//Relation's variables/struct
type Index_Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type Relation struct {
	Index []Index_Relation `json:"index"`
}
var templates, err1 = template.ParseGlob("./Frontend/*.html")
//this fucntion fetch Data from APIs
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
//this function Serves the home page by fetching artists data for the home page
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
		http.Error(w, "Internal Server Error", 500)
	}
}
//this function prepare the data for the artist page
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
	//fetching artist's data
	var artists []Artist
	err = FetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//fetching Locations' data
	var locations Location
	err = FetchData("https://groupietrackers.herokuapp.com/api/locations", &locations)
	if err != nil {
		log.Printf("Error fetching locations: %v", err)
		http.Error(w, "Internal Server Error11", http.StatusInternalServerError)
		return
	}
	//fetching Dates' data
	var dates Date
	err = FetchData("https://groupietrackers.herokuapp.com/api/dates", &dates)
	if err != nil {
		log.Printf("Error fetching dates: %v", err)
		http.Error(w, "Internal Server Error2", http.StatusInternalServerError)
		return
	}
	//fetching Relations' data
	var relations Relation
	err = FetchData("https://groupietrackers.herokuapp.com/api/relation", &relations)
	if err != nil {
		log.Printf("Error fetching relations: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//Preparing artist's data to be returned
	var artist Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a
			break
		}
	}
	//check the artist's id
	if artist.ID == 0 {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}
	//Preparing artist's data to be returned
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
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}
	//Preparing Dates' data to be returned
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
	//Preparing Relations' data to be returned
	var DL struct{
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
    Artist    Artist
		Locations []string `json:"locations"`
		Dates []string `json:"dates"`
		Relations_map map[string][]string
	}
	data := Data{
		Artist:    artist,
		Locations: location.Locations,
		Dates: wakt.Dates,
		Relations_map: DL.DatesLocations,
	}
	//checking errors
	if err1 != nil{
		log.Printf("Internal server error: %v", err1)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	//check if the user's input has the artist is + another link
	if len(splitted) >= 2{
		http.Error(w, "Page not found", 404)
		return
	}
	err = templates.ExecuteTemplate(w, "artist.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", 500)
	}
}
//This fucntion check the link (/css) & serve the css
func Test(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path == "/css/") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	// Serve the CSS file
	cssFile := "Frontend/css/" + r.URL.Path[len("/css/"):] // this : r.URL.Path[len("/css/"):] gives me (for example) index.css
	http.ServeFile(w, r, cssFile)
}