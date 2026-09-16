# pr-draft

## Scope

Do not modify any files and do not run `git push`, `gh pr create`, or any
other write command.

Propose shell commands only. Do not write changelog, README, or `docs/` as part
of this command (`/changelog`, `/readme`, `/docs`).

1. Determine the pull request base branch.
   - Default to `main`.
   - If `main` does not exist locally or on the remote, use `master` when that
     is the repository's default branch.
   - If the base cannot be determined confidently, ask the user before proposing
     commands.

2. Review everything on the current branch that would be included in a PR to
   the base branch.
   - Inspect `git status`.
   - Inspect commits in `<base>..HEAD` with `git log`.
   - Inspect the cumulative diff with `git diff <base>...HEAD`.
   - Note unpushed commits and uncommitted changes separately.

3. If there are uncommitted changes:
   - Do not propose a PR yet.
   - Tell the user to commit or discard them first.
   - Optionally recommend `/commit` for atomic commit proposals.

4. If the current branch is the base branch:
   - Do not propose a PR.
   - Tell the user to create or switch to a feature branch first.

5. If there are no commits on the branch since `<base>`:
   - Do not propose a PR.
   - Report that there is nothing to merge.

6. Draft the pull request from the actual branch content.
   - Title: concise, imperative, Conventional Commits style when it fits
     (`type: description` or `type(scope): description`).
   - Body: use this structure:

     ```markdown
     ## Summary
     - <1-3 bullets of user-facing or maintainer-facing changes>

     ## Test plan
     - [ ] <concrete verification step>
     ```

   - Derive summary bullets from commits and the diff, not from assumptions.
   - Test plan items must be actionable and relevant to the changes.
   - Do not invent tests, CI results, or manual steps that cannot be inferred
     from the branch.

7. Before proposing commands, verify push state.
   - If the branch has no upstream or is ahead of its remote tracking branch,
     include `git push -u origin HEAD` (or an equivalent explicit push) before
     `gh pr create`.
   - If the branch is already up to date with its remote, omit the push.

8. Output the proposed commands in this order when applicable:
   - optional push command
   - `gh pr create` with `--base`, `--title`, and `--body`
   - Use a HEREDOC for the PR body so Markdown formatting is preserved.

   Example shape:

   ```bash
   git push -u origin HEAD

   gh pr create --base main --title "feat(installer): polish TUI usability" --body "$(cat <<'EOF'
   ## Summary
   - ...

   ## Test plan
   - [ ] ...

   EOF
   )"
   ```

9. Do not use:
   - `git push --force` unless the user explicitly requested it
   - broad or destructive git commands
   - placeholder titles or bodies when the branch content is known

10. After the commands, briefly summarize:
    - current branch and base branch
    - commits included (`<base>..HEAD`)
    - whether a push is required
    - anything that blocked proposing a PR

11. NEVER execute `git push`, `gh pr create`, or related write commands
    yourself.
    Only print the commands for the user to review and run manually.
