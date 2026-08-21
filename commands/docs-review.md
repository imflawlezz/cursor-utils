# docs-review

Audit the project's documentation against the current implementation and
report actionable problems. Do not modify any files.

1. Understand the project before reviewing its documentation.

   Inspect, where relevant:
   - source code
   - project structure
   - package manifests and dependencies
   - build configuration
   - scripts
   - tests
   - configuration files
   - README.md
   - `docs/`
   - architecture documentation
   - ADRs
   - API specifications
   - contribution guidelines

   Treat the implementation and configuration as the primary source of truth.

2. Identify all relevant documentation sources.

   Review:
   - `README.md`
   - files under `docs/`
   - architecture documents
   - ADRs
   - API documentation
   - contribution documentation
   - other repository Markdown files that function as technical documentation

   Do not treat `CHANGELOG.md` as general documentation, but check it when
   documentation claims depend on historical release behavior.

3. Check documentation accuracy.

   Look for:
   - incorrect descriptions of current behavior
   - outdated APIs
   - removed features still being documented
   - missing newly introduced behavior
   - incorrect commands
   - obsolete scripts
   - invalid configuration options
   - nonexistent environment variables
   - incorrect file paths
   - outdated dependency information
   - incorrect supported-platform claims
   - incorrect compatibility requirements
   - stale examples
   - incorrect architecture descriptions

4. Check documentation completeness.

   Identify important project concepts that are currently undocumented,
   especially:
   - public APIs
   - developer workflows
   - important configuration
   - architecture boundaries
   - non-obvious constraints
   - integration points
   - persistence behavior
   - build requirements
   - testing procedures
   - deployment procedures
   - significant technical decisions

   Do not report missing documentation for trivial or self-explanatory
   implementation details.

5. Check documentation structure.

   Identify:
   - duplicate sources of truth
   - documents covering overlapping topics
   - large catch-all documents that should be split
   - unnecessary documentation files
   - misplaced documentation
   - architecture information mixed into unrelated guides
   - reference information presented as tutorials
   - procedural guides presented as architecture explanations

   Use the following documentation categories as a conceptual reference:

   - Tutorial: learning-oriented
   - How-to: task-oriented
   - Reference: factual and precise
   - Explanation: conceptual or architectural

   Do not require a specific directory structure if the project's current
   structure is reasonable.

6. Check README boundaries.

   Ensure `README.md` remains focused on:
   - project overview
   - purpose
   - primary features
   - requirements
   - installation
   - initial usage
   - links to deeper documentation

   Flag large amounts of detailed architecture, API reference, or specialized
   development documentation that would be better placed in `docs/`.

   Do not flag concise technical information that is genuinely useful during
   initial setup or usage.

7. Check architecture documentation.

   Verify that documented:
   - layers
   - modules
   - responsibilities
   - dependency direction
   - data flow
   - integration boundaries

   match the current implementation.

   Distinguish between:
   - actual architectural inconsistencies
   - intentional exceptions
   - missing documentation

   Do not infer an architectural violation merely because the project does not
   follow a textbook architecture.

8. Check Architecture Decision Records.

   Verify that ADRs:
   - describe real decisions
   - contain sufficient context
   - clearly state the decision
   - describe meaningful consequences
   - reflect the historical decision accurately

   Flag ADRs that contradict the current architecture without explaining that
   they are historical.

   Do not require an ADR for routine or trivial decisions.

9. Check examples and commands.

   Where possible, verify every documented:
   - shell command
   - script
   - API example
   - configuration example
   - file path
   - code snippet

   Flag examples that would not work with the current project.

10. Check internal links and references.

    Verify:
    - relative links
    - referenced documentation files
    - anchors where practical
    - referenced source files
    - referenced scripts
    - cross-references between documentation

    Flag broken or stale references.

11. Check writing quality.

    Identify documentation that is:
    - vague
    - unnecessarily verbose
    - repetitive
    - inconsistent
    - overly conversational
    - promotional
    - filled with generic filler
    - written in an AI-narration style

    Flag phrases such as:
    - "Welcome to..."
    - "In this guide, we will..."
    - "Let's take a look at..."
    - "This powerful..."
    - "seamlessly"
    - "robust"
    - "elegant"
    - "effortlessly"

    Do not flag normal technical prose merely because it is informal.

12. Check maintainability.

    Identify documentation that is tightly coupled to unstable implementation
    details without a clear reason.

    Prefer documentation of:
    - stable interfaces
    - stable workflows
    - architectural boundaries
    - important constraints
    - externally observable behavior

13. Check for conflicting sources of truth.

    When two documents describe the same concept differently:
    - identify the conflict
    - determine which source is consistent with the implementation
    - recommend which document should become authoritative

14. Do not modify any files.

15. Report only actionable findings.

16. Group findings by severity:

    Error
    A documented fact is incorrect, a command/example is broken, or the
    documentation materially conflicts with the implementation.

    Warning
    Important documentation is missing, outdated, duplicated, poorly
    structured, or likely to cause confusion.

    Suggestion
    The documentation is technically correct but could be clearer, more
    maintainable, or better organized.

17. For every finding, include:
    - file
    - relevant section or heading
    - problem
    - evidence from the current project
    - recommended correction

18. Do not report speculative problems.

    Every finding must be supported by the repository or by an observable
    inconsistency between documentation sources.

19. At the end, provide a short overall assessment:
    - Documentation status: Good / Needs attention / Significantly outdated
    - Most important issue
    - Recommended next action

20. If no actionable issues are found, explicitly report that the documentation
    is consistent with the current implementation.
