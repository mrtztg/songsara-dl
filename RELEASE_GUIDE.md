# GitHub Release Setup Guide

This guide will walk you through setting up automatic GitHub releases for the SongSara downloader.

## Step 1: Create GitHub Repository

1. Go to [GitHub](https://github.com) and create a new repository
2. Name it `songsara-dl`
3. Make it public or private (your choice)
4. **Don't** initialize with README, .gitignore, or license (we already have these)

## Step 2: Update Repository Configuration

### Update `.goreleaser.yml`

Replace `yourusername` in the file with your actual GitHub username:

```yaml
release:
  github:
    owner: YOUR_GITHUB_USERNAME  # Replace this
    name: songsara-dl
```

### Update `README.md`

Replace all instances of `yourusername` in the download links with your actual GitHub username:

```markdown
- **macOS**: [songsara-dl_Darwin_x86_64.tar.gz](https://github.com/YOUR_GITHUB_USERNAME/songsara-dl/releases/latest/download/songsara-dl_Darwin_x86_64.tar.gz)
```

## Step 3: Push to GitHub

```bash
# Add the remote repository (replace YOUR_GITHUB_USERNAME)
git remote add origin https://github.com/YOUR_GITHUB_USERNAME/songsara-dl.git

# Push to GitHub
git push -u origin master
```

## Step 4: Create Your First Release

### Option A: Using Git Tags (Recommended)

```bash
# Create and push a tag
git tag v1.0.0
git push origin v1.0.0
```

This will automatically trigger the GitHub Actions workflow and create a release.

### Option B: Manual Release

1. Go to your repository on GitHub
2. Click on "Releases" in the right sidebar
3. Click "Create a new release"
4. Choose a tag (e.g., `v1.0.0`)
5. Add a title and description
6. Upload the built binaries from the `dist/` folder

## Step 5: Verify the Release

1. Go to your repository's "Releases" page
2. You should see your release with all the platform binaries
3. Test downloading one of the binaries

## Step 6: Future Releases

For future releases, simply:

```bash
# Make your changes
git add .
git commit -m "New feature or fix"

# Create a new tag
git tag v1.1.0
git push origin v1.1.0
```

The GitHub Actions workflow will automatically:
- Build binaries for all platforms
- Create a new GitHub release
- Upload all the binaries
- Generate release notes

## Troubleshooting

### If GitHub Actions fails:

1. Check the Actions tab in your repository
2. Ensure the repository has the correct permissions
3. Verify your `.goreleaser.yml` configuration

### If binaries aren't uploaded:

1. Check that the GitHub token has the correct permissions
2. Verify the release configuration in `.goreleaser.yml`
3. Check the GoReleaser logs in the Actions tab

## Local Testing

Before pushing to GitHub, you can test locally:

```bash
# Test the build process
make snapshot

# Test with a local release (requires GITHUB_TOKEN)
export GITHUB_TOKEN=your_token_here
goreleaser release --snapshot --clean
```

## Release Notes

The release will automatically include:
- Changelog from git commits
- Download links for all platforms
- Release notes template

## Supported Platforms

Your release will include binaries for:
- macOS (Intel & Apple Silicon)
- Linux (x86_64 & ARM64)
- Windows (x86_64)

Each platform will have its own download link in the release. 