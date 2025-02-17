package main

import (
	"ascii-art-web-export-file/ascii"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// PageData holds template data
type PageData struct {
	Text     string
	Banner   string
	ASCIIArt string
	Error    string
	Banners  []string
}

// Global vars for exporting
var lastGeneratedASCII string
var lastBanner string

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/404", handler404)

	// New route for exporting ASCII art
	http.HandleFunc("/download", downloadHandler)

	fmt.Println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// ----------------------------------------------------
// 404 handler
// ----------------------------------------------------
func handler404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data, err := os.ReadFile("templates/404.html")
	if err != nil {
		w.Write([]byte("404 Not Found"))
		return
	}
	w.Write(data)
}

// ----------------------------------------------------
// homePage handler
// ----------------------------------------------------
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Template Parsing Error:", err)
		return
	}

	// Get banners from bannerList.go (already strips .txt)
	banners, err := ascii.BannerList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Banner list error:", err)
		return
	}

	data := PageData{
		Banners: banners,
	}

	if r.Method == http.MethodGet {
		if r.URL.Path != "/" {
			handler404(w, r)
			return
		}
	}

	if r.Method == http.MethodPost {
		text := r.FormValue("text")
		banner := r.FormValue("banner")

		// Validate
		if text == "" || banner == "" {
			data.Error = "Error: Text and banner fields cannot be empty."
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, data)
			return
		}

		// ASCII-only check
		for _, char := range text {
			if char > 127 {
				data.Error = "HTTP status 400 - Bad Request: Non-ASCII characters are not supported."
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, data)
				return
			}
		}

		// Generate ASCII
		asciiArt, err := ascii.GenerateASCIIArt(text, banner)
		if err != nil {
			log.Println("ASCII Generation Error:", err)
			data.Error = err.Error()

			if err.Error() == "Banner file not found" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			tmpl.Execute(w, data)
			return
		}

		// Store for exporting
		lastGeneratedASCII = asciiArt
		lastBanner = banner

		data.Text = text
		data.Banner = banner
		data.ASCIIArt = asciiArt
		w.WriteHeader(http.StatusOK)
	}

	tmpl.Execute(w, data)
}

// ----------------------------------------------------
// Export Handler (/download?format=txt|json|html)
// ----------------------------------------------------
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure we have ASCII art
	if lastGeneratedASCII == "" {
		http.Error(w, "No ASCII art generated yet", http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	switch format {
	case "json":
		exportJSON(w, lastGeneratedASCII, lastBanner)
	case "html":
		exportHTML(w, lastGeneratedASCII)
	default:
		exportTXT(w, lastGeneratedASCII)
	}
}

func exportTXT(w http.ResponseWriter, asciiArt string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(asciiArt)))
	w.Write([]byte(asciiArt))
}

func exportJSON(w http.ResponseWriter, asciiArt, banner string) {
	lines := strings.Split(asciiArt, "\n")
	data := map[string]interface{}{
		"ascii_art": lines,
		"banner":    banner,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonData, _ := json.MarshalIndent(data, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonData)))
	w.Write(jsonData)
}

func exportHTML(w http.ResponseWriter, asciiArt string) {
	htmlContent := fmt.Sprintf("<html><body><pre>%s</pre></body></html>", asciiArt)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.html")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(htmlContent)))
	w.Write([]byte(htmlContent))
}
