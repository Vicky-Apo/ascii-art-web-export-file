package ascii

import (
	"html/template"
	"log"
	"net/http"
	"os"
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

// HomePage handler
func HomePage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {

	// Get banners from bannerList.go
	banners, err := BannerList()
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
		if r.URL.Path != "/" {
			handler404(w, r)
			return
		}
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
		asciiArt, err := GenerateASCIIArt(text, banner)
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

		data.Text = text
		data.Banner = banner
		data.ASCIIArt = asciiArt
		w.WriteHeader(http.StatusOK)
	}

	tmpl.Execute(w, data)
}

// 404 handler
func handler404(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data, err := os.ReadFile("templates/404.html")
	if err != nil {
		w.Write([]byte("404 Not Found"))
		return
	}
	w.Write(data)
}

// Export Handler
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure we have ASCII art
	if lastGeneratedASCII == "" {
		http.Error(w, "No ASCII art generated yet", http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	switch format {
	case "html":
		ExportHTML(w, lastGeneratedASCII)
	default:
		ExportTXT(w, lastGeneratedASCII)
	}
}
