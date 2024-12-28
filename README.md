## Usage with pre-commit

Add this to your `.pre-commit-config.yaml`:

```yaml
repos:
- repo: https://github.com/bartekwitczak/tmrw-fix
  rev: v1.0.0
  hooks:
    - id: tmrw-fix
```

## Installation

```bash
go install github.com/bartekwitczak/tmrw-fix@latest
```