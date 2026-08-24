#!/usr/bin/env bash
# Builds release binaries locally, pins their URLs and hashes in setup.cfg,
# commits and tags the release, pushes, and uploads the binaries to GitHub.
set -euo pipefail

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <tag>  (e.g. $0 v0.3.0)" >&2
    exit 1
fi

TAG="$1"
if [[ ! "$TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Tag must look like v0.3.0, got '$TAG'" >&2
    exit 1
fi
VERSION="${TAG#v}"  # PEP 440 version for setup.cfg
REPO="ahans/bazel-build-file-sorter"

# --- preflight -------------------------------------------------------------

if [[ "$(git branch --show-current)" != "main" ]]; then
    echo "Releases must be made from main." >&2
    exit 1
fi
if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
    echo "Working tree has uncommitted changes." >&2
    exit 1
fi
git fetch --quiet origin main --tags
if [[ "$(git rev-parse HEAD)" != "$(git rev-parse origin/main)" ]]; then
    echo "main is not in sync with origin/main." >&2
    exit 1
fi
if git rev-parse -q --verify "refs/tags/$TAG" >/dev/null; then
    echo "Tag $TAG already exists." >&2
    exit 1
fi
gh auth status >/dev/null

# --- build -----------------------------------------------------------------

make clean
./build-all.sh "$TAG"

# Sanity check: the host binary reports the release version.
host_bin="bazel-build-file-sorter-$(go env GOHOSTOS)-$(go env GOHOSTARCH)"
if [[ "$("./$host_bin" -v)" != "$TAG" ]]; then
    echo "$host_bin reports version '$("./$host_bin" -v)', expected '$TAG'" >&2
    exit 1
fi

# --- setup.cfg -------------------------------------------------------------

# Each entry: <goos> <goarch> <sys_platform> <platform_machine>
# Marker values are what Python reports, not Go's names.
PLATFORMS=(
    "linux   amd64 linux  x86_64"
    "linux   arm64 linux  aarch64"
    "darwin  amd64 darwin x86_64"
    "darwin  arm64 darwin arm64"
    "windows amd64 win32  AMD64"
    "windows arm64 win32  ARM64"
)

assets=()
{
    # Keep everything up to and including "download_scripts =", with the version bumped.
    sed -e '/^download_scripts =/q' -e "s/^version = .*/version = $VERSION/" setup.cfg

    for p in "${PLATFORMS[@]}"; do
        read -r goos goarch sys_platform machine <<< "$p"
        ext=''
        if [[ "$goos" == "windows" ]]; then
            ext='.exe'
        fi
        asset="bazel-build-file-sorter-${goos}-${goarch}${ext}"
        assets+=("$asset")
        sha256=$(shasum -a 256 "$asset" | cut -d' ' -f1)
        cat << EOF
    [bazel-build-file-sorter${ext}]
    group = bazel-build-file-sorter-binary
    marker = sys_platform == "${sys_platform}" and platform_machine == "${machine}"
    url = https://github.com/${REPO}/releases/download/${TAG}/${asset}
    sha256 = ${sha256}
EOF
    done
} > setup.cfg.new
mv setup.cfg.new setup.cfg

sed -i '' "s/^    rev: v[0-9].*/    rev: $TAG/" README.md

# --- commit, tag, publish --------------------------------------------------

git add setup.cfg README.md
git commit -m "Release $TAG"
git tag -a "$TAG" -m "Release $TAG"

echo
git show --stat HEAD
read -r -p "Push $TAG and publish the GitHub release? [y/N] " answer
if [[ "$answer" != "y" ]]; then
    echo "Aborted. Undo with: git tag -d $TAG && git reset --hard HEAD~1" >&2
    exit 1
fi

git push --atomic origin main "$TAG"
gh release create "$TAG" --repo "$REPO" --verify-tag --title "$TAG" --generate-notes "${assets[@]}"
