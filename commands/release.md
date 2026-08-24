# release

## Scope

May update only:
- `CHANGELOG.md`: replace `[Unreleased]` with a dated version section and
  update version reference links at the bottom of that file

Must not:
- rewrite historical changelog entries
- modify `README.md`, `docs/`, `LICENSE`, source, or configuration
- create a Git tag, GitHub release, commit, or push
- bump version numbers in package manifests or other files

If other files need a version bump, report that as a follow-up. Do not edit
them.

1. Read the current `CHANGELOG.md` and inspect the `[Unreleased]` section.

2. Determine the appropriate Semantic Versioning bump based on the documented
   changes:
   - MAJOR: breaking changes that require users to change their usage,
     configuration, API integration, or workflow.
   - MINOR: backward-compatible new functionality.
   - PATCH: backward-compatible bug fixes and small corrections.

3. Consider the project's existing version history and conventions when
   determining the new version.

4. If the appropriate version cannot be determined confidently:
   - Do not modify the changelog.
   - Explain the ambiguity.
   - Ask the user to choose the intended version.

5. If `[Unreleased]` contains no changes:
   - Do not create an empty release.
   - Report that there are no changes to release.

6. Replace `[Unreleased]` with a new version section using:
   `## [X.Y.Z] — YYYY-MM-DD`

7. Use the current date as the release date unless the user explicitly
   provides another date.

8. Preserve all existing changelog entries and formatting.
   - Do not rewrite historical releases.
   - Do not re-categorize historical entries.
   - Do not remove information from previous releases.

9. Update the reference links at the bottom of `CHANGELOG.md`:
   - `[Unreleased]` must compare the newly released version against `HEAD`.
   - Add a reference link for the new version.
   - Preserve all existing version links.
   - Follow the URL structure already used by the project.

10. Do not create a Git tag.
11. Do not create a GitHub release.
12. Do not create a commit.
13. Do not push anything.

14. After making the change, report:
   - the selected version
   - the previous version
   - the release date
   - the reason for the SemVer bump
   - the changes included in the release
