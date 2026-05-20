# Contributing to DifLog

Thank you for your interest in contributing~ To maintain a high standard of code quality and ensure a smooth development workflow, please follow guidelines mentioned in this document.

---

## Quick Links

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Architecture Overview](ARCHITECTURE.md)

---

## Before you start

1. **Validate Requirement**: Check if someone is already working on the similar topic
2. **Open issue**: Discuss with community before initiating development. After discussion, Open an issue plan development.
3. **Branching Strategy**: We follow GitFlow strategy for development. Use the issue id for your branches.

---

## Development Setup

### Prerequisites

- **Go 1.24.3** - `go version` to validate your version
- **Git** - For version control
- **Terminal Emulators** (optional) - terminal emulators like kitty to test the TUI in better conditions.

### Initiate the development

```bash
git clone https://github.com/xedeveloper/DifLog.git
cd DifLog
go mod download
go build -o diflog .
./diflog 
```

### Create the branch

The project follows [GitFlow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) strategy. Each branch name should be connected to the respective open issue.
for e.g. If there's development on Issue 1 which is for a feature development, then the branch name should be `feature/ISSUE-1`
For hotfix:
`hotfix/ISSUE-1`

See the [GitFlow documentation](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) for reference.

### Run Test

For each development, it is required to write test cases for the logic layer updated.

```bash
# Run all test cases
go test ./...
```

---

## Project Structure

The project follows CLEAN architecture for the development. Please refer to [Architecture Documentation](/ARCHITECTURE.md)
---

## Coding Standards

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Run `go fmt` before committing

### Comments

**None by default.** Comments are allowed only when logic is complicated to describe with the function name.
Function, Variable names should be self explanatory to avoid the inline comments.

### Error Handling

- Provide explanatory context to the `fmt.Errorf` block
- Errors should be handled & should not directly affect the user journey.

---

## Pull Request

### Pull Request Prerequisites

- [ ] All test passed locally
- [ ] Code is formatted
- [ ] Commit message follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

### PR title format

The PR title should consist of your issue number followed with title of issue
For e.g.
`ISSUE-1: Support for Copilot CLI`

### Pull request description

Pull request template is created & should be followed for the following development.

---

## Code of Conduct

We follow the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

**Thanks for contributing!**
Your every line of code helps us make the DifLog better for developers tackling with AI memory issues.
