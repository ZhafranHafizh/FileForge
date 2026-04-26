# FileForge — Offline Local File Conversion & Compression CLI

## 1. Project Overview

Build an open-source terminal-based tool for converting, compressing, and managing files locally and offline.

The tool should be written in **Go** and act as a CLI orchestrator for proven local engines such as:

- Ghostscript
- Poppler
- qpdf
- ImageMagick
- LibreOffice headless
- Tesseract OCR
- optional image optimizers such as jpegoptim, pngquant, oxipng, and cwebp

The main goal is not to build every converter from scratch, but to provide a clean, reliable, cross-platform CLI experience around existing offline tools.

Project name: **FileForge**

---

## 2. Core Principles

The project must follow these principles:

1. **Local-first**
   - No cloud processing.
   - No file upload.
   - No external API calls.
   - All processing must happen on the user's machine.

2. **Offline-capable**
   - After dependencies are installed, the tool should work without internet access.

3. **Open-source friendly**
   - Simple installation.
   - Clear documentation.
   - Clean repository structure.
   - Easy contribution flow.

4. **Safe by default**
   - Never overwrite input files unless explicitly requested.
   - Always validate input and output paths.
   - Provide clear error messages.
   - Use temporary directories safely.

5. **Composable CLI**
   - Commands should be predictable and script-friendly.
   - Support both single-file and batch workflows.

---

## 3. Recommended Tech Stack

### Language

Use:

```txt
Go
```

Reason:

- Produces standalone binaries.
- Good for CLI tools.
- Easy cross-platform builds.
- Good process execution support.
- Great GitHub Actions and GoReleaser support.

### CLI Libraries

Use:

```txt
github.com/spf13/cobra
github.com/spf13/viper
github.com/schollz/progressbar/v3
```

Optional terminal UI polish:

```txt
github.com/charmbracelet/lipgloss
github.com/charmbracelet/bubbles
github.com/charmbracelet/bubbletea
```

For MVP, keep it simple. Use Cobra + progressbar first.

---

## 4. External Local Engines

The CLI should detect and call these local binaries when needed.

### Required for MVP

```txt
qpdf
ghostscript / gs
pdftoppm
pdfinfo
magick / convert
```

### Optional for Later Versions

```txt
libreoffice
tesseract
jpegoptim
pngquant
oxipng
cwebp
exiftool
```

### Engine Responsibilities

| Feature | Engine |
|---|---|
| PDF compress | Ghostscript |
| PDF merge | qpdf |
| PDF split | qpdf |
| PDF to image | Poppler: pdftoppm |
| PDF info | Poppler: pdfinfo |
| Image convert | ImageMagick |
| Image compress | ImageMagick, jpegoptim, pngquant, oxipng, cwebp |
| Image to PDF | ImageMagick or img2pdf alternative |
| Office to PDF | LibreOffice headless |
| OCR | Tesseract |

---

## 5. Initial CLI Design

The binary should be named:

```bash
fileforge
```

### Global Commands

```bash
fileforge --help
fileforge version
fileforge doctor
```

### Convert Commands

```bash
fileforge convert input.png --to jpg --out output.jpg
fileforge convert input.jpg --to webp --out output.webp
fileforge convert input.pdf --to jpg --out ./pages
fileforge convert input.docx --to pdf --out output.pdf
```

### Compress Commands

```bash
fileforge compress image.jpg --quality 80 --out image-compressed.jpg
fileforge compress image.png --level medium --out image-compressed.png
fileforge compress report.pdf --preset ebook --out report-compressed.pdf
```

### PDF Commands

```bash
fileforge pdf merge a.pdf b.pdf --out merged.pdf
fileforge pdf split book.pdf --pages 1-5 --out chapter.pdf
fileforge pdf info document.pdf
```

### OCR Commands

```bash
fileforge ocr scan.pdf --lang eng --out result.txt
fileforge ocr scan.png --lang ind --out result.txt
```

OCR can be added after the MVP.

---

## 6. MVP Scope

Build the MVP with these features only:

### MVP Feature List

1. `fileforge doctor`
   - Check whether required external binaries exist.
   - Print detected versions.
   - Print missing tools with install hints.

2. `fileforge convert`
   - Image to image:
     - jpg
     - jpeg
     - png
     - webp
   - PDF to image:
     - jpg
     - png
   - Image to PDF.

3. `fileforge compress`
   - Compress image.
   - Compress PDF.

4. `fileforge pdf merge`
   - Merge multiple PDFs into one.

5. `fileforge pdf split`
   - Extract page range from PDF.

6. Basic docs:
   - README
   - installation guide
   - examples
   - contribution guide

Do not implement DOCX/PPTX/XLSX, OCR, or advanced batch processing in the first MVP unless the base commands are already stable.

---

## 7. Non-MVP / Future Features

Add these after MVP:

1. Office conversion
   - DOCX to PDF
   - PPTX to PDF
   - XLSX to PDF
   - Requires LibreOffice.

2. PDF to DOCX
   - Use LibreOffice where possible.
   - Warn users that layout may not be perfect.

3. OCR
   - Image to text.
   - Scanned PDF to searchable PDF.
   - Requires Tesseract.

4. Batch processing

```bash
fileforge batch ./input --to webp --out ./output
fileforge batch ./pdfs --compress --preset ebook --out ./compressed
```

5. Presets config

```yaml
pdf:
  compress:
    default_preset: ebook

image:
  compress:
    default_quality: 80
```

6. Metadata removal

```bash
fileforge clean image.jpg --out clean.jpg
fileforge clean document.pdf --out clean.pdf
```

7. Plugin system.

---

## 8. Repository Structure

Create the repository using this structure:

```txt
fileforge/
├── cmd/
│   ├── root.go
│   ├── version.go
│   ├── doctor.go
│   ├── convert.go
│   ├── compress.go
│   └── pdf.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── runner/
│   │   ├── runner.go
│   │   └── runner_test.go
│   ├── engine/
│   │   ├── detector.go
│   │   ├── ghostscript.go
│   │   ├── imagemagick.go
│   │   ├── poppler.go
│   │   └── qpdf.go
│   ├── convert/
│   │   ├── image.go
│   │   ├── pdf.go
│   │   └── image_to_pdf.go
│   ├── compress/
│   │   ├── image.go
│   │   └── pdf.go
│   ├── pdf/
│   │   ├── merge.go
│   │   ├── split.go
│   │   └── info.go
│   ├── config/
│   │   └── config.go
│   ├── validation/
│   │   ├── file.go
│   │   └── format.go
│   └── errors/
│       └── errors.go
├── docs/
│   ├── installation.md
│   ├── usage.md
│   ├── engines.md
│   └── roadmap.md
├── scripts/
│   ├── install-deps-macos.sh
│   ├── install-deps-ubuntu.sh
│   └── install-deps-windows.ps1
├── testdata/
│   ├── sample.pdf
│   ├── sample.jpg
│   └── sample.png
├── .github/
│   └── workflows/
│       ├── test.yml
│       └── release.yml
├── .gitignore
├── go.mod
├── go.sum
├── main.go
├── README.md
├── CONTRIBUTING.md
├── LICENSE
└── Makefile
```

---

## 9. Command Behavior Requirements

### General Rules

- Input file must exist.
- Output path must be validated.
- If output exists, refuse to overwrite unless `--force` is passed.
- Parent output directory should be created if it does not exist.
- Use safe temporary directories.
- Print clear, human-readable errors.
- Return proper exit codes.

### Exit Codes

Use these exit codes:

```txt
0 = success
1 = general error
2 = invalid input
3 = missing dependency
4 = conversion failed
5 = compression failed
```

### Global Flags

Support these global flags:

```bash
--verbose
--quiet
--force
--config
```

### Output Flag

Every command that creates a file should support:

```bash
--out
```

### Force Overwrite

Example:

```bash
fileforge compress input.pdf --out output.pdf --force
```

Without `--force`, never overwrite.

---

## 10. Dependency Detection

Implement:

```bash
fileforge doctor
```

Expected output example:

```txt
FileForge Doctor

Required tools:
[OK] qpdf          11.9.0
[OK] ghostscript   10.03.1
[OK] pdftoppm      24.02.0
[OK] pdfinfo       24.02.0
[OK] magick        7.1.1

Optional tools:
[Missing] libreoffice
[Missing] tesseract

Status:
Ready for MVP features.
```

If required dependency is missing, show install suggestions.

Example:

```txt
Missing dependency: qpdf

Install:
macOS:   brew install qpdf
Ubuntu:  sudo apt install qpdf
Windows: winget install qpdf
```

---

## 11. Engine Implementation Notes

### Runner Package

Create a reusable process runner.

Responsibilities:

- Execute local binaries.
- Capture stdout.
- Capture stderr.
- Return exit code.
- Support verbose logging.
- Prevent shell injection by using `exec.Command` with argument arrays.
- Do not concatenate shell strings.

Example concept:

```go
runner.Run(ctx, "qpdf", []string{"--version"})
```

Do not use:

```go
exec.Command("sh", "-c", userInput)
```

### Engine Detector

Create functions:

```go
IsAvailable(binary string) bool
Version(binary string, args ...string) (string, error)
Require(binary string) error
```

### Path Validation

Create helpers:

```go
EnsureInputFile(path string) error
EnsureOutputPath(path string, force bool) error
EnsureOutputDir(path string) error
Extension(path string) string
```

---

## 12. Feature Implementation Details

## 12.1 Image Convert

Command:

```bash
fileforge convert input.png --to jpg --out output.jpg
```

Engine:

```txt
ImageMagick
```

Command concept:

```bash
magick input.png output.jpg
```

Requirements:

- Support jpg, jpeg, png, webp.
- Validate `--to`.
- Validate output extension.
- Refuse unsupported formats.
- Never overwrite without `--force`.

---

## 12.2 PDF to Image

Command:

```bash
fileforge convert input.pdf --to jpg --out ./pages
fileforge convert input.pdf --to png --out ./pages
```

Engine:

```txt
pdftoppm
```

Command concept:

```bash
pdftoppm -jpeg input.pdf ./pages/page
pdftoppm -png input.pdf ./pages/page
```

Requirements:

- If output is directory, create it.
- Output files should be named predictably:

```txt
page-1.jpg
page-2.jpg
page-3.jpg
```

Optional flags:

```bash
--dpi 150
--first-page 1
--last-page 5
```

Default DPI:

```txt
150
```

---

## 12.3 Image to PDF

Command:

```bash
fileforge convert image.jpg --to pdf --out image.pdf
```

Engine:

```txt
ImageMagick
```

Command concept:

```bash
magick image.jpg image.pdf
```

Support:

- jpg
- jpeg
- png
- webp

---

## 12.4 PDF Compress

Command:

```bash
fileforge compress report.pdf --preset ebook --out report-compressed.pdf
```

Engine:

```txt
Ghostscript
```

Preset mapping:

```txt
screen   = low quality, smallest size
ebook    = medium quality
printer  = high quality
prepress = highest quality
default  = general purpose
```

Command concept:

```bash
gs -sDEVICE=pdfwrite \
   -dCompatibilityLevel=1.4 \
   -dPDFSETTINGS=/ebook \
   -dNOPAUSE \
   -dQUIET \
   -dBATCH \
   -sOutputFile=output.pdf \
   input.pdf
```

Requirements:

- Validate preset.
- Output file should be smaller if possible, but do not guarantee.
- Show before/after file size.
- Show percentage reduction.
- If output is larger, warn the user.

---

## 12.5 Image Compress

Command:

```bash
fileforge compress image.jpg --quality 80 --out image-compressed.jpg
fileforge compress image.png --level medium --out image-compressed.png
```

Engine MVP:

```txt
ImageMagick
```

Command concept:

```bash
magick input.jpg -quality 80 output.jpg
```

Requirements:

- Support jpg, jpeg, png, webp.
- `--quality` default: 80.
- Valid quality range: 1-100.
- Show before/after file size.
- Show percentage reduction.

Later, use specialized optimizers:

- jpegoptim for JPEG
- pngquant or oxipng for PNG
- cwebp for WEBP

---

## 12.6 PDF Merge

Command:

```bash
fileforge pdf merge a.pdf b.pdf c.pdf --out merged.pdf
```

Engine:

```txt
qpdf
```

Command concept:

```bash
qpdf --empty --pages a.pdf b.pdf c.pdf -- merged.pdf
```

Requirements:

- Require at least two input PDFs.
- Validate all inputs.
- Validate output path.
- Preserve input order.

---

## 12.7 PDF Split

Command:

```bash
fileforge pdf split book.pdf --pages 1-5 --out chapter.pdf
```

Engine:

```txt
qpdf
```

Command concept:

```bash
qpdf book.pdf --pages book.pdf 1-5 -- chapter.pdf
```

Requirements:

- Validate page range format.
- Support:
  - `1`
  - `1-5`
  - `1,3,5`
  - `1-3,7,10-12`
- Validate output path.
- Print clear errors for invalid range.

---

## 12.8 PDF Info

Command:

```bash
fileforge pdf info document.pdf
```

Engine:

```txt
pdfinfo
```

Command concept:

```bash
pdfinfo document.pdf
```

Requirements:

- Print page count.
- Print PDF size.
- Print creation date if available.
- Print encrypted status if available.

---

## 13. Testing Plan

### Unit Tests

Add unit tests for:

- Path validation.
- Format detection.
- Output overwrite behavior.
- Page range parsing.
- Engine detection.
- Command argument generation.

### Integration Tests

Integration tests should be optional because they require external binaries.

Use build tags:

```bash
go test ./...
go test -tags=integration ./...
```

Integration tests should cover:

- Image conversion.
- PDF compression.
- PDF merge.
- PDF split.
- PDF to image.

### Test Data

Put small test files in:

```txt
testdata/
```

Use tiny sample files to keep repository lightweight.

---

## 14. Documentation Requirements

Create these docs:

### README.md

Must include:

- Project description.
- Feature list.
- Installation.
- Dependency requirements.
- Quick start.
- Usage examples.
- Roadmap.
- License.

### docs/installation.md

Include:

```txt
macOS
Ubuntu/Debian
Windows
```

### docs/engines.md

Explain what external engines are used and why.

### docs/usage.md

Include all command examples.

### docs/roadmap.md

Split into:

```txt
MVP
v0.2
v0.3
v1.0
```

### CONTRIBUTING.md

Include:

- How to set up development.
- How to run tests.
- How to add a new command.
- Coding conventions.
- Pull request checklist.

---

## 15. GitHub Actions

Create:

```txt
.github/workflows/test.yml
.github/workflows/release.yml
```

### test.yml

Run:

```bash
go test ./...
go vet ./...
go fmt check
```

On:

```txt
push
pull_request
```

### release.yml

Use GoReleaser later.

For MVP, create a simple workflow that builds binaries for:

```txt
linux-amd64
linux-arm64
darwin-amd64
darwin-arm64
windows-amd64
```

---

## 16. Makefile

Create a Makefile with:

```makefile
APP_NAME=fileforge

.PHONY: build run test fmt vet clean

build:
	go build -o bin/$(APP_NAME) .

run:
	go run .

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin/
```

---

## 17. Suggested Milestones

## Milestone 1 — Project Skeleton

Tasks:

- Initialize Go module.
- Add Cobra.
- Create root command.
- Add version command.
- Add Makefile.
- Add README draft.
- Add basic GitHub Actions test workflow.

Definition of done:

```bash
fileforge --help
fileforge version
go test ./...
```

---

## Milestone 2 — Engine Detection

Tasks:

- Add runner package.
- Add engine detector.
- Implement `fileforge doctor`.
- Detect qpdf, ghostscript, pdftoppm, pdfinfo, and magick.
- Add install hints.

Definition of done:

```bash
fileforge doctor
```

prints dependency status clearly.

---

## Milestone 3 — Validation Layer

Tasks:

- Add path validation helpers.
- Add format validation helpers.
- Add overwrite protection.
- Add file size helper.
- Add page range parser.

Definition of done:

- Invalid files return clear errors.
- Existing output files are protected.
- Page range parser is tested.

---

## Milestone 4 — Image Convert & Compress

Tasks:

- Implement image to image conversion.
- Implement image compression.
- Add quality flag.
- Add before/after file size display.
- Add tests for command argument generation.

Definition of done:

```bash
fileforge convert input.png --to jpg --out output.jpg
fileforge compress input.jpg --quality 80 --out compressed.jpg
```

---

## Milestone 5 — PDF Core Features

Tasks:

- Implement PDF compression.
- Implement PDF merge.
- Implement PDF split.
- Implement PDF info.
- Add tests for command argument generation.

Definition of done:

```bash
fileforge compress input.pdf --preset ebook --out compressed.pdf
fileforge pdf merge a.pdf b.pdf --out merged.pdf
fileforge pdf split input.pdf --pages 1-3 --out output.pdf
fileforge pdf info input.pdf
```

---

## Milestone 6 — PDF/Image Conversion

Tasks:

- Implement PDF to image.
- Implement image to PDF.
- Add DPI option.
- Add page range options for PDF to image.

Definition of done:

```bash
fileforge convert input.pdf --to jpg --out ./pages
fileforge convert image.jpg --to pdf --out image.pdf
```

---

## Milestone 7 — Documentation & Release Prep

Tasks:

- Complete README.
- Complete installation docs.
- Complete usage docs.
- Add MIT or Apache-2.0 license.
- Add release workflow.
- Create v0.1.0 tag.

Definition of done:

- New user can install dependencies.
- New user can build the tool.
- New user can run all MVP commands.
- GitHub release includes binaries.

---

## 18. Coding Guidelines

Follow these rules:

1. Keep command logic thin.
2. Put business logic inside `internal/`.
3. Do not call external binaries directly from command files.
4. Use context-aware process execution.
5. Avoid shell string execution.
6. Prefer explicit errors.
7. Keep functions small and testable.
8. Add tests for parsing and validation.
9. Use standard Go formatting.
10. Do not introduce unnecessary abstractions too early.

---

## 19. Security Guidelines

Important:

- Never pass untrusted user input through shell strings.
- Use `exec.CommandContext`.
- Use argument slices, not shell concatenation.
- Validate all file paths.
- Do not overwrite files by default.
- Avoid writing temp files outside safe temp directories.
- Clean up temporary files.
- Do not send telemetry.
- Do not make network requests.
- Do not upload files anywhere.

---

## 20. Initial Commands for Codex CLI

Use these instructions to start implementation.

```txt
Create a Go CLI project named FileForge.

Use Cobra for the CLI structure.

Implement the repository structure described in this document.

Start with:
1. Go module initialization.
2. Root command.
3. Version command.
4. Doctor command.
5. Runner package.
6. Engine detection package.
7. Basic validation helpers.
8. README draft.
9. Makefile.
10. GitHub Actions test workflow.

Do not implement every feature at once.
Prioritize clean structure and Milestone 1-2 first.
```

---

## 21. First Implementation Prompt for Codex CLI

Use this prompt:

```txt
You are working on an open-source Go CLI tool named FileForge.

FileForge is an offline, local-first file conversion and compression CLI. It wraps local binaries such as qpdf, ghostscript, poppler, and ImageMagick.

Please implement the initial project skeleton:

- Initialize a Go module.
- Use Cobra for command structure.
- Create main.go.
- Create cmd/root.go.
- Create cmd/version.go.
- Create cmd/doctor.go.
- Create internal/runner/runner.go.
- Create internal/engine/detector.go.
- Create internal/validation/file.go.
- Create README.md.
- Create Makefile.
- Create .github/workflows/test.yml.

The CLI binary should be named fileforge.

Commands required now:

1. fileforge --help
2. fileforge version
3. fileforge doctor

The doctor command should check whether these binaries are available:

- qpdf
- gs
- pdftoppm
- pdfinfo
- magick

It should print OK or Missing for each one.

Use exec.CommandContext with argument slices.
Do not use shell execution.
Do not implement conversion yet.
Keep the code clean, idiomatic, and testable.
```

---

## 22. Second Implementation Prompt for Codex CLI

After Milestone 1 and 2 are done, use this prompt:

```txt
Continue implementing FileForge.

Add validation and utility packages:

- file existence validation
- output path validation
- overwrite protection
- output directory creation
- extension detection
- file size helper
- page range parser

Add tests for:

- existing input file
- missing input file
- output overwrite protection
- valid page ranges
- invalid page ranges

Do not implement conversion yet.
Focus on correctness and clean tests.
```

---

## 23. Third Implementation Prompt for Codex CLI

After validation is done, use this prompt:

```txt
Continue implementing FileForge.

Implement image conversion and image compression.

Commands:

fileforge convert input.png --to jpg --out output.jpg
fileforge convert input.jpg --to webp --out output.webp
fileforge compress input.jpg --quality 80 --out compressed.jpg

Requirements:

- Use ImageMagick through the magick binary.
- Support jpg, jpeg, png, and webp.
- Validate input file.
- Validate output path.
- Do not overwrite existing files unless --force is passed.
- Show before and after file size for compression.
- Use internal runner package for process execution.
- Add tests for command argument generation and validation.
```

---

## 24. Fourth Implementation Prompt for Codex CLI

After image features are done, use this prompt:

```txt
Continue implementing FileForge.

Implement PDF core commands.

Commands:

fileforge compress input.pdf --preset ebook --out compressed.pdf
fileforge pdf merge a.pdf b.pdf --out merged.pdf
fileforge pdf split input.pdf --pages 1-3 --out output.pdf
fileforge pdf info input.pdf

Requirements:

- Use Ghostscript for PDF compression.
- Use qpdf for merge and split.
- Use pdfinfo for info.
- Support PDF compression presets:
  - screen
  - ebook
  - printer
  - prepress
  - default
- Validate input files.
- Validate output paths.
- Do not overwrite existing files unless --force is passed.
- Show before and after file size for compression.
- Add tests for command argument generation and validation.
```

---

## 25. Fifth Implementation Prompt for Codex CLI

After PDF core commands are done, use this prompt:

```txt
Continue implementing FileForge.

Implement PDF/image conversion.

Commands:

fileforge convert input.pdf --to jpg --out ./pages
fileforge convert input.pdf --to png --out ./pages
fileforge convert image.jpg --to pdf --out image.pdf

Requirements:

- Use pdftoppm for PDF to image.
- Use ImageMagick for image to PDF.
- Support --dpi, default 150.
- Support --first-page and --last-page for PDF to image.
- Create output directory if it does not exist.
- Use predictable page output naming.
- Validate input files.
- Validate output paths.
- Use internal runner package.
- Add tests for command argument generation and validation.
```

---

## 26. License Recommendation

Use one of these:

### MIT

Best if you want maximum adoption and simple usage.

### Apache-2.0

Best if you want explicit patent protection.

Recommendation:

```txt
MIT License
```

---

## 27. Final Expected MVP Commands

By the end of MVP, these commands should work:

```bash
fileforge --help
fileforge version
fileforge doctor

fileforge convert input.png --to jpg --out output.jpg
fileforge convert input.jpg --to webp --out output.webp
fileforge convert input.pdf --to jpg --out ./pages
fileforge convert image.jpg --to pdf --out image.pdf

fileforge compress input.jpg --quality 80 --out compressed.jpg
fileforge compress input.pdf --preset ebook --out compressed.pdf

fileforge pdf info input.pdf
fileforge pdf merge a.pdf b.pdf --out merged.pdf
fileforge pdf split input.pdf --pages 1-3 --out output.pdf
```

---

## 28. Notes for Future Maintainers

This project should remain a wrapper/orchestrator, not a full reimplementation of PDF/image/office engines.

The value of FileForge is:

- clean CLI UX
- local/offline guarantee
- safe defaults
- useful presets
- cross-platform documentation
- consistent command structure
- automation-friendly behavior

Keep the core simple and reliable before adding advanced features.
