package ascii

import (
	"net/http"
	"os"
)

// 404 handler
func Handler404(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data, err := os.ReadFile("templates/404.html")
	if err != nil {
		w.Write([]byte("404 Not Found"))
		return
	}
	w.Write(data)
}
