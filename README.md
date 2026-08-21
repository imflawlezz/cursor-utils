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
  changes since the latest SemVer tag, following Keep a Changelog.
- `/changelog-draft` — Propose an `[Unreleased]` section without modifying
  files.
- `/changelog-review` — Audit `CHANGELOG.md` for Keep a Changelog, SemVer, and
  consistency issues; report findings only.
- `/release` — Turn `[Unreleased]` into a dated SemVer version section and
  update reference links; does not create tags, GitHub releases, commits, or
  pushes.

### Code quality

- `/cleanup` — Remove or rewrite low-value comments (including AI-style
  narration) without changing program behavior.

### Architecture

- `/clean-arch` — Review architectural boundaries against the project's actual
  structure and apply the smallest coherent refactor when a fix is warranted.

### Documentation

- `/readme` — Create or update `README.md` from the current project state.
- `/readme-review` — Audit `README.md` against the implementation; report
  findings only.
- `/docs` — Maintain technical documentation from the current implementation
  and project conventions.
- `/docs-review` — Audit project documentation against the implementation;
  report findings only.

## Requirements

- [Cursor](https://cursor.com)

## Installation

Copy the command files into Cursor's global commands directory:

```bash
mkdir -p ~/.cursor/commands
cp /path/to/cursor-utils/commands/*.md ~/.cursor/commands/
```

The repository does not need to live under `~/.cursor/`.

| Role | Path |
| --- | --- |
| Source | `cursor-utils/commands/*.md` |
| Installed | `~/.cursor/commands/*.md` |

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

Copy the updated `commands/*.md` files into `~/.cursor/commands/` again,
overwriting the previous copies if desired.

## Uninstalling

Delete the specific Markdown files you copied into `~/.cursor/commands/`.
Do not remove unrelated commands in that directory.

## Repository layout

```text
cursor-utils/
├── commands/     # Slash command Markdown files
├── README.md
├── CHANGELOG.md
└── LICENSE
```

## Contributing

Commands live in `commands/` as Markdown files. The basename (without `.md`) is
the slash command name.

When adding or changing a command:

1. Keep the prompt factual and specific to the intended workflow.
2. Do not claim behavior the command does not implement.
3. Update `CHANGELOG.md` under `[Unreleased]`.

## Versioning

This project follows [Semantic Versioning](https://semver.org/).

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## License

[MIT](LICENSE)
