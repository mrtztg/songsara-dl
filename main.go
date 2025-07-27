package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// Version information - set by GoReleaser
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "songsara-dl [URL1] [URL2] [URL3]...",
	Short: "Download songs from SongSara music platform",
	Long: `A CLI tool to download entire albums or playlists from SongSara.
Supports concurrent downloads with a maximum of 10 concurrent downloads at a time.

Examples:
  # Download a single album
  songsara-dl "https://songsara.net/59021/"

  # Download multiple albums
  songsara-dl "https://songsara.net/59021/" "https://songsara.net/12345/"

  # Download with custom concurrency
  songsara-dl -c 5 "https://songsara.net/59021/" "https://songsara.net/12345/"

  # Download to custom directory with verbose output
  songsara-dl -o /path/to/music -v "https://songsara.net/59021/"`,
	Example: `  songsara-dl "https://songsara.net/59021/"
  songsara-dl -c 8 -o /music -v "https://songsara.net/59021/" "https://songsara.net/12345/"`,
	Args: cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please provide at least one SongSara URL")
		}

		downloader := NewSongSaraDownloader()

		// Get flags and set them on the downloader
		if concurrency, err := cmd.Flags().GetInt("concurrency"); err == nil {
			downloader.concurrency = concurrency
		}
		if outputDir, err := cmd.Flags().GetString("output"); err == nil {
			downloader.outputDir = outputDir
		}
		if verbose, err := cmd.Flags().GetBool("verbose"); err == nil {
			downloader.verbose = verbose
		}
		if dryRun, err := cmd.Flags().GetBool("dry-run"); err == nil {
			downloader.dryRun = dryRun
		}
		if skipExisting, err := cmd.Flags().GetBool("skip-existing"); err == nil {
			downloader.skipExisting = skipExisting
		}
		if timeout, err := cmd.Flags().GetInt("timeout"); err == nil {
			downloader.timeout = timeout
			downloader.client.Timeout = time.Duration(timeout) * time.Second
		}

		// Show summary of what will be downloaded
		if len(args) > 1 {
			fmt.Printf("Will download %d albums/playlists:\n", len(args))
			for i, url := range args {
				fmt.Printf("  %d. %s\n", i+1, url)
			}
			fmt.Println()
		}

		return downloader.Download(args)
	},
}

func init() {
	rootCmd.Flags().IntP("concurrency", "c", 10, "Maximum number of concurrent downloads")
	rootCmd.Flags().StringP("output", "o", "downloads", "Output directory for downloaded files")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolP("dry-run", "n", false, "Show what would be downloaded without actually downloading")
	rootCmd.Flags().BoolP("skip-existing", "s", true, "Skip existing files (default: true)")
	rootCmd.Flags().IntP("timeout", "t", 30, "HTTP timeout in seconds")

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("songsara-dl version %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Date: %s\n", date)
		},
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
