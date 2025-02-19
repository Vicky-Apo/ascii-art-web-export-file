package ascii

import (
	"fmt"
	"net/http"
)

func ExportTXT(w http.ResponseWriter, asciiArt string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(asciiArt)))
	w.Write([]byte(asciiArt))
}

func ExportHTML(w http.ResponseWriter, asciiArt string) {
	htmlContent := fmt.Sprintf("<html><body><pre>%s</pre></body></html>", asciiArt)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.html")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(htmlContent)))
	w.Write([]byte(htmlContent))
}
