package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/html"
)

func (a *Application) sendText(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := a.bot.Send(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

// Download media from URL and send to Telegram
func (a *Application) downloadAndSendMedia(chatId int64, url string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Download the file
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download media: %v", err)
	}
	defer resp.Body.Close()

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download media: status %d", resp.StatusCode)
	}

	// Read the content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read media content: %v", err)
	}

	// Create FileBytes from the downloaded content
	fileBytes := tgbotapi.FileBytes{
		Name:  getFilenameFromURL(url),
		Bytes: content,
	}

	// Determine media type and send accordingly
	if isImageURL(url) {
		photo := tgbotapi.NewPhoto(chatId, fileBytes)
		_, err = a.bot.Send(photo)
	} else if isVideoURL(url) {
		video := tgbotapi.NewVideo(chatId, fileBytes)
		_, err = a.bot.Send(video)
	} else {
		// Try sending as document if type is unknown
		document := tgbotapi.NewDocument(chatId, fileBytes)
		_, err = a.bot.Send(document)
	}

	return err
}

// Extract filename from URL
func getFilenameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]
		// Remove query parameters if any
		if idx := strings.Index(filename, "?"); idx != -1 {
			filename = filename[:idx]
		}
		return filename
	}
	return "media_file"
}

// Helper function to check if URL is an image
func isImageURL(url string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}
	lowerURL := strings.ToLower(url)
	for _, ext := range imageExtensions {
		if strings.HasSuffix(lowerURL, ext) {
			return true
		}
	}
	return false
}

// Helper function to check if URL is a video
func isVideoURL(url string) bool {
	videoExtensions := []string{".mp4", ".avi", ".mov", ".webm", ".mkv", ".flv"}
	lowerURL := strings.ToLower(url)
	for _, ext := range videoExtensions {
		if strings.HasSuffix(lowerURL, ext) {
			return true
		}
	}
	return false
}

// fetchHTML fetches HTML content from the given URL and returns it as a string.
func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// extractPreviewReddItLinks parses HTML and returns all <a> hrefs with origin "preview.redd.it"
func extractPreviewReddItLinks(r io.Reader) ([]string, error) {
	var links []string
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "video" || n.Data == "img") {
			var shouldFetch = false
			for _, attr := range n.Attr {
				if attr.Key == "fetchpriority" {
					shouldFetch = true
				}
			}
			if shouldFetch {
				for _, attr := range n.Attr {
					if attr.Key == "src" && strings.HasPrefix(attr.Val, "https://preview.redd.it") {
						if strings.Contains(attr.Val, "-") {
							links = append(links, attr.Val)
						}
					}

				}
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, nil
}

func removeDuplicateLinks(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func getRandomItem(slice []string) string {
	if len(slice) == 0 {
		return ""
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Get random index
	randomIndex := rand.Intn(len(slice))
	return slice[randomIndex]
}
