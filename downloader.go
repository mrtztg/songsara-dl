package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
)

type Song struct {
	Title string
	URL   string
}

type Album struct {
	Title string
	Songs []Song
}

type SongSaraDownloader struct {
	client       *http.Client
	concurrency  int
	outputDir    string
	verbose      bool
	dryRun       bool
	skipExisting bool
	timeout      int
}

func NewSongSaraDownloader() *SongSaraDownloader {
	return &SongSaraDownloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DisableCompression: false, // Enable compression handling
			},
		},
		concurrency:  10,
		outputDir:    "downloads",
		verbose:      false,
		dryRun:       false,
		skipExisting: true,
		timeout:      30,
	}
}

func (d *SongSaraDownloader) Download(urls []string) error {
	// Create output directory
	if err := os.MkdirAll(d.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	totalURLs := len(urls)
	successfulDownloads := 0
	failedDownloads := 0

	fmt.Printf("Starting download of %d album(s)/playlist(s)...\n", totalURLs)

	// Process each URL
	for i, url := range urls {
		if totalURLs > 1 {
			fmt.Printf("\n[%d/%d] Processing: %s\n", i+1, totalURLs, url)
		} else if d.verbose {
			fmt.Printf("Processing URL: %s\n", url)
		}

		album, err := d.scrapeAlbum(url)
		if err != nil {
			fmt.Printf("❌ Error scraping %s: %v\n", url, err)
			failedDownloads++
			continue
		}

		if err := d.downloadAlbum(album); err != nil {
			fmt.Printf("❌ Error downloading album %s: %v\n", album.Title, err)
			failedDownloads++
			continue
		}

		successfulDownloads++
	}

	// Print summary
	fmt.Printf("\n" + strings.Repeat("=", 50) + "\n")
	fmt.Printf("Download Summary:\n")
	fmt.Printf("  Total URLs: %d\n", totalURLs)
	fmt.Printf("  Successful: %d\n", successfulDownloads)
	fmt.Printf("  Failed: %d\n", failedDownloads)
	fmt.Printf("  Output directory: %s\n", d.outputDir)
	fmt.Printf(strings.Repeat("=", 50) + "\n")

	if failedDownloads > 0 {
		return fmt.Errorf("%d download(s) failed", failedDownloads)
	}

	return nil
}

func (d *SongSaraDownloader) scrapeAlbum(pageURL string) (*Album, error) {
	// Create request with proper headers
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Make request
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract album title
	albumTitle := d.extractAlbumTitle(doc)
	if albumTitle == "" {
		albumTitle = "Unknown Album"
	}

	if d.verbose {
		fmt.Printf("Found album title: %s\n", albumTitle)
	}

	// Extract songs
	songs := d.extractSongs(doc)

	if d.verbose {
		fmt.Printf("Found %d songs\n", len(songs))
		if len(songs) == 0 {
			fmt.Println("No songs found. Trying to debug...")
			// Print some debug info
			fmt.Printf("Page title: %s\n", doc.Find("title").Text())
			fmt.Printf("H1 elements: %d\n", doc.Find("h1").Length())
			fmt.Printf("Audio elements: %d\n", doc.Find("audio").Length())
			fmt.Printf("Links with audio extensions: %d\n", doc.Find("a[href*='.mp3'], a[href*='.m4a'], a[href*='.wav'], a[href*='.flac']").Length())

			// Print first 1000 characters of the HTML to see what we're getting
			html, _ := doc.Html()
			if len(html) > 1000 {
				fmt.Printf("First 1000 characters of HTML:\n%s\n", html[:1000])
			} else {
				fmt.Printf("Full HTML (length: %d):\n%s\n", len(html), html)
			}

			// Check if we're getting a blocked page or error page
			if strings.Contains(html, "blocked") || strings.Contains(html, "captcha") || strings.Contains(html, "403") || strings.Contains(html, "404") {
				fmt.Println("WARNING: Page appears to be blocked or showing an error!")
			}
		}
	}

	return &Album{
		Title: albumTitle,
		Songs: songs,
	}, nil
}

func (d *SongSaraDownloader) extractAlbumTitle(doc *goquery.Document) string {
	// Try the selector from .cursorrules
	title := doc.Find(".AL-Si").Text()
	if title != "" {
		return strings.TrimSpace(title)
	}

	// Try h1 tags (common for page titles)
	title = doc.Find("h1").First().Text()
	if title != "" {
		return strings.TrimSpace(title)
	}

	// Try to find title in meta tags
	title = doc.Find("meta[property='og:title']").AttrOr("content", "")
	if title != "" {
		return strings.TrimSpace(title)
	}

	// Try to find title in page title
	title = doc.Find("title").Text()
	if title != "" {
		return strings.TrimSpace(title)
	}

	// Fallback selectors
	selectors := []string{
		".title",
		".album-title",
		"[class*='title']",
		"[class*='album']",
		".playlist-title",
		"[class*='playlist']",
	}

	for _, selector := range selectors {
		title = doc.Find(selector).First().Text()
		if title != "" {
			return strings.TrimSpace(title)
		}
	}

	return ""
}

func (d *SongSaraDownloader) extractSongFromAudioElement(s *goquery.Selection) Song {
	// Extract title from various possible sources
	title := s.AttrOr("title", "")
	if title == "" {
		title = s.AttrOr("alt", "")
	}
	if title == "" {
		// Try to find title in parent or sibling elements
		title = s.Parent().Find("[title]").AttrOr("title", "")
	}
	if title == "" {
		title = s.Parent().Text()
	}

	// Extract URL from src attribute
	url := s.AttrOr("src", "")
	if url == "" {
		// Try to find source in child elements
		url = s.Find("source").AttrOr("src", "")
	}

	return Song{
		Title: strings.TrimSpace(title),
		URL:   strings.TrimSpace(url),
	}
}

func (d *SongSaraDownloader) extractSongFromLink(s *goquery.Selection) Song {
	// Extract title from link text or title attribute
	title := s.Text()
	if title == "" {
		title = s.AttrOr("title", "")
	}
	if title == "" {
		title = s.AttrOr("alt", "")
	}

	// Extract URL from href attribute
	url := s.AttrOr("href", "")

	return Song{
		Title: strings.TrimSpace(title),
		URL:   strings.TrimSpace(url),
	}
}

func (d *SongSaraDownloader) extractSongs(doc *goquery.Document) []Song {
	var songs []Song

	// Try the selector from .cursorrules
	doc.Find("#aramplayer .audioplayer-audios li").Each(func(i int, s *goquery.Selection) {
		song := d.extractSongFromElement(s)
		if song.Title != "" && song.URL != "" {
			songs = append(songs, song)
		}
	})

	// If no songs found, try alternative selectors
	if len(songs) == 0 {
		selectors := []string{
			"li[data-title]",
			".track",
			".song",
			"[class*='track']",
			"[class*='song']",
			// Try to find any list items that might contain songs
			"li",
			".playlist-item",
			"[class*='playlist']",
			// Look for any elements with audio-related attributes
			"[data-src]",
			"[src*='.mp3']",
			"[src*='.m4a']",
			"[src*='.wav']",
			"[src*='.flac']",
		}

		for _, selector := range selectors {
			doc.Find(selector).Each(func(i int, s *goquery.Selection) {
				song := d.extractSongFromElement(s)
				if song.Title != "" && song.URL != "" {
					songs = append(songs, song)
				}
			})
			if len(songs) > 0 {
				break
			}
		}
	}

	// If still no songs found, try to extract from any audio elements
	if len(songs) == 0 {
		doc.Find("audio").Each(func(i int, s *goquery.Selection) {
			song := d.extractSongFromAudioElement(s)
			if song.Title != "" && song.URL != "" {
				songs = append(songs, song)
			}
		})
	}

	// If still no songs, try to find any links that might be download links
	if len(songs) == 0 {
		doc.Find("a[href*='.mp3'], a[href*='.m4a'], a[href*='.wav'], a[href*='.flac']").Each(func(i int, s *goquery.Selection) {
			song := d.extractSongFromLink(s)
			if song.Title != "" && song.URL != "" {
				songs = append(songs, song)
			}
		})
	}

	return songs
}

func (d *SongSaraDownloader) extractSongFromElement(s *goquery.Selection) Song {
	// Extract title from data-title attribute
	title := s.AttrOr("data-title", "")
	if title == "" {
		// Try to find title in child elements
		title = s.Find("[data-title]").AttrOr("data-title", "")
	}
	if title == "" {
		// Try to find any text content
		title = strings.TrimSpace(s.Text())
	}

	// Extract download URL
	var downloadURL string

	// Try the selector from .cursorrules
	sourceDiv := s.Find("div.audioplayer-source")
	if sourceDiv.Length() > 0 {
		downloadURL = sourceDiv.AttrOr("data-src", "")
	}

	// If not found, try alternative selectors
	if downloadURL == "" {
		selectors := []string{
			"[data-src]",
			"[src]",
			"audio source",
			"a[href*='.mp3']",
			"a[href*='.m4a']",
			"a[href*='.wav']",
		}

		for _, selector := range selectors {
			element := s.Find(selector).First()
			if element.Length() > 0 {
				downloadURL = element.AttrOr("data-src", element.AttrOr("src", element.AttrOr("href", "")))
				if downloadURL != "" {
					break
				}
			}
		}
	}

	return Song{
		Title: strings.TrimSpace(title),
		URL:   strings.TrimSpace(downloadURL),
	}
}

func (d *SongSaraDownloader) downloadAlbum(album *Album) error {
	if len(album.Songs) == 0 {
		return fmt.Errorf("no songs found in album")
	}

	// Create album directory
	albumDir := filepath.Join(d.outputDir, d.sanitizeFilename(album.Title))
	if err := os.MkdirAll(albumDir, 0755); err != nil {
		return fmt.Errorf("failed to create album directory: %w", err)
	}

	if d.dryRun {
		fmt.Printf("Dry run - would download album: %s (%d songs)\n", album.Title, len(album.Songs))
		return nil
	}

	fmt.Printf("Downloading album: %s (%d songs)\n", album.Title, len(album.Songs))

	// Create progress bar
	bar := progressbar.Default(int64(len(album.Songs)), "Downloading songs")

	// Create semaphore for concurrency control
	semaphore := make(chan struct{}, d.concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	// Download songs concurrently
	for i, song := range album.Songs {
		wg.Add(1)
		go func(index int, s Song) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			if err := d.downloadSong(s, albumDir, index+1); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("song '%s': %w", s.Title, err))
				mu.Unlock()
			}

			bar.Add(1)
		}(i, song)
	}

	wg.Wait()
	bar.Finish()

	// Report errors
	if len(errors) > 0 {
		fmt.Printf("\nErrors occurred during download:\n")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Printf("Album '%s' download completed!\n", album.Title)
	return nil
}

func (d *SongSaraDownloader) downloadSong(song Song, albumDir string, trackNumber int) error {
	if song.URL == "" {
		return fmt.Errorf("no download URL found")
	}

	// Determine file extension
	ext := d.getFileExtension(song.URL)
	if ext == "" {
		ext = ".mp3" // Default extension
	}

	// Create filename
	filename := fmt.Sprintf("%02d - %s%s", trackNumber, d.sanitizeFilename(song.Title), ext)
	filepath := filepath.Join(albumDir, filename)

	// Skip if file already exists
	if d.skipExisting {
		if _, err := os.Stat(filepath); err == nil {
			if d.verbose {
				fmt.Printf("Skipping existing file: %s\n", filename)
			}
			return nil
		}
	}

	// If dry run, just show what would be downloaded
	if d.dryRun {
		if d.verbose {
			fmt.Printf("Would download: %s\n", filename)
		}
		return nil
	}

	// Download file
	resp, err := d.client.Get(song.URL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	// Create file
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if d.verbose {
		fmt.Printf("Downloaded: %s\n", filename)
	}

	return nil
}

func (d *SongSaraDownloader) sanitizeFilename(filename string) string {
	// Remove filesystem-incompatible characters
	re := regexp.MustCompile(`[/*?:"<>|]`)
	filename = re.ReplaceAllString(filename, "")

	// Remove extra whitespace
	filename = strings.TrimSpace(filename)

	// Replace multiple spaces with single space
	re = regexp.MustCompile(`\s+`)
	filename = re.ReplaceAllString(filename, " ")

	return filename
}

func (d *SongSaraDownloader) getFileExtension(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	path := parsedURL.Path
	if strings.HasSuffix(path, ".mp3") {
		return ".mp3"
	} else if strings.HasSuffix(path, ".m4a") {
		return ".m4a"
	} else if strings.HasSuffix(path, ".wav") {
		return ".wav"
	} else if strings.HasSuffix(path, ".flac") {
		return ".flac"
	}

	return ""
}
