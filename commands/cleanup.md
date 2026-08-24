# cleanup

## Scope

May modify only comments and documentation comments in source files (and
tests) that are in the requested set or in the current project changes.

Must not:
- change program behavior, APIs, architecture, naming, or unrelated formatting
- rewrite `README.md`, `CHANGELOG.md`, `docs/`, or other project Markdown
- create or update ADRs
- execute git commits

If Markdown documentation needs work, recommend `/readme` or `/docs`.

1. Review the requested files or the current project changes.

2. Clean up comments and documentation comments without changing program
   behavior.

3. Remove comments that are:
   - obvious from the code itself
   - redundant or repetitive
   - outdated
   - inaccurate
   - temporary and no longer relevant
   - narrating what the code is doing line by line
   - generated as AI-style narration
   - unnecessarily verbose
   - restating a function, variable, type, or method name without adding
     useful information

4. Pay particular attention to AI-style comments such as:
   - "This function handles..."
   - "We need to..."
   - "Here we..."
   - "Let's..."
   - "Now we..."
   - "This ensures that..."
   - comments describing obvious implementation steps
   - comments written as if an assistant is explaining the code to another
     assistant

5. Preserve comments that provide information that cannot be reasonably
   inferred from the code, especially:
   - non-obvious business rules
   - invariants
   - important constraints
   - compatibility requirements
   - performance considerations
   - security considerations
   - workarounds for external bugs or platform behavior
   - architectural decisions
   - reasons behind non-obvious implementation choices
   - public API documentation

6. Rewrite retained comments when necessary.

7. Retained comments should:
   - be concise
   - use professional technical language
   - explain why rather than merely what
   - describe constraints or intent
   - read like documentation written by an experienced developer
   - avoid conversational or AI-generated phrasing

8. Follow the project's existing documentation and comment conventions.
   Preserve required documentation comments for public APIs, generated
   documentation, or language-specific tooling.

9. Do not add comments merely to replace comments that were removed.
   Prefer self-documenting code when possible.

10. Do not change program behavior, APIs, architecture, formatting unrelated
    to comments, or naming unless explicitly requested.

11. After the cleanup, briefly summarize:
    - comments removed
    - comments rewritten
    - important comments intentionally preserved
