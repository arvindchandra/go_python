# Pre-checkin Validator

A production-parity PR gate for CI/CD that runs critical tests in a containerized environment before merge. The project uses a Go CLI for orchestration and Python helpers for test selection, flaky-test quarantine, and JUnit reporting. This design follows common Go CLI layout guidance, uses Docker-based test execution, and is wired for GitHub Actions required status checks on pull requests.[web:144][web:145][web:183][web:184]

## Features

- Detects changed files and maps them to services/tests.
- Builds a merge-gating test plan from `validator.yaml`.
- Runs tests through Docker Compose in a production-like container environment.
- Retries suspected flaky tests once in future extensions.
- Honors a quarantine registry so known flakes stay visible without blocking merges.
- Produces JUnit and JSON summaries for GitHub Actions artifacts.

## Project layout

- `cmd/validator`: Go CLI entrypoint.
- `internal/config`: YAML config loading.
- `internal/gitdiff`: Changed-file discovery.
- `internal/selector`: Test planning.
- `internal/runner`: Docker/test execution.
- `internal/report`: JUnit/JSON summary writing.
- `scripts`: Python helpers.
- `.github/workflows`: GitHub Actions PR validation workflow.
- `sample_project`: toy Python app and tests to prove the flow.

## Quick start

```bash
mkdir -p .artifacts

go run ./cmd/validator run   --config validator.yaml   --changed-files sample_project/app.py,sample_project/tests/test_app.py   --output .artifacts
```

## Push to GitHub

```bash
git init
git add .
git commit -m "Initial pre-checkin validator"
git branch -M main
git remote add origin https://github.com/<your-org>/<your-repo>.git
git push -u origin main
```

## Enable merge gating in GitHub

1. Push the repository to GitHub.
2. Let the workflow in `.github/workflows/precheckin.yml` run on `push` to `main` at least once so the check name appears in branch protection. GitHub only shows recent checks that have actually run on the branch when you configure required status checks.[web:192][web:186]
3. In GitHub go to **Settings → Branches**.
4. Add a protection rule for `main`.
5. Enable **Require a pull request before merging** and **Require status checks to pass before merging**.
6. Select the `precheckin-validator` job as the required check. GitHub protected branches can require specific status checks before merges are allowed.[web:183][web:184][web:196]

## GitHub Actions behavior

The workflow runs on both `pull_request` and `push` to `main`, uploads `.artifacts` as a workflow artifact, and uses a unique job name so it can be selected cleanly in branch protection. GitHub recommends unique job names for required checks, and artifact upload is the standard mechanism for preserving test outputs from workflows.[web:183][web:194]
