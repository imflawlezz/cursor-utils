# changelog-draft

## Scope

Do not modify any files.

Output only a proposed `[Unreleased]` section for `CHANGELOG.md`.

Do not draft README, `docs/`, release notes outside the changelog, git tags,
or version numbers.

1. Determine the latest release tag in Git.
   - Prefer the latest reachable semantic-version tag matching `vX.Y.Z`.
   - Use that tag as the baseline.
   - If no semantic-version tag exists, use the initial commit.

2. Review the changes from the baseline tag to `HEAD`.
   - Inspect both commit history and the cumulative diff.
   - Consider the existing `CHANGELOG.md` and its `[Unreleased]` section.

3. Produce a proposed `[Unreleased]` section following Keep a Changelog
   1.1.0 and the project's existing changelog style.

4. Use only the sections that are actually needed:
   - Added
   - Changed
   - Deprecated
   - Removed
   - Fixed
   - Security

5. Write user-facing changelog entries.
   - Do not copy commit messages verbatim.
   - Describe observable behavior and meaningful project changes.
   - Group related implementation changes.
   - Avoid trivial internal changes.
   - Preserve useful details such as shortcuts, compatibility requirements,
     migration requirements, supported formats, and notable UX changes.

6. Compare the proposed entries with the existing `[Unreleased]` section.
   - Do not duplicate existing entries.
   - Clearly distinguish already documented changes from newly discovered ones.

7. Match the existing changelog's terminology, Markdown formatting, emphasis,
   and level of detail.

8. Do not modify any files.

9. Output only the proposed changelog content for `[Unreleased]`.
