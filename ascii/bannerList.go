package ascii

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func BannerList() ([]string, error) {
	bannerFiles, err := os.ReadDir("./banners")
	if err != nil {
		log.Printf("Could not read banners directory %v", err)
		return nil, err
	}

	var banners []string
	for _, entry := range bannerFiles {
		if !entry.IsDir() {
			name := entry.Name()
			ext := filepath.Ext(name)
			baseName := strings.TrimSuffix(name, ext)
			banners = append(banners, baseName)
		}
	}
	return banners, nil
}
