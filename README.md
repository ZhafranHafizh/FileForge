# FileForge

Local-first, offline file conversion and compression from your terminal.

![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)
![Local-first](https://img.shields.io/badge/Local--first-yes-success)
![Offline](https://img.shields.io/badge/Offline-ready-success)

FileForge is a Go-based CLI for converting, compressing, and managing files locally without uploading them to the cloud. It wraps proven local tools such as ImageMagick, Ghostscript, qpdf, and Poppler to provide a consistent command-line workflow for everyday file tasks.

FileForge supports both:

- command mode
- interactive terminal wizard mode

## Key Features

- Image to image conversion
- Image compression
- PDF compression
- PDF merge
- PDF split
- PDF info
- PDF to image
- Image to PDF
- Interactive terminal wizard
- Output folder mode with `--output-dir`
- Local/offline processing

## Why FileForge?

- No cloud upload: files stay on your machine.
- Privacy-friendly: useful for sensitive documents and local workflows.
- Scriptable: predictable CLI commands for automation and shell usage.
- Offline-ready: after dependencies are installed, it works without internet access.
- Practical: useful for developers, students, sysadmins, and content creators.

## Installation

### Build from source

```bash
git clone https://github.com/YOUR_USERNAME/fileforge.git
cd fileforge
go build -buildvcs=false -o bin/fileforge .
```

On Windows:

```powershell
go build -buildvcs=false -o bin/fileforge.exe .
```

If Go reports a VCS stamping error such as `error obtaining VCS status`, use the `-buildvcs=false` form above. FileForge's Makefile and helper scripts already do this.

### Windows dependencies

Install required local engines with `winget`:

```powershell
winget install ImageMagick.ImageMagick
winget install ArtifexSoftware.GhostScript
winget install QPDF.QPDF
winget install oschwartz10612.Poppler
```

Notes:

- Ghostscript may be installed as `gswin64c` or `gswin32c`.
- FileForge checks `gs`, `gswin64c`, and `gswin32c` automatically.

### macOS dependencies

Install required local engines with Homebrew:

```bash
brew install imagemagick ghostscript qpdf poppler
```

### Ubuntu/Debian dependencies

Install required local engines with `apt`:

```bash
sudo apt update
sudo apt install imagemagick ghostscript qpdf poppler-utils
```

## Dependency Overview

FileForge currently depends on these local binaries:

- `magick`
  - image conversion
  - image compression
  - image to PDF
- Ghostscript: `gs`, `gswin64c`, or `gswin32c`
  - PDF compression
- `qpdf`
  - PDF merge
  - PDF split
- `pdftoppm`
  - PDF to image
- `pdfinfo`
  - PDF info

Use `fileforge doctor` to verify your setup.

## Command Mode

### Convert

```bash
fileforge convert input.png --to jpg --out output.jpg
fileforge convert input.jpg --to webp --output-dir ./output
fileforge convert input.pdf --to jpg --output-dir ./output
fileforge convert input.pdf --to png --output-dir ./output --dpi 200 --first-page 1 --last-page 3
fileforge convert photo.jpg --to pdf --output-dir ./output
```

### Compress

```bash
fileforge compress image.jpg --quality 80 --output-dir ./output
fileforge compress report.pdf --preset ebook --output-dir ./output
```

### PDF Tools

```bash
fileforge pdf merge a.pdf b.pdf --output-dir ./output
fileforge pdf split book.pdf --pages 1-5 --output-dir ./output
fileforge pdf info document.pdf
```

### Utilities

```bash
fileforge doctor
fileforge version
```

## Interactive Mode

Start the interactive wizard with either:

```bash
go run .
```

or:

```bash
fileforge
fileforge interactive
```

Interactive mode currently supports:

- Image to image conversion
- Image compression
- PDF compression
- PDF merge
- PDF split
- PDF info
- PDF to image
- Image to PDF
- Doctor / dependency check

Interactive mode is a guided terminal wizard with grouped menus and local-first workflows. It uses the same internal services as command mode, so you get the same behavior whether you script FileForge or step through it manually.

Interactive mode includes:

- Grouped menus for Convert, Compress, PDF Tools, System Check, and Coming Soon
- Drag-and-drop path support for files dropped into the terminal
- A default output folder of `./FileForge-Output` when you leave the output prompt empty
- Styled summaries before execution, plus overwrite confirmation and final confirmation
- Fully local and offline processing after dependencies are installed

## Output Folder Mode

FileForge supports both explicit output paths and generated output locations.

- Use `--out` when you want to provide the exact output file or directory.
- Use `--output-dir` when you want FileForge to generate sensible output names inside one folder.
- Do not use both at the same time.

Examples:

```bash
fileforge compress image.jpg --quality 80 --output-dir ./output
fileforge convert photo.png --to webp --output-dir ./output
fileforge pdf merge a.pdf b.pdf --output-dir ./output
```

After a successful operation, FileForge prints:

```text
The process is complete. Please access the following path to view the results.
<output-path>
```

## Running FileForge

From source:

```bash
go run .
```

Build and run on Windows:

```powershell
go build -buildvcs=false -o bin/fileforge.exe .
.\bin\fileforge.exe
```

Build and run on macOS/Linux:

```bash
go build -buildvcs=false -o bin/fileforge .
./bin/fileforge
```

## Development

### Windows PowerShell

```powershell
.\scripts\build-windows.ps1
```

### macOS/Linux

```bash
sh ./scripts/build-unix.sh
```

### Makefile

```bash
make fmt
make test
make vet
make build
```

## Doctor

`fileforge doctor` checks the local engines required for currently implemented features and explains what each tool is used for.

Example coverage:

- `magick`: image conversion, image compression, image to PDF
- `qpdf`: PDF merge and PDF split
- `ghostscript`: PDF compression
- `pdftoppm`: PDF to image
- `pdfinfo`: PDF info

## Known Limitations

- OCR is not implemented yet
- Office conversion is not implemented yet
- PDF to image currently uses Poppler's generated `page-*` filenames
- PDF to image overwrite protection is directory/prefix based
- Release packaging is still basic

## Roadmap

- OCR with Tesseract
- Office conversion with LibreOffice
- Metadata cleanup
- Batch processing
- Better release packaging

## License

MIT. See [LICENSE](LICENSE).
