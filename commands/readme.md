# readme

## Scope

May create or update only:
- the repository's primary `README.md`

Must not modify:
- `CHANGELOG.md` (`/changelog`, `/release`)
- files under `docs/` (`/docs`)
- `LICENSE`
- source code, tests, configuration, or build files
- Cursor command, rule, skill, agent, or hook files

You may read those files as context. If detailed documentation belongs in
`docs/`, put a short link in the README and recommend `/docs`; do not create
or edit `docs/` files.

1. Review the entire project before writing or updating `README.md`.

2. Inspect relevant source code, project configuration, package manifests,
   build configuration, scripts, existing documentation, examples, and
   project metadata.

3. Treat the actual project implementation as the source of truth.

4. Update or create `README.md` so that it accurately describes the current
   state of the project.

5. Do not invent:
   - features
   - commands
   - configuration options
   - APIs
   - dependencies
   - supported platforms
   - installation steps
   - screenshots
   - benchmarks
   - compatibility claims
   - project goals

6. Remove or update information that is no longer true.

7. Write for a technically competent developer.

8. Use a professional, concise engineering style:
   - factual
   - direct
   - specific
   - restrained
   - technically accurate

9. Avoid AI-generated writing patterns:
   - "Welcome to..."
   - "In this guide..."
   - "Whether you're..."
   - "This powerful..."
   - "seamlessly"
   - "robust"
   - "elegant"
   - "effortlessly"
   - unnecessary adjectives
   - marketing language
   - repetitive summaries
   - generic explanations of obvious concepts

10. Prefer concrete information over promotional descriptions.

11. Adapt the README structure to the project rather than using a fixed
    template.

12. Include sections that are actually useful for the project, such as:
    - Overview
    - Features
    - Requirements
    - Installation
    - Usage
    - Configuration
    - Development
    - Architecture
    - Testing
    - Building
    - Deployment
    - License

13. Do not add empty, generic, or unnecessary sections.

14. Preserve useful existing README content when it remains accurate.

15. Keep examples synchronized with the actual project.

16. Use commands and code snippets that can be verified against the project.

17. If information required for a README section cannot be established from
    the repository, omit the section rather than inventing content.

18. Update only `README.md` directly.

19. After updating the README, briefly summarize the sections added, changed,
    removed, any information that could not be verified, and any follow-up
    recommended for `/docs` or `/changelog`.
