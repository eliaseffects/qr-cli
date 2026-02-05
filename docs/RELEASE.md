# Release Guide

## Homebrew Tap

1. Create a tap repo: `github.com/eliaseffects/homebrew-tap`.
2. Ensure the tap repo has a `Formula/` directory committed.
3. GoReleaser will publish the formula on tag push (configured in `.goreleaser.yml`).

## Release Smoke Checklist

1. Update `CHANGELOG.md`.
2. Tag and push the release:
   - `git tag vX.Y.Z`
   - `git push origin vX.Y.Z`
3. Verify GitHub Actions release job succeeds.
4. Download a release asset and verify:
   - `qr --version` shows the tag.
   - `qr "https://example.com" -o test.png` creates a valid PNG.
   - `qr "https://example.com" --format svg -o test.svg` creates a valid SVG.
5. Verify Homebrew tap:
   - `brew tap eliaseffects/tap`
   - `brew install qr-cli`
   - `qr --version`
6. Verify `scripts/install.sh` works on macOS/Linux.
7. (Optional) Run `go test ./...` before announcing.
