# go-mail-testserver

A lightweight SMTP test server with HTTP API for integration testing. Capture and inspect emails sent during tests without actually sending them to real recipients.

## Features

- **SMTP Server**: Receives emails on port 1025 (configurable)
- **HTTP API**: Query received emails via REST API on port 8025 (configurable)
- **In-Memory Storage**: Fast, ephemeral email storage perfect for testing
- **Simple Integration**: Easy to integrate into your test suites
- **No Dependencies**: Single binary with minimal external requirements

## Quick Start

### Prerequisites

- Go 1.25.5 or later
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/joukojo/go-mail-testserver.git
cd go-mail-testserver

# Build the server
go build -o mail-testserver ./cmd/mail-testserver

# Run the server
./mail-testserver
```

The server will start with:
- SMTP server listening on `:1025`
- HTTP API listening on `:8025`

### Configuration

Configure the server using environment variables:

```bash
# Custom SMTP port
export SMTP_ADDR=":2525"

# Custom HTTP port
export HTTP_ADDR=":9025"

# Run with custom configuration
./mail-testserver
```

## Usage

### Sending Emails

Send emails to the SMTP server using any SMTP client:

```bash
# Using telnet
telnet localhost 1025

# Using a mail client in your application
# Configure your SMTP settings to:
# Host: localhost
# Port: 1025
# No authentication required
```

### Retrieving Emails via HTTP API

The HTTP API provides several endpoints to interact with received emails:

#### Get All Messages

```bash
curl http://localhost:8025/api/v1/messages
```

Response:
```json
[
  {
    "id": 1,
    "from": "sender@example.com",
    "to": ["recipient@example.com"],
    "subject": "Test Email",
    "body": "Email body content",
    "date": "2026-01-04T10:30:00Z"
  }
]
```

#### Get Specific Message

```bash
curl http://localhost:8025/api/v1/messages/1
```

#### Get Raw Message

```bash
curl http://localhost:8025/api/v1/messages/1/raw
```

#### Clear All Messages

```bash
curl -X POST http://localhost:8025/api/v1/messages/clear
```

### Integration Testing Example

```go
// Example test in Go
func TestEmailSending(t *testing.T) {
    // Your application sends email to localhost:1025
    
    // Wait a bit for email to be processed
    time.Sleep(100 * time.Millisecond)
    
    // Query the test server
    resp, err := http.Get("http://localhost:8025/api/v1/messages")
    if err != nil {
        t.Fatal(err)
    }
    
    var messages []Message
    json.NewDecoder(resp.Body).Decode(&messages)
    
    // Assert email was sent correctly
    if len(messages) != 1 {
        t.Errorf("Expected 1 message, got %d", len(messages))
    }
}
```

## Development

### Setting Up Development Environment

1. **Clone the repository**:
   ```bash
   git clone https://github.com/joukojo/go-mail-testserver.git
   cd go-mail-testserver
   ```

2. **Install development tools**:
   ```bash
   make init
   ```
   
   This will install:
   - `goimports` - Import formatter
   - `gofumpt` - Go formatter
   - `govulncheck` - Vulnerability checker
   - `golangci-lint` - Linter
   - Git hooks for pre-commit checks

3. **Run tests**:
   ```bash
   make test
   ```

### Project Structure

```
.
├── cmd/
│   ├── mail-testserver/    # Main server application
│   └── mali-testclient/    # Test client (if needed)
├── internal/
│   ├── commonssmtp/        # SMTP server implementation
│   └── httpapi/            # HTTP API and storage
├── apidocs/
│   └── openapi.yml         # API documentation
├── go.mod
├── Makefile
└── README.md
```

### Development Workflow

1. **Make your changes** to the relevant files

2. **Format your code**:
   ```bash
   make fmt
   ```

3. **Run tests**:
   ```bash
   make test
   ```

4. **Run all checks** (format, lint, test, vulnerability check):
   ```bash
   make ci
   ```

5. **Commit your changes**:
   ```bash
   git add .
   git commit -m "Description of your changes"
   ```
   
   The pre-commit hook will automatically run checks before committing.

### Adding New Features

#### Adding a New HTTP API Endpoint

1. **Open** [internal/httpapi/api.go](internal/httpapi/api.go)

2. **Add your handler function**:
   ```go
   func (s *Server) handleNewFeature(w http.ResponseWriter, r *http.Request) {
       // Your implementation
   }
   ```

3. **Register the route** in the `Start()` method:
   ```go
   mux.HandleFunc("/api/v1/newfeature", s.handleNewFeature)
   ```

4. **Update the OpenAPI spec** in [apidocs/openapi.yml](apidocs/openapi.yml)

5. **Add tests** in `internal/httpapi/api_test.go`

#### Extending Storage

1. **Open** [internal/httpapi/storage.go](internal/httpapi/storage.go)

2. **Add new methods** to the `Storage` struct

3. **Write tests** in [internal/httpapi/storage_test.go](internal/httpapi/storage_test.go)

#### Modifying SMTP Behavior

1. **Open** [internal/commonssmtp/server.go](internal/commonssmtp/server.go)

2. **Modify the SMTP handler** logic

3. **Test thoroughly** with various SMTP clients

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./internal/httpapi

# Run specific test
go test -v -run TestSpecificFunction ./internal/httpapi
```

### Code Quality Checks

```bash
# Format code
make fmt

# Run linter
make lint

# Check for vulnerabilities
make vuln

# Run all CI checks
make ci
```

## API Documentation

Full API documentation is available in [apidocs/openapi.yml](apidocs/openapi.yml). You can view it using any OpenAPI viewer.

### Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/messages` | Get all received messages |
| GET | `/api/v1/messages/{id}` | Get specific message by ID |
| GET | `/api/v1/messages/{id}/raw` | Get raw message content |
| POST | `/api/v1/messages/clear` | Clear all messages |

## Reporting Issues

We appreciate your help in making this project better! If you encounter any issues or have suggestions:

### Before Reporting

1. **Check existing issues** at [GitHub Issues](https://github.com/joukojo/go-mail-testserver/issues) to avoid duplicates
2. **Verify the issue** can be reproduced
3. **Check** you're using the latest version

### How to Report an Issue

1. **Go to** [GitHub Issues](https://github.com/joukojo/go-mail-testserver/issues)

2. **Click** "New Issue"

3. **Provide the following information**:
   - **Clear title**: Brief description of the issue
   - **Description**: Detailed explanation of the problem
   - **Steps to reproduce**: Exact steps to reproduce the issue
   - **Expected behavior**: What you expected to happen
   - **Actual behavior**: What actually happened
   - **Environment**:
     - OS: (e.g., macOS 14.1, Ubuntu 22.04)
     - Go version: `go version`
     - Server version/commit
   - **Logs**: Relevant server logs or error messages
   - **Configuration**: Any custom environment variables

### Issue Template Example

```markdown
**Description**
Brief description of the issue

**Steps to Reproduce**
1. Start the server with `./mail-testserver`
2. Send an email with subject "Test"
3. Call GET /api/v1/messages

**Expected Behavior**
The email should appear in the response

**Actual Behavior**
The response is empty

**Environment**
- OS: macOS 14.1
- Go version: 1.25.5
- Server version: commit abc123

**Logs**
```
[paste relevant logs here]
```
```

### Security Issues

For security-related issues, please **DO NOT** create a public issue. Instead:
- Email the maintainers directly
- Provide full details of the vulnerability
- Allow reasonable time for a fix before public disclosure

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Make sure all tests pass and code is formatted before submitting.

## License

[Add your license information here]

## Support

- **Issues**: [GitHub Issues](https://github.com/joukojo/go-mail-testserver/issues)
- **Discussions**: [GitHub Discussions](https://github.com/joukojo/go-mail-testserver/discussions)

## Acknowledgments

Built with:
- [emersion/go-smtp](https://github.com/emersion/go-smtp) - SMTP server implementation