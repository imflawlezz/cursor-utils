# commit

1. Review the current uncommitted changes, including both staged and unstaged changes.

2. Split changes into logically atomic commits.
   - Each commit must represent one coherent change: one feature, one fix,
     one refactor, one documentation change, etc.
   - Do not group unrelated changes together.
   - Do not split a single logical change across multiple commits.

3. For each proposed commit, output:
   - `git add <specific files>`
   - `git commit -m "<message>"` or, when a body is useful, the appropriate
     multi-line `git commit` command.

4. Never use:
   - `git add .`
   - `git add -A`
   - broad or unspecified file patterns when specific files can be identified.

5. Commit message rules:
   - Use Conventional Commits format: `type: description` or
     `type(scope): description`
   - Use a scope when it adds useful context, especially in large or diverse
     projects.
   - Omit the scope when it provides no meaningful information.
   - English only.
   - Use imperative mood ("add", not "added" or "adds").
   - Keep the subject concise and descriptive.
   - Lowercase after the colon.
   - Do not add a trailing period to the subject.
   - Prefer a single-line commit message for simple changes.
   - A commit body is allowed when additional context is genuinely useful.
   - If a body is used, keep it focused on why the change was made or important
     implementation context.
   - Do not use a body just to restate the subject.

6. Allowed commit types:
   - feat
   - fix
   - refactor
   - perf
   - test
   - docs
   - build
   - ci
   - chore

7. Before proposing commits, inspect the actual diff and use the project's
   existing structure and conventions to determine the correct scope and type.

8. NEVER execute `git add` or `git commit` yourself.
   Only print the commands for the user to review and run manually.
