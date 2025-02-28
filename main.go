package main

import (
	"ascii-art-web-export-file/ascii"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	indexTmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Println("Template Parsing Error:", err)
		return
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ascii.HomePage(w, r, indexTmpl)
	})
	http.HandleFunc("/download", ascii.DownloadHandler)

	fmt.Println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
