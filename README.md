# Chat Scraper CLI

A command-line tool for extracting and processing chat conversations.

## Prerequisites

- Go 1.19 or higher
- Git

## Installation and Build

### In WSL (Windows Subsystem for Linux)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/seu-usuario/chat-scraper-cli.git
   cd chat-scraper-cli
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build and installation:**
   ```bash
   go install
   ```

4. **Verify installation:**
   ```bash
   chat-scraper-cli --help
   ```

### On Windows

#### Option 1: Go installed natively on Windows

1. **Clone the repository:**
   ```cmd
   git clone https://github.com/seu-usuario/chat-scraper-cli.git
   cd chat-scraper-cli
   ```

2. **Configure environment variable:**
   
   **Via command line (temporary):**
   ```cmd
   set WINDOWS_DIR=%USERPROFILE%\go\bin
   ```
   
   **Via system environment variables (permanent):**
   - Open "System Environment Variables"
   - Click on "Environment Variables"
   - Under "User variables", click "New"
   - Variable name: `WINDOWS_DIR`
   - Variable value: `%USERPROFILE%\go\bin`
   - Click "OK"

3. **Install dependencies:**
   ```cmd
   go mod download
   ```

4. **Build for Windows:**
   ```cmd
   go build -o chat-scraper-cli.exe
   ```

5. **Move executable to PATH directory:**
   ```cmd
   move chat-scraper-cli.exe %WINDOWS_DIR%
   ```

6. **Verify installation:**
   ```cmd
   chat-scraper-cli.exe --help
   ```

#### Option 2: Go installed only in WSL (Cross-Platform Build)

1. **In WSL, clone the repository:**
   ```bash
   git clone https://github.com/seu-usuario/chat-scraper-cli.git
   cd chat-scraper-cli
   ```

2. **Configure the environment variable in WSL:**
   `windows_path` is the directory where your CLI executable will be created. Ensure the directory exists before setting this environment variable.

   ```bash
   export WINDOWS_DIR=[windows_path]
   ```
3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Cross-platform build for Windows:**
   ```bash
   GOOS=windows GOARCH=amd64 go build -o $WINDOWS_DIR/chat-scraper-cli.exe
   ```

5. **Open the CLI output directory**
   Open PowerShell or Windows Terminal, and navigate to the `windows_path` defined in step 2.

6. **Verify installation on Windows:**
   ```cmd
   chat-scraper-cli.exe --help
   ```

**Note:** The `WINDOWS_DIR` variable in WSL should point to the Go bin directory on Windows, which is usually `/mnt/c/Users/$USER/go/bin` (where `$USER` is your Windows username).

## Usage

```bash
# Basic example
chat-scraper-cli [options]

# Show help
chat-scraper-cli --help
```

## Development

### Project Structure

```
chat-scraper-cli/
├── main.go          # Application entry point
├── go.mod           # Go dependencies
├── go.sum           # Dependency checksums
└── temp_chats/      # Temporary directory for chats
```

### Development Commands

```bash
# Run tests
go test ./...

# Run the project
go run main.go

# Build for different platforms
go build -o chat-scraper-cli-linux   # Linux
go build -o chat-scraper-cli-windows # Windows
go build -o chat-scraper-cli-darwin  # macOS
```

## Troubleshooting

### Common Issues

1. **"command not found" in WSL:**
   - Check if Go is installed: `go version`
   - Check if `$GOPATH/bin` directory is in PATH

2. **"command not found" on Windows:**
   - Check if `WINDOWS_DIR` variable is configured: `echo %WINDOWS_DIR%`
   - Check if the directory is in system PATH
   - Restart terminal after configuring environment variables

3. **Dependency errors:**
   ```bash
   go mod tidy
   go mod download
   ```

## Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the license specified in the `LICENSE` file.

## Support

If you encounter any issues or have questions, please open an issue in the repository.
