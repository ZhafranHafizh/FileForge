# FileForge

FileForge is a local-first, offline-capable CLI for file conversion, compression, and PDF utilities. It wraps trusted local binaries like ImageMagick, Ghostscript, qpdf, and Poppler so files stay on your machine.

## Local / Offline Guarantee

- No uploads
- No cloud processing
- No telemetry
- No network calls during normal use
- All processing uses local binaries already installed on your system

## Current Features

- Image to image conversion:
  - `jpg`
  - `jpeg`
  - `png`
  - `webp`
- Image compression
- PDF compression
- PDF merge
- PDF split
- PDF info
- PDF to image
- Image to PDF
- Interactive terminal mode
- Dependency doctor command

## Supported Commands

### Command Mode

```bash
fileforge convert input.png --to jpg --out output.jpg
fileforge convert input.jpg --to webp --out output.webp
fileforge convert input.pdf --to jpg --out ./pages
fileforge convert input.pdf --to png --out ./pages --dpi 200 --first-page 1 --last-page 3
fileforge convert image.jpg --to pdf --out image.pdf

fileforge compress input.jpg --quality 80 --out compressed.jpg
fileforge compress input.pdf --preset ebook --out compressed.pdf

fileforge pdf merge a.pdf b.pdf c.pdf --out merged.pdf
fileforge pdf split input.pdf --pages 1-3 --out chapter.pdf
fileforge pdf info input.pdf

fileforge doctor
fileforge version
```

### Interactive Mode

Start the wizard with either:

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

Not implemented yet in interactive mode:

- OCR
- Office conversion

## Required Local Engines

FileForge uses these binaries for currently implemented features:

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

## Install Dependencies

### Windows

Using `winget`:

```powershell
winget install ImageMagick.ImageMagick
winget install ArtifexSoftware.GhostScript
winget install QPDF.QPDF
winget install oschwartz10612.Poppler
```

Ghostscript may be installed as `gswin64c` or `gswin32c`. FileForge checks those automatically.

### macOS

Using Homebrew:

```bash
brew install imagemagick ghostscript qpdf poppler
```

### Ubuntu / Debian

Using `apt`:

```bash
sudo apt update
sudo apt install imagemagick ghostscript qpdf poppler-utils
```

## Build From Source

### Standard Go Build

```bash
go build -buildvcs=false -o bin/fileforge .
```

### Windows Binary

```powershell
go build -buildvcs=false -o bin/fileforge.exe .
```

If you hit a VCS stamping error such as `error obtaining VCS status`, use the `-buildvcs=false` form above. The included scripts and Makefile already do this.

## Development Commands

### Windows PowerShell

```powershell
.\scripts\build-windows.ps1
```

This runs:

```powershell
go fmt ./...
go test ./...
go vet ./...
go build -buildvcs=false -o bin/fileforge.exe .
```

### macOS / Linux

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

## Running FileForge

### From source

```bash
go run .
```

### Build and run on Windows

```powershell
go build -buildvcs=false -o bin/fileforge.exe .
.\bin\fileforge.exe
```

### Build and run on macOS / Linux

```bash
go build -buildvcs=false -o bin/fileforge .
./bin/fileforge
```

## Doctor Output

`fileforge doctor` checks all required engines for currently implemented features and explains why each tool matters. Example categories:

- `magick`: image conversion, image compression, image to PDF
- `qpdf`: PDF merge and PDF split
- `ghostscript`: PDF compression
- `pdftoppm`: PDF to image
- `pdfinfo`: PDF info

## Known Limitations

- OCR is not implemented yet
- Office conversion is not implemented yet
- PDF to image currently keeps Poppler’s generated `page-*` naming
- PDF to image overwrite protection is directory/prefix-based
- Release packaging is still manual
- No metadata cleanup yet
- No batch processing yet

## Roadmap

- OCR with Tesseract
- Office conversion with LibreOffice
- Metadata cleanup
- Batch processing
- Better release packaging

## License

License file is not added yet. The implementation plan currently recommends MIT.
