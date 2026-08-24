# changelog-review

## Scope

Review only `CHANGELOG.md`. Do not modify any files.

Do not audit `README.md` (`/readme-review`) or `docs/` (`/docs-review`) except
as brief handoffs when a changelog claim depends on them.

1. Review the current `CHANGELOG.md` against:
   - Keep a Changelog 1.1.0
   - Semantic Versioning
   - the project's existing changelog conventions

2. Compare the changelog with the actual Git history and project changes where
   useful, especially changes since the latest release tag.

3. Check for:
   - missing meaningful user-facing changes
   - duplicated entries
   - incorrect categories
   - overly fragmented entries
   - unrelated changes grouped together
   - vague or unhelpful descriptions
   - implementation details presented as user-facing changes
   - inconsistent terminology
   - inconsistent Markdown formatting
   - incorrect version ordering
   - incorrect release dates
   - missing or broken version reference links
   - incorrect `[Unreleased]` comparison links
   - possible Semantic Versioning violations

4. Pay particular attention to whether changes are categorized correctly:
   - Added for new functionality
   - Changed for modifications to existing functionality
   - Deprecated for functionality that still exists but should no longer be used
   - Removed for removed functionality
   - Fixed for bug fixes
   - Security for security-related changes

5. Do not modify any files.

6. Report only actionable findings.

7. Group findings by severity:
   - Error
   - Warning
   - Suggestion

8. For every finding, explain:
   - what is wrong
   - why it matters
   - what concrete correction is recommended

9. If the changelog is already correct, explicitly state that no issues were
   found.
