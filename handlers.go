package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"path/filepath"

	gowebly "github.com/gowebly/helpers"
)

// indexViewHandler handles a view for the index page.
func indexViewHandler(w http.ResponseWriter, r *http.Request) {

	// Define paths to the user templates.
	indexPage := filepath.Join("templates", "pages", "index.html")

	// Parse user templates or return error.
	tmpl, err := gowebly.ParseTemplates(indexPage) // gowebly helper for parse user templates
	if err != nil {
		// If not, return HTTP 500 error.
		slog.Error(err.Error(), "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Execute (render) all templates or return error.
	if err := tmpl.Execute(w, nil); err != nil {
		// If not, return HTTP 500 error.
		slog.Error(err.Error(), "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func fetchGoogleResults(w http.ResponseWriter, r *http.Request) {
	searchString := r.FormValue("search")

	// if the search string is empty, return a 400 Bad Request
	if searchString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := Search(searchString)

	jsonResults, _ := json.Marshal(result)

	// encode JSON data as a data: URL
	dataURL := "data:application/json;charset=utf-8," + url.PathEscape(string(jsonResults))

	// return an HTML response containing a download link with the data: URL and the number of results
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`
		<p>Number of results: %d</p>
		<br/>
		<a href="%s" download="searchResults.json" class="btn">Download JSON</a>
    `, len(result), dataURL)))
}
