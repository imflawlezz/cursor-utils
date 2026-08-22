# Installer

Go module for the cursor-utils TUI. It fetches a GitHub **content** tag tarball
and copies installer-owned files into a Cursor root. Users should download a
release binary; this file is for building and changing the installer.

Install and command docs: [../README.md](../README.md).

## Build

Requires [Go](https://go.dev/) 1.25. `CGO_ENABLED=0`. From this directory:

| Target | Output |
| --- | --- |
| `make build` | `dist/cursor_utils_installer_<os>_<arch>_<version>` for the host (Windows adds `.exe`) |
| `make build-all` | darwin/linux/windows × amd64/arm64 under `dist/` |
| `make test` | unit tests |
| `make tidy` | `go mod tidy` |
| `make clean` | remove `dist/` |

Version comes from `VERSION` (default `1.1.0`) via ldflags. Dots in the
filename become hyphens (`1.1.0` → `1-1-0`).

```bash
make build
./dist/cursor_utils_installer_$(go env GOOS)_$(go env GOARCH)_1-1-0$(go env GOEXE)
```

Live GitHub test (optional):

```bash
CURSOR_UTILS_LIVE=1 go test ./internal/installer -run TestLiveGitHubInstallUninstall
```

## CLI

Subcommands are the same regardless of the binary filename:

```text
<binary>           Start the interactive installer
<binary> version   Print installer version
<binary> help      Show help
```

```bash
./cursor_utils_installer_darwin_arm64_1-1-0 version
```

`version` also accepts `--version` and `-v`. `help` also accepts `--help` and `-h`.
Needs an interactive TTY.

## Tags

Content vs installer tags are defined in [../README.md](../README.md#versioning).
The TUI lists `v*` tags only and downloads
`https://github.com/imflawlezz/cursor-utils/archive/refs/tags/<tag>.tar.gz`.
Do not tag an installer-only release as `v1.0.0`.

## Layout

```text
cmd/cursor-utils     entrypoint
internal/config      version, GitHub repo, manifest paths
internal/github      tag list and tarball fetch
internal/content     component registry and path safety
internal/installer   plan, apply, rollback
internal/manifest    ownership records
internal/platform    OS/arch and Cursor-root checks
internal/prefs       saved Cursor root (best-effort)
internal/semver      content vs installer tags
internal/tui         Bubble Tea UI
```

The engine writes installed files, backups, and the manifest. The TUI
best-effort saves the chosen Cursor root to `~/.cursor-utils.json` (save
errors are ignored). `?` in the TUI lists keys.

## Ownership

Written files are listed in `<cursor-root>/.cursor-utils/manifest.json`.
Updates and removes use that list only. The installer does not delete a
component directory or files it does not own. Overwrites are backed up under
`<cursor-root>/.cursor-utils/backups/`.
