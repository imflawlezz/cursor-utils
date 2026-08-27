# cursor-utils

Reusable Cursor development workflows as Markdown slash commands.

## Contents

| Repository path | Installed path |
| --- | --- |
| `commands/` | `~/.cursor/commands/` |

## Commands

Invoke installed commands through Cursor's `/` command interface.

### Git

- `/commit` — Review staged and unstaged changes, propose atomic Conventional
  Commits with explicit `git add` / `git commit` commands, and do not execute
  them.

### Changelog

- `/changelog` — Update only the `[Unreleased]` section of `CHANGELOG.md` from
  changes since the latest SemVer tag, following Keep a Changelog. Does not
  edit README or `docs/`.
- `/changelog-draft` — Propose an `[Unreleased]` section without modifying
  files.
- `/changelog-review` — Audit `CHANGELOG.md` for Keep a Changelog, SemVer, and
  consistency issues; report findings only.
- `/release` — Turn `[Unreleased]` into a dated SemVer version section and
  update changelog reference links; does not create tags, GitHub releases,
  commits, or pushes, and does not bump versions in other files.

### Code quality

- `/cleanup` — Remove or rewrite low-value comments in source without changing
  program behavior or rewriting project Markdown.

### Architecture

- `/clean-arch` — Review architectural boundaries against the project's actual
  structure and apply the smallest coherent refactor when a fix is warranted.
  Does not rewrite README, changelog, or `docs/`.

### Documentation

- `/readme` — Create or update only `README.md` from the current project state.
- `/readme-review` — Audit `README.md` against the implementation; report
  findings only.
- `/docs` — Maintain technical documentation under `docs/` (and equivalent
  trees). Does not edit `README.md` or `CHANGELOG.md`.
- `/docs-review` — Audit `docs/` against the implementation; report findings
  only.

## Requirements

- [Cursor](https://cursor.com) to invoke installed commands
- macOS, Linux, or Windows (`amd64` or `arm64`), an interactive terminal, and
  network access to GitHub to run the installer

## Installation

The supported installer is a standalone TUI. It does not require Go, Git, or
Node at runtime. To build from source, see
[installer/README.md](installer/README.md).

Download a release archive for your OS and architecture from the latest
[`installer-v*`](https://github.com/imflawlezz/cursor-utils/releases) GitHub
Release, extract it, then run the binary:

```bash
tar -xzf cursor-utils-installer_1.2.0_darwin_arm64.tar.gz
chmod +x cursor-utils
./cursor-utils
```

Windows archives are `.zip` (`windows_amd64` or `windows_arm64`). Replace
`darwin_arm64` with `darwin_amd64`, `linux_arm64`, `linux_amd64`,
`windows_amd64`, or `windows_arm64`. Release assets are also listed in
`checksums.txt` on the same GitHub Release.
The installer copies tagged `commands/` files into your Cursor configuration
directory (default `~/.cursor/commands/` on macOS and Linux,
`%USERPROFILE%\.cursor\commands` on Windows) and records ownership in a
manifest so later updates and uninstalls do not touch your own files.

You can still copy files by hand if you prefer:

```bash
mkdir -p ~/.cursor/commands
cp /path/to/cursor-utils/commands/*.md ~/.cursor/commands/
```

On Windows, copy the `.md` files into `%USERPROFILE%\.cursor\commands`. The
repository does not need to live under the Cursor directory.

| Role | Path |
| --- | --- |
| Source | `cursor-utils/commands/*.md` |
| Installed | `~/.cursor/commands/*.md` |
| Installer manifest | `~/.cursor/.cursor-utils/manifest.json` |

On Windows, replace `~/.cursor` with `%USERPROFILE%\.cursor`.

## Cursor compatibility

**Target:** Cursor 3.x

**Developed / tested against:** Cursor 3.17.8 (Universal)

Cursor's command format and customization APIs may change between releases.
This toolkit does not claim compatibility with untested Cursor versions.
After upgrading Cursor—especially across a major version—verify that Markdown
commands in `~/.cursor/commands/` are still discovered and invoked as expected.

Official documentation:

- [Customize Cursor](https://cursor.com/docs/customize-cursor)
- [Plugins reference (commands format)](https://cursor.com/docs/reference/plugins)

## Usage

Open Agent chat in Cursor and run a command by name, for example:

```text
/commit
/changelog
/readme
```

Command behavior is defined by the corresponding Markdown file under
`commands/`.

## Updating

Use the installer and select a newer content version. It writes the files you
select from that version. Files it does not own are left unchanged unless you
choose to overwrite them.

To update by hand, copy the updated `commands/*.md` files into
`~/.cursor/commands/` again, overwriting the previous copies if desired.

## Uninstalling

Use the installer, select the files to remove, and choose **Remove selected**,
or delete only the Markdown files you copied into `~/.cursor/commands/`.

Do not remove unrelated commands in that directory. The installer never
deletes the `commands` directory itself or files it does not own.

## Repository layout

```text
cursor-utils/
├── commands/             # Slash command Markdown files
├── installer/            # Standalone TUI installer (Go)
│   └── README.md
├── README.md
├── CHANGELOG.md
└── LICENSE
```

## Contributing

Commands live in `commands/` as Markdown files. The basename (without `.md`) is
the slash command name.

When adding or changing a command:

1. Keep the prompt factual and specific to the intended workflow.
2. State a **Scope** of files the command may modify or review, and files it
   must not touch. Recommend a handoff to the owning command instead of
   overlapping.
3. Do not claim behavior the command does not implement.
4. Update `CHANGELOG.md` under `[Unreleased]`.

## Versioning

This project follows [Semantic Versioning](https://semver.org/) with two
independent series:

| What | Tag | Example |
| --- | --- | --- |
| Cursor content (`commands/`, …) | `vMAJOR.MINOR.PATCH` | `v0.1.0` |
| Installer binary | `installer-vMAJOR.MINOR.PATCH` | `installer-v1.0.0` |

The TUI installer only installs content tags. Tag an installer-only release as
`installer-v1.0.0`, not `v1.0.0`.

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## License

[MIT](LICENSE)
