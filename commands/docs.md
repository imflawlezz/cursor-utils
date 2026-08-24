# docs

Maintain technical documentation in `docs/` based on the current
implementation, existing documentation, and established project conventions.

The goal is accurate, useful, maintainable documentation written in the style
of an experienced software engineer. Do not generate documentation merely to
increase the amount of documentation in the repository.

## Scope

May create, update, or remove only:
- files under `docs/`
- Architecture Decision Records in the project's established ADR location
  (typically `docs/decisions/` or equivalent)
- other technical documentation trees the project already uses for the same
  purpose as `docs/` (guides, reference, architecture)

Must not modify:
- `README.md` or other README files (`/readme`)
- `CHANGELOG.md` (`/changelog`, `/release`)
- `LICENSE`
- source code, tests, configuration, or build files
- Cursor command, rule, skill, agent, or hook files

You may read those files as context. If they need changes, report a handoff
to the owning command instead of editing them.

1. Understand the project before changing documentation.

   Inspect, where relevant:
   - source code
   - project structure
   - package manifests and dependencies
   - build configuration
   - scripts and development tooling
   - tests
   - configuration files
   - existing `docs/` content
   - architecture documentation
   - ADRs or decision records
   - API specifications
   - generated documentation configuration
   - project-specific documentation conventions
   - `README.md` and `CHANGELOG.md` as read-only context only

   Treat the implementation and configuration as the primary source of truth.
   Treat existing documentation as authoritative only when it is consistent
   with the implementation.

2. Determine what documentation actually needs to change.

   Look for:
   - newly introduced user or developer workflows
   - changed behavior
   - changed APIs
   - changed configuration
   - changed architecture
   - new dependencies or requirements
   - new development procedures
   - removed or deprecated functionality
   - outdated examples
   - outdated commands
   - missing explanations for important non-obvious behavior

   Do not create or modify documentation when there is no meaningful
   documentation change to make.

3. Determine the appropriate documentation type.

   Use the following documentation model:

   - Tutorial:
     A learning-oriented, step-by-step path for someone unfamiliar with the
     project or technology.

   - How-to:
     A focused procedure for accomplishing a specific task.

   - Reference:
     Precise factual information such as APIs, configuration options,
     commands, environment variables, schemas, file formats, or supported
     values.

   - Explanation:
     Conceptual, architectural, or design-oriented documentation explaining
     how something works, why it works that way, or what constraints shaped
     the design.

   Use the category that best matches the reader's goal.

   Do not force every project to contain all four categories.

4. Inspect existing documentation before creating files.

   Search for documents covering the same or closely related topic.

   If an appropriate document already exists:
   - update it instead of creating a duplicate
   - preserve useful existing content
   - reorganize it only when necessary for correctness or usability

   If several documents overlap substantially:
   - consolidate them when appropriate
   - avoid multiple competing sources of truth

5. Choose the simplest documentation structure that fits the project.

   Do not create directories or documents merely because they are common in
   documentation templates.

   Possible structures include:

   `docs/guides/`
   `docs/reference/`
   `docs/architecture/`
   `docs/decisions/`

   Use only the directories that provide real value.

   A small project may need no `docs/` tree at all if `README.md` and
   `CHANGELOG.md` are sufficient. In that case, make no documentation files
   and do not edit README or CHANGELOG.

   A larger project may benefit from a more structured documentation tree
   under `docs/`.

6. Keep README and technical documentation separate.

   `README.md` is owned by `/readme`. Do not edit it.

   Detailed architecture, API reference, design explanations, development
   procedures, and extensive technical reference belong in `docs/` when they
   would make the README unnecessarily long.

   Do not duplicate large sections of documentation between README and `docs/`.
   If the README needs a short link to new `docs/` pages, report that as a
   follow-up for `/readme`. Do not add the link yourself.

7. Keep architecture documentation separate from decision records.

   Architecture documentation describes the current system:
   - layers
   - modules
   - responsibilities
   - dependency direction
   - data flow
   - important integration boundaries

   Architecture Decision Records describe significant decisions and their
   reasoning:
   - context
   - decision
   - consequences

   Do not use an ADR as a substitute for current architecture documentation.

8. Create an ADR only for a significant architectural or technical decision.

   Do not create ADRs for:
   - trivial implementation choices
   - routine refactoring
   - obvious code organization
   - temporary implementation details
   - decisions that are unlikely to matter to future maintainers

   If the project already has an ADR convention, follow it.

   Otherwise use:

   `# <Decision>`

   `Status: <Proposed|Accepted|Deprecated|Superseded>`
   `Date: YYYY-MM-DD`

   `## Context`

   `## Decision`

   `## Consequences`

   Use sequential numeric filenames when the project uses numbered ADRs, for
   example:

   `docs/decisions/001-use-avfoundation.md`

9. Write documentation from verified facts.

   Never invent:
   - APIs
   - commands
   - configuration options
   - environment variables
   - dependencies
   - supported platforms
   - compatibility requirements
   - performance characteristics
   - benchmarks
   - architecture boundaries
   - workflows
   - examples
   - limitations

   If something cannot be established from the repository, omit it or clearly
   identify it as requiring confirmation.

10. Prefer concrete examples when they improve understanding.

    Examples must:
    - correspond to the actual implementation
    - use valid commands and syntax
    - use realistic values
    - remain concise
    - not duplicate information unnecessarily

11. Write in professional technical prose.

    Documentation should be:
    - factual
    - concise
    - precise
    - direct
    - consistent
    - easy to scan

    Prefer short paragraphs, meaningful headings, lists, tables, and code
    blocks when they improve readability.

12. Avoid AI-generated writing patterns.

    Do not use:
    - marketing language
    - exaggerated claims
    - unnecessary adjectives
    - generic introductions
    - repetitive summaries
    - conversational narration
    - "Welcome to..."
    - "In this guide, we will..."
    - "Let's take a look at..."
    - "This powerful..."
    - "seamlessly"
    - "robust"
    - "elegant"
    - "effortlessly"
    - "it's important to note that..." unless genuinely necessary

    Do not write documentation as if an AI assistant is explaining the project
    to the reader.

13. Explain why when the reason is important.

    Prefer documenting:
    - non-obvious constraints
    - architectural boundaries
    - external platform limitations
    - compatibility requirements
    - security considerations
    - performance trade-offs
    - persistence decisions
    - important workarounds
    - rejected alternatives when documented through an ADR

    Do not document obvious implementation details that can be understood
    directly from the source code.

14. Keep documentation focused on stable concepts.

    Document public interfaces, workflows, architecture, configuration, and
    important behavior.

    Avoid documenting implementation details that are likely to change
    frequently unless those details are themselves relevant to developers.

15. Keep in-scope documentation internally consistent.

    When updating a concept, check related files under `docs/` for:
    - outdated terminology
    - contradictory behavior
    - obsolete commands
    - old file paths
    - removed features
    - outdated architecture descriptions
    - stale examples

    Update related `docs/` files when necessary rather than leaving conflicting
    sources of truth inside `docs/`.

    If `README.md` or `CHANGELOG.md` conflict with the implementation, note the
    conflict in the summary and recommend `/readme` or `/changelog`. Do not
    edit those files.

16. Do not modify out-of-scope files.

    Limit changes to `docs/` (and equivalent in-scope technical docs) that are:
    - directly relevant to the current task
    - demonstrably outdated
    - necessary to keep in-scope cross-references consistent

17. Preserve existing project conventions.

    Follow the project's:
    - naming conventions
    - Markdown style
    - heading hierarchy
    - terminology
    - link style
    - code formatting
    - documentation structure
    - ADR format
    - language conventions

    Do not introduce a new documentation convention when an established one
    already exists.

18. Validate documentation after writing it.

    Check:
    - referenced files exist
    - referenced commands exist
    - referenced scripts exist
    - referenced configuration options exist
    - internal links point to valid files or sections
    - examples match the current implementation
    - terminology is consistent
    - headings are logically structured
    - no duplicate source of truth was introduced

19. Keep documentation maintainable.

    Prefer a small number of focused documents over large catch-all files.

    Split a document when it covers several unrelated concerns.

    Do not split documents merely to create more files.

20. If the documentation is already accurate and complete, make no changes.

21. Modify only in-scope documentation files directly.

22. After completing the task, provide a concise summary containing:
    - documentation files created
    - documentation files updated
    - documentation files removed, if any
    - the documentation type used for each new document
    - important inconsistencies that were corrected
    - anything that could not be verified from the repository
    - recommended follow-up commands for out-of-scope files, if any
