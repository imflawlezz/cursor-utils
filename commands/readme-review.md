# readme-review

## Scope

Review only `README.md`. Do not modify any files.

Do not audit as primary targets:
- `CHANGELOG.md` (`/changelog-review`)
- files under `docs/` (`/docs-review`)
- source comments (`/cleanup`)
- code architecture (`/clean-arch`)

If those files have problems, report a brief handoff to the owning command.

1. Review `README.md` against the current project implementation.

2. Verify, where possible:
   - project description
   - features
   - installation instructions
   - requirements
   - commands
   - configuration
   - usage examples
   - supported platforms
   - dependencies
   - build and development instructions
   - testing instructions
   - architecture descriptions
   - license information

3. Identify:
   - outdated information
   - missing important information
   - incorrect commands
   - broken examples
   - undocumented current behavior
   - references to removed features
   - claims that cannot be verified from the repository
   - unnecessary or overly verbose sections
   - AI-style or marketing-oriented language

4. Do not modify any files.

5. Report actionable findings grouped by:
   - Error
   - Warning
   - Suggestion

6. For every finding, provide a concrete correction.

7. If the README accurately represents the current project, report that no
   issues were found.
