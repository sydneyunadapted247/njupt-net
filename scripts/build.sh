#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

MODE="${1:-local}"

mkdir -p dist

build_target() {
	local goos="$1"
	local goarch="$2"
	local ext="${3:-}"
	local out="dist/njupt-net-${goos}-${goarch}${ext}"

	echo "Building ${goos}/${goarch} -> ${out}"
	CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
		go build -trimpath -ldflags "-s -w" -o "${out}" ./cmd/njupt-net
}

case "$MODE" in
	local)
		build_target windows amd64 .exe
		build_target linux arm64
		;;
	all)
		build_target windows amd64 .exe
		build_target linux amd64
		build_target linux arm64
		build_target darwin arm64
		;;
	*)
		echo "Unknown mode: $MODE"
		echo "Usage: ./scripts/build.sh [local|all]"
		exit 1
		;;
esac

echo "Done. Artifacts are in ./dist"
