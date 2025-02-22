package Web

import (
	"encoding/json"
	"io"
	"net/http"
)

// this fucntion fetch Data from APIs
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

// This fucntion check the link (/css) & serve the css
func Test(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/css/" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	// Serve the CSS file
	cssFile := "Frontend/css/" + r.URL.Path[len("/css/"):] // this : r.URL.Path[len("/css/"):] gives me (for example) index.css
	http.ServeFile(w, r, cssFile)
}
