#!/bin/bash

# SongSara Downloader - Usage Examples
# This script demonstrates various ways to use the songsara-dl tool

echo "SongSara Downloader - Usage Examples"
echo "===================================="

# Check if the binary exists
if [ ! -f "./songsara-dl" ]; then
    echo "Error: songsara-dl binary not found. Please run 'make build' first."
    exit 1
fi

echo ""
echo "1. Show help information:"
echo "   ./songsara-dl --help"
echo ""

echo "2. Download a single album (basic usage):"
echo "   ./songsara-dl \"https://songsara.net/album-url\""
echo ""

echo "3. Download with custom concurrency (5 concurrent downloads):"
echo "   ./songsara-dl -c 5 \"https://songsara.net/album-url\""
echo ""

echo "4. Download to custom output directory:"
echo "   ./songsara-dl -o /path/to/music \"https://songsara.net/album-url\""
echo ""

echo "5. Enable verbose output:"
echo "   ./songsara-dl -v \"https://songsara.net/album-url\""
echo ""

echo "6. Download multiple albums:"
echo "   ./songsara-dl \"https://songsara.net/album1\" \"https://songsara.net/album2\""
echo ""

echo "7. Combine all options:"
echo "   ./songsara-dl -c 8 -o /music -v \"https://songsara.net/album-url\""
echo ""

echo "8. Dry run to see what would be downloaded:"
echo "   ./songsara-dl -n -v \"https://songsara.net/album-url\""
echo ""

echo "9. Download with custom timeout:"
echo "   ./songsara-dl -t 60 \"https://songsara.net/album-url\""
echo ""

echo "10. Test with help (no actual download):"
echo "    ./songsara-dl --help"
echo ""

echo "Note: Replace 'https://songsara.net/album-url' with actual SongSara URLs"
echo "The tool will create a 'downloads/' directory (or custom directory) with organized folders." 