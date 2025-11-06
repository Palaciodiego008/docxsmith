# Contributing to DocxSmith

First off, thank you for considering contributing to DocxSmith! It's people like you that make DocxSmith such a great tool.

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct of being respectful and inclusive.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include code samples and error messages**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and explain the behavior you expected to see instead**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go coding style
* Include thoughtful comments in your code
* Write tests for new features
* End all files with a newline
* Ensure all tests pass before submitting

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/docxsmith.git`
3. Create a branch: `git checkout -b feature/my-new-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit your changes: `git commit -am 'Add some feature'`
7. Push to the branch: `git push origin feature/my-new-feature`
8. Submit a pull request

## Development Guidelines

### Code Style

* Follow standard Go conventions
* Use `gofmt` to format your code
* Keep functions small and focused
* Write clear comments for exported functions
* Use meaningful variable names

### Testing

* Write tests for all new features
* Maintain or improve code coverage
* Use table-driven tests where appropriate
* Test edge cases and error conditions

### Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

## Project Structure

```
docxsmith/
├── cmd/docxsmith/      # CLI application
├── pkg/docx/           # Core library
├── examples/           # Usage examples
├── testdata/          # Test fixtures
└── internal/          # Internal packages
```

## Building

```bash
# Build the CLI
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run all pre-commit checks
make pre-commit
```

## Adding New Features

When adding a new feature:

1. **Design First**: Open an issue to discuss the feature before implementing
2. **Write Tests**: Add tests before writing the implementation
3. **Implement**: Write clean, documented code
4. **Test**: Ensure all tests pass
5. **Document**: Update README.md and add examples if needed
6. **Submit**: Create a pull request with a clear description

## Areas for Contribution

We're especially interested in contributions in these areas:

* **Image Support**: Adding and manipulating images in documents
* **Headers/Footers**: Support for document headers and footers
* **Styles**: Advanced style management
* **Charts**: Support for charts and graphs
* **Comments**: Document comments and annotations
* **Performance**: Optimizations for large documents
* **Documentation**: More examples and tutorials
* **Tests**: Increasing test coverage

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to DocxSmith!
