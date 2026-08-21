# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] — 2026-08-21

### Added

- Initial release of reusable Cursor slash commands:
  - `/commit` — propose atomic Conventional Commits without executing them
  - `/changelog` — update the `[Unreleased]` section of `CHANGELOG.md`
  - `/changelog-draft` — propose `[Unreleased]` changelog entries without
    modifying files
  - `/changelog-review` — audit `CHANGELOG.md` for Keep a Changelog and SemVer
    issues
  - `/release` — convert `[Unreleased]` into a dated SemVer release section
  - `/cleanup` — remove or rewrite low-value comments without changing behavior
  - `/clean-arch` — review architecture boundaries and apply minimal refactors
    when appropriate
  - `/readme` — create or update `README.md` from the current project
  - `/readme-review` — audit `README.md` against the current implementation
  - `/docs` — maintain technical documentation based on the current project
  - `/docs-review` — audit project documentation against the implementation

[Unreleased]: https://github.com/imflawlezz/cursor-utils/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/imflawlezz/cursor-utils/releases/tag/v0.1.0
