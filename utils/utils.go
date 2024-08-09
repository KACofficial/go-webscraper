package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var scrapingUrl string
var scrapingUrlFolder string
var foundImages []string

type HTMLText struct {
	H1 []string `json:"h1"`
	H2 []string `json:"h2"`
	H3 []string `json:"h3"`
	H4 []string `json:"h4"`
	H5 []string `json:"h5"`
	H6 []string `json:"h6"`
	P  []string `json:"p"`
}

var completeHTMLText HTMLText = HTMLText{
	H1: []string{},
	H2: []string{},
	H3: []string{},
	H4: []string{},
	H5: []string{},
	H6: []string{},
	P:  []string{},
}

func SetupScrapeEnv(url string) {
	scrapingUrl = url
	dirname, err := extractMainDomain(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	scrapingUrlFolder = dirname
	folders := filepath.Join(dirname, "images")
	err = os.MkdirAll(folders, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Made path:", folders)
}

func extractMainDomain(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()
	parts := strings.Split(host, ".")

	// If the domain has at least two parts, return the second-to-last part
	if len(parts) >= 2 {
		return parts[len(parts)-2], nil
	}

	// If the domain is malformed or too short, return an empty string or an error
	return "", fmt.Errorf("invalid domain: %s", rawURL)
}

func ResolveURL(baseURL, imgSrc string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Check if the imgSrc is already an absolute URL
	imgURL, err := url.Parse(imgSrc)
	if err != nil {
		return "", err
	}

	// Resolve the imgSrc relative to the baseURL
	resolvedURL := base.ResolveReference(imgURL)

	return resolvedURL.String(), nil
}

func DownloadImage(src string) {
	src = strings.Split(src, "?")[0]

	if contains(foundImages, src) {
		return
	}

	fmt.Println("Downloading", src)

	foundImages = append(foundImages, src)

	downloadFile(filepath.Join(scrapingUrlFolder, "images", path.Base(src)), src)

}

func AddToTextStruct(element string, value string) {
	switch element {
	case "h1":
		if contains(completeHTMLText.H1, value) {
			return
		}
		completeHTMLText.H1 = append(completeHTMLText.H1, value)
	case "h2":
		if contains(completeHTMLText.H2, value) {
			return
		}
		completeHTMLText.H2 = append(completeHTMLText.H2, value)
	case "h3":
		if contains(completeHTMLText.H3, value) {
			return
		}
		completeHTMLText.H3 = append(completeHTMLText.H3, value)
	case "h4":
		if contains(completeHTMLText.H4, value) {
			return
		}
		completeHTMLText.H4 = append(completeHTMLText.H4, value)
	case "h5":
		if contains(completeHTMLText.H5, value) {
			return
		}
		completeHTMLText.H5 = append(completeHTMLText.H5, value)
	case "h6":
		if contains(completeHTMLText.H6, value) {
			return
		}
		completeHTMLText.H6 = append(completeHTMLText.H6, value)
	case "p":
		if contains(completeHTMLText.P, value) {
			return
		}
		completeHTMLText.P = append(completeHTMLText.P, value)
	}
}

func SaveTextStructJSONFile() {
	file, err := os.Create(filepath.Join(scrapingUrlFolder, "text.json"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(completeHTMLText)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func downloadFile(filepath string, url string) error {
	// Get the data
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}
