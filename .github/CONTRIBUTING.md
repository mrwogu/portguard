# Contributing to PortGuard

First off, thank you for considering contributing to PortGuard! It's people like you that make PortGuard such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include logs and configuration files** (sanitize sensitive data)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and explain which behavior you expected to see instead**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go coding style
* Include thoughtfully-worded, well-structured tests
* Document new code
* End all files with a newline

## Development Process

1. Fork the repo
2. Create a new branch from `main`:
   ```bash
   git checkout -b feature/my-new-feature
   ```
3. Make your changes
4. Run tests:
   ```bash
   go test ./...
   ```
5. Run linter:
   ```bash
   golangci-lint run
   ```
6. Commit your changes:
   ```bash
   git commit -am 'Add some feature'
   ```
7. Push to the branch:
   ```bash
   git push origin feature/my-new-feature
   ```
8. Create a new Pull Request

## Coding Conventions

* Use `gofmt` to format your code
* Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
* Write clear commit messages
* Add comments for exported functions and types
* Keep functions small and focused
* Write tests for new functionality

## Testing

Run all tests:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

