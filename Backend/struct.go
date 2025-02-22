package Web

import "text/template"


// artist struct
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// Location's variables/struct
type Index []struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}
type Location struct {
	Index Index `json:"index"`
}

// Date's variables/struct
type Index_Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Date struct {
	Index_Date []Index_Date `json:"index"`
}

// Relation's variables/struct
type Index_Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type Relation struct {
	Index []Index_Relation `json:"index"`
}

var templates, err1 = template.ParseGlob("./Frontend/*.html")
