# changelog

1. Determine the latest release tag in Git.
   - Prefer the latest reachable semantic-version tag matching `vX.Y.Z`.
   - Treat this tag as the baseline for the changelog update.
   - Do not use an arbitrary number of recent commits as the baseline.

2. Review all changes between the latest release tag and `HEAD`.
   - Inspect commits in `<latest-tag>..HEAD`.
   - Inspect the cumulative diff between `<latest-tag>` and `HEAD`.
   - Use both commit history and the actual diff to understand what changed.
   - Focus on meaningful user-facing changes rather than merely counting
     commits.

3. If no semantic-version release tag exists:
   - Review the project's history from the initial commit.
   - Use the existing `CHANGELOG.md` as additional context.
   - Do not invent a previous release tag.

4. Read the existing `CHANGELOG.md` before making changes.
   - Preserve the project's existing style, terminology, level of detail,
     formatting, and structure.
   - Follow Keep a Changelog 1.1.0 conventions.
   - Follow the project's existing Semantic Versioning conventions.

5. Update only the `[Unreleased]` section.
   - Never modify released versions or their historical entries.
   - If `[Unreleased]` does not exist, create it immediately before the latest
     released version.

6. Treat the existing `[Unreleased]` section as already documented work.
   - Preserve existing entries.
   - Compare existing entries with the changes since the latest release tag.
   - Add only missing or materially changed entries.
   - Do not create duplicate entries.
   - Do not rewrite existing entries without a good reason.

7. Categorize changes using only the sections that are actually needed:
   - Added
   - Changed
   - Deprecated
   - Removed
   - Fixed
   - Security

8. Write changelog entries for users and maintainers, not for Git history.
   - Describe what changed and its observable effect.
   - Do not simply copy commit messages.
   - Do not expose implementation details unless they are relevant to users,
     developers, compatibility, migration, performance, or architecture.
   - Group related implementation changes into one meaningful entry.
   - Do not merge unrelated changes into one entry.
   - Do not create entries for trivial internal changes that have no meaningful
     impact.

9. Preserve important details when they are relevant:
   - keyboard shortcuts
   - changed behavior
   - configuration changes
   - compatibility requirements
   - migration requirements
   - supported formats
   - notable performance changes
   - meaningful UX changes
   - user-visible architectural or dependency changes

10. Match the existing changelog's writing style.
    - Keep entries concise but informative.
    - Use Markdown consistently.
    - Preserve existing emphasis conventions such as bold shortcuts,
      commands, filenames, and important terms.
    - Prefer concrete descriptions over vague phrases such as "improve UX"
      or "various fixes".

11. Do not invent information that cannot be established from the repository,
    Git history, or existing changelog.

12. Do not create a version number or release date.
    All new changes belong under `[Unreleased]`.

13. After updating the file, briefly summarize:
    - which categories were added or updated
    - which meaningful changes were documented
    - whether any changes were intentionally omitted and why

14. Make the changelog update directly in `CHANGELOG.md`.
