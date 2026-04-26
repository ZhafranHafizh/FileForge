$ErrorActionPreference = "Stop"

if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

go fmt ./...
go test ./...
go vet ./...
go build -buildvcs=false -o bin/fileforge.exe .
