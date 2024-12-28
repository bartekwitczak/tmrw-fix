# TMRW-FIX

TMRW-FIX is a tool designed to help developers manage and track upcoming code changes. When you know you'll need to make changes to certain parts of code tomorrow but want to commit your current work today, you can mark these spots with `tmrw-fix` comments.

Concept is inspired by Hemingway's Bridge technique. It helps you keep going on the next day. Read more about it [here](https://www.bartekwitczak.com/posts/hemingways-bridge).

## How it works

1. Add `tmrw-fix` comments in your code where future changes are needed - it will help you keep going on the next day
2. The tool scans your codebase and shows all locations marked with `tmrw-fix`
3. When integrated with pre-commit hooks, it prevents pushing code containing `tmrw-fix` comments
4. This ensures you don't accidentally push incomplete work or forget about planned changes

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