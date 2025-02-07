package main

import (
	"ascii-art-web/ascii"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	Text     string
	Banner   string
	ASCIIArt string
	Error    string
	Banners  []string
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/404", handler404)

	fmt.Println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handler404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data, err := os.ReadFile("templates/404.html")
	if err != nil {
		w.Write([]byte("404 Not Found"))
		return
	}
	w.Write(data)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Template Parsing Error:", err)
		return
	}

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

		if text == "" || banner == "" {
			data.Error = "Error: Text and banner fields cannot be empty."
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, data)
			return
		}

		for _, char := range text {
			if char > 127 { // ASCII range check
				data.Error = "HTTP status 400 - Bad Request: Non-ASCII characters are not supported."
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, data)
				return
			}
		}

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

		data.Text = text
		data.Banner = banner
		data.ASCIIArt = asciiArt
		w.WriteHeader(http.StatusOK)
	}

	tmpl.Execute(w, data)
}
