# FileForge

FileForge is a local-first, offline-capable CLI for file conversion, compression, and PDF utilities. It wraps proven local binaries such as qpdf, Ghostscript, Poppler, and ImageMagick instead of sending files to any external service.

This repository currently implements Milestone 1 and Milestone 2 from `fileforge_implementation_plan.md`:

- Cobra-based CLI skeleton
- `fileforge version`
- `fileforge doctor`
- reusable process runner
- engine detection helpers
- basic file validation helpers

Conversion, compression, merge, split, OCR, and office conversion are intentionally not implemented yet.

## Requirements

- Go 1.22+
- Local binaries for feature detection:
  - `qpdf`
  - `gs`
  - `pdftoppm`
  - `pdfinfo`
  - `magick`

## Installation Hints

Doctor will print platform-specific hints when dependencies are missing. Common install commands:

- macOS: `brew install qpdf ghostscript poppler imagemagick`
- Ubuntu/Debian: `sudo apt install qpdf ghostscript poppler-utils imagemagick`
- Windows: `winget install QPDF.QPDF ArtifexSoftware.GhostScript oschwartz10612.Poppler ImageMagick.ImageMagick`

## Quick Start

```bash
go run . --help
go run . version
go run . doctor
```

## Commands Available Now

```bash
fileforge --help
fileforge version
fileforge doctor
```

## Global Flags

All current commands inherit these global flags:

- `--verbose`
- `--quiet`
- `--force`
- `--config`

## Development

```bash
make fmt
make test
make vet
make build
```

## Roadmap

Next planned milestone:

- Milestone 3: validation layer expansion, including page range parsing and more format helpers

After that:

- Milestone 4: image convert and compress
- Milestone 5: PDF core features
- Milestone 6: PDF/image conversion

## License

License has not been added yet. The implementation plan currently recommends MIT.
