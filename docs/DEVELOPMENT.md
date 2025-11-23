# ffmpeg-statigo Developer Guide

## Why Manual Setup?

Go modules are designed for source code, not 90MB+ binary files. Attempting to use `go get` with large static libraries leads to:
- Read-only module cache errors
- Complex workarounds that violate Go's design principles
- Poor user experience

Our approach is simple and transparent: clone locally, download libraries, use a replace directive.

## Prerequisites

- Go 1.24+
- gcc (for CGO)
- Git

## Installation

Two installation options are available.

### Option A: Git Submodule (Recommended)

Use this approach for better project integration - keeps dependencies within your project structure.

#### 1. Add as Submodule

Add ffmpeg-statigo as a submodule in your project.

```bash
git submodule add https://github.com/linuxmatters/ffmpeg-statigo vendor/ffmpeg-statigo
```

#### 2. Configure Replace Directive

Add the replace directive pointing to the submodule and get the dependency.

```bash
go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=./vendor/ffmpeg-statigo
go get github.com/linuxmatters/ffmpeg-statigo
```

#### 3. Download Static Libraries

Download the FFmpeg static libraries (~90MB).

```bash
cd vendor/ffmpeg-statigo
go run ./cmd/download-lib
cd ../..
```

#### 4. Build Your Project

```bash
go build
```

#### 5. Commit the Submodule

Only commit the submodule reference. The downloaded files (`lib/*/libffmpeg.a`) are automatically ignored by ffmpeg-statigo's `.gitignore`, so they won't be committed to your repo.

```bash
git add .gitmodules vendor/ffmpeg-statigo
git commit -m "Add ffmpeg-statigo as submodule"
```

#### 6. Team Collaboration

When others clone your project:

**Clone with submodules:**

```bash
git clone --recursive <your-project>
```

**Or if already cloned:**

```bash
git submodule update --init --recursive
cd vendor/ffmpeg-statigo
go run ./cmd/download-lib
cd ../..
```

### Option B: External Directory

Use this if you prefer not to use submodules or want to share ffmpeg-statigo across multiple projects.

#### 1. Clone to Shared Location

Create a dedicated directory for ffmpeg-statigo.

```bash
mkdir -p ~/go-libs
git clone https://github.com/linuxmatters/ffmpeg-statigo ~/go-libs/ffmpeg-statigo
```

Download the static libraries:

```bash
cd ~/go-libs/ffmpeg-statigo
go run ./cmd/download-lib
```

#### 2. Configure Replace Directive

In your project, add the replace directive pointing to use the absolute path to ffmpeg-statigo and get the depenency

```bash
go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=$HOME/go-libs/ffmpeg-statigo
go get github.com/linuxmatters/ffmpeg-statigo
```

#### 3. Build Your Project

```bash
go build
```

## Using FFmpeg in Your Code

```go
package main

import (
    "fmt"
    "github.com/linuxmatters/ffmpeg-statigo"
)

func main() {
    version := ffmpeg.AVCodecVersion()
    fmt.Printf("FFmpeg version: %d.%d.%d\n",
        (version>>16)&0xff, (version>>8)&0xff, version&0xff)
}
```

## Building Your Project

```bash
go build
```

That's it! Your project now has access to a complete FFmpeg static library.

## For ffmpeg-statigo Development

If you're working on ffmpeg-statigo itself:

```bash
# Clone and enter the repository
git clone https://github.com/linuxmatters/ffmpeg-statigo.git
cd ffmpeg-statigo

# Download libraries for your platform
go run ./cmd/download-lib

# Build and test
go build ./...
go test ./...
```

## CI/CD Integration

### GitHub Actions with Submodules

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
      with:
        submodules: recursive

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    - name: Download ffmpeg-statigo libraries
      run: |
        cd vendor/ffmpeg-statigo
        go run ./cmd/download-lib

    - name: Build project
      run: go build ./...
```

### GitHub Actions without Submodules

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Setup ffmpeg-statigo
      run: |
        git clone https://github.com/linuxmatters/ffmpeg-statigo
        cd ffmpeg-statigo
        go run ./cmd/download-lib

    - name: Build project
      run: |
        go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=./ffmpeg-statigo
        go build ./...
```

## Troubleshooting

### "cannot find -lffmpeg" Error

This means the linker can't find the FFmpeg library. Check:

1. **Libraries are downloaded**:
   - For submodules: Verify `vendor/ffmpeg-statigo/lib/<platform>_<arch>/libffmpeg.a` exists
   - For external: Verify `~/go-libs/ffmpeg-statigo/lib/<platform>_<arch>/libffmpeg.a` exists
2. **Replace directive is correct**:
   - For submodules: Use relative path `./vendor/ffmpeg-statigo`
   - For external: Path should be absolute or use `$HOME`
3. **Platform matches**: Libraries must match your build target (linux_amd64, darwin_arm64, etc.)

### Cross-Compilation

Download libraries for other platforms:

```bash
# For submodules
cd vendor/ffmpeg-statigo

# For external directory
cd ~/go-libs/ffmpeg-statigo

# Download for Linux ARM64
GOOS=linux GOARCH=arm64 go run ./cmd/download-lib

# Download for macOS Intel
GOOS=darwin GOARCH=amd64 go run ./cmd/download-lib

# Download all platforms
for os in linux darwin; do
  for arch in amd64 arm64; do
    GOOS=$os GOARCH=$arch go run ./cmd/download-lib
  done
done
```

### Module Cache Issues

If you accidentally try `go get` without the replace directive, Go will cache an incomplete version. Clear it with:

```bash
go clean -modcache
```

## Library Details

- **Size**: ~100MB per platform
- **Verification**: SHA256 checksums from GitHub releases
- **Platforms**: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
- **Source**: Built by [github.com/linuxmatters/ffmpeg-builder](https://github.com/linuxmatters/ffmpeg-builder)

## Quick Reference

### Using Git Submodules (Recommended)

```bash
# In your project
git submodule add https://github.com/linuxmatters/ffmpeg-statigo vendor/ffmpeg-statigo
cd vendor/ffmpeg-statigo && go run ./cmd/download-lib && cd ../..
go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=./vendor/ffmpeg-statigo
go get github.com/linuxmatters/ffmpeg-statigo
go build

# For team members
git clone --recursive <project>
cd vendor/ffmpeg-statigo && go run ./cmd/download-lib
```

### External Directory

```bash
# One-time setup
git clone https://github.com/linuxmatters/ffmpeg-statigo ~/go-libs/ffmpeg-statigo
cd ~/go-libs/ffmpeg-statigo && go run ./cmd/download-lib

# In your project
go mod edit -replace github.com/linuxmatters/ffmpeg-statigo=$HOME/go-libs/ffmpeg-statigo
go get github.com/linuxmatters/ffmpeg-statigo
```
