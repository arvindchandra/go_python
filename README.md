# Pre-checkin Validator

A production-parity PR gate for CI/CD that runs critical tests in a containerized environment before merge.

## Quick start

```bash
mkdir -p .artifacts
go run ./cmd/validator run --config validator.yaml --changed-files sample_project/app.py,sample_project/tests/test_app.py --output .artifacts
```
