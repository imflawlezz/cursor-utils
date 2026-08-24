# clean-arch

## Scope

May modify only source, tests, and module structure required for the smallest
coherent architectural fix that preserves behavior.

Must not modify:
- `README.md` (`/readme`)
- `CHANGELOG.md` (`/changelog`)
- files under `docs/` or ADRs (`/docs`)
- comments unrelated to the architectural change (`/cleanup`)

If an intentional exception should be recorded in documentation, report it and
recommend `/docs`. Do not write the ADR or README yourself.

1. Review the project's current architecture before making any changes.

2. Do not assume that the project follows a textbook Clean Architecture
   implementation.

3. Identify the project's actual architectural boundaries, responsibilities,
   and dependency direction from the source code, project structure, build
   configuration, and existing documentation.

4. If the project explicitly follows Clean Architecture, evaluate it against
   the principles actually used by the project.

5. If the project does not follow Clean Architecture, determine whether
   introducing or enforcing Clean Architecture would provide a meaningful
   benefit before proposing structural changes.

6. Check for architectural violations such as:
   - Domain depending on Infrastructure or UI frameworks
   - business rules implemented in Presentation
   - Infrastructure concerns leaking into Domain
   - Application logic coupled directly to concrete infrastructure
   - incorrect dependency direction
   - unnecessary framework dependencies in inner layers
   - duplicated business logic across layers
   - domain models polluted with persistence or presentation concerns
   - use cases bypassed by Presentation or Infrastructure
   - abstractions introduced without a meaningful boundary
   - circular dependencies between architectural layers

7. Distinguish real architectural violations from intentional and justified
   design decisions.

8. Do not introduce abstractions, protocols, interfaces, wrappers, or layers
   solely to satisfy a theoretical Clean Architecture rule.

9. Prefer the simplest architecture that preserves:
   - separation of responsibilities
   - dependency direction
   - testability
   - replaceability of infrastructure
   - maintainability

10. When a violation is found, determine whether it should be:
    - fixed immediately
    - refactored as part of a larger change
    - documented as an intentional exception
    - left unchanged because the proposed abstraction would add unnecessary
      complexity

11. Before modifying architecture, inspect all relevant usages and dependencies
    to avoid creating partial migrations or inconsistent layer boundaries.

12. Do not perform broad architectural rewrites automatically.

13. If changes are appropriate, make the smallest coherent refactoring that
    resolves the identified issue while preserving behavior.

14. After the review, report:
    - the current architectural structure
    - identified violations
    - intentional exceptions
    - recommended improvements
    - changes actually made, if any
