package ascii

import (
	"fmt"
	"os"
	"strings"
)

func GenerateASCIIArt(input, banner string) (string, error) {
	fileName := ""
	switch banner {
	case "standard":
		fileName = "./banners/standard.txt"
	case "shadow":
		fileName = "./banners/shadow.txt"
	case "thinkertoy":
		fileName = "./banners/thinkertoy.txt"
	default:
		return "", fmt.Errorf("500: Invalid banner selected")
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("error reading banner file: %v", err)
	}

	normalizedContent := strings.ReplaceAll(string(content), "\r", "")
	fontLines := strings.Split(normalizedContent, "\n")
	inputLines := strings.Split(input, "\n")

	var result strings.Builder
	pendingEmptyLine := false
	for _, line := range inputLines {
		if line == "" {
			if !pendingEmptyLine {
				result.WriteString("\n")
				pendingEmptyLine = true
			}
			continue
		}
		pendingEmptyLine = false
		asciiArt := convertToASCIIArt(line, fontLines)
		result.WriteString(asciiArt + "\n")
	}

	return result.String(), nil
}

func convertToASCIIArt(text string, fontLines []string) string {
	if text == "" {
		return ""
	}
	var result strings.Builder
	for row := 1; row <= 8; row++ {
		lineOutput := ""
		for _, char := range text {
			if char < 32 || char > 126 {
				continue // Skip non-printable characters
			}
			asciiIndex := (int(char)-32)*9 + row
			if asciiIndex < len(fontLines) {
				lineOutput += fontLines[asciiIndex]
			}
		}
		result.WriteString(lineOutput + "\n")
	}
	return result.String()
}
