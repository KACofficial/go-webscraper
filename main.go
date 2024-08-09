package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"webscrape/utils"

	"github.com/gocolly/colly"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a url: https://")
	BaseURLNoHTTP, _ := reader.ReadString('\n')
	BaseURLNoHTTP = BaseURLNoHTTP[:len(BaseURLNoHTTP)-1]
	fmt.Print(BaseURLNoHTTP)
	BaseURL := "https://" + BaseURLNoHTTP
	// BaseURL = BaseURL[:len(BaseURL)-1]

	c := colly.NewCollector()
	c.IgnoreRobotsTxt = true

	fmt.Printf("Would you like it to ONLY visit %s? (y/n): ", BaseURLNoHTTP)
	resp, _ := reader.ReadString('\n')
	resp = strings.TrimSpace(strings.ToLower(resp))

	if strings.ToLower(resp) == "y" {
		c.AllowedDomains = append(c.AllowedDomains, "github.com")
	}

	utils.SetupScrapeEnv(BaseURL)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		imgSrc, _ := utils.ResolveURL(BaseURL, e.Attr("src"))
		utils.DownloadImage(imgSrc)
		// fmt.Println("Found img:", imgSrc)
		//return
	})

	c.OnHTML("h1, h2, h3, h4, h5, h6, p", func(e *colly.HTMLElement) {
		utils.AddToTextStruct(e.Name, e.Text)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		absURL := href
		if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "?") || strings.HasPrefix(href, "#") {
			absURL = utils.JoinURL(BaseURL, href)
		}
		fmt.Println("Checking", absURL)
		e.Request.Visit(absURL)
		//return
	})

	go func() {
		<-sigCh // Wait for an interrupt signal
		fmt.Println("Received interrupt signal. Exiting...")
		utils.SaveTextStructJSONFile()
		os.Exit(0) // Exit the program
	}()

	c.Visit(BaseURL)
}
