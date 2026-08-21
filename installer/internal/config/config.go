package config

import "path/filepath"

// Version is set at release via ldflags.
var Version = "1.0.0"

const (
	AppName = "cursor-utils"

	GitHubOwner = "imflawlezz"
	GitHubRepo  = "cursor-utils"

	ManifestFormatVersion = 1

	// Cursor ignores this directory; it is not commands, rules, skills, agents, or hooks.
	ManifestDir  = ".cursor-utils"
	ManifestFile = "manifest.json"
	BackupDir    = "backups"

	// Binary-only tags. Never offered as Cursor content; content stays vMAJOR.MINOR.PATCH.
	InstallerTagPrefix = "installer-v"
)

type App struct {
	Name        string
	Version     string
	GitHubOwner string
	GitHubRepo  string
}

func Default() App {
	return App{
		Name:        AppName,
		Version:     Version,
		GitHubOwner: GitHubOwner,
		GitHubRepo:  GitHubRepo,
	}
}

func (a App) RepoSlug() string {
	return a.GitHubOwner + "/" + a.GitHubRepo
}

func (a App) RepositoryURL() string {
	return "https://github.com/" + a.RepoSlug()
}

func (a App) AuthorURL() string {
	return "https://github.com/" + a.GitHubOwner
}

func ManifestPath(cursorRoot string) string {
	return filepath.Join(cursorRoot, ManifestDir, ManifestFile)
}

func ManifestDirPath(cursorRoot string) string {
	return filepath.Join(cursorRoot, ManifestDir)
}

func BackupRoot(cursorRoot string) string {
	return filepath.Join(cursorRoot, ManifestDir, BackupDir)
}
