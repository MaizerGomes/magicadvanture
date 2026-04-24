#!/bin/bash
# Homebrew Deploy Script for Magic Adventure
# Usage: ./deploy-homebrew.sh <version>

set -e

VERSION=${1:-latest}

if [ "$VERSION" = "latest" ]; then
    VERSION=$(gh release list --limit 1 --jq '.[0].tagName')
    VERSION=${VERSION#v}
fi

echo "Deploying Magic Adventure v$VERSION to Homebrew..."

# Verify release exists
RELEASE_URL="https://github.com/MaizerGomes/magicadvanture/releases/tag/v$VERSION"
if ! curl -sI "$RELEASE_URL" | head -1 | grep -q "200"; then
    echo "Error: Release v$VERSION not found"
    exit 1
fi

# Get checksums with redirect following
echo "Getting checksums..."

get_checksum() {
    local url=$1
    curl -sL --fail "$url" -o /tmp/magicadventure-temp || {
        echo "Error downloading: $url"
        exit 1
    }
    shasum -a256 /tmp/magicadventure-temp | cut -d' ' -f1
    rm -f /tmp/magicadventure-temp
}

ARM_SHA=$(get_checksum "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-arm")
INTEL_SHA=$(get_checksum "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-64")

echo "ARM SHA: $ARM_SHA"
echo "Intel SHA: $INTEL_SHA"

# Verify checksums are valid (not empty, not error messages)
if [ ${#ARM_SHA} -ne 64 ] || [ ${#INTEL_SHA} -ne 64 ]; then
    echo "Error: Invalid checksum received"
    exit 1
fi

# Update formula
cat > /tmp/magicadventure.rb << EOF
class Magicadventure < Formula
  desc "A terminal RPG with slot-based saves, turn-based combat, and optional multiplayer"
  homepage "https://github.com/MaizerGomes/magicadvanture"
  version "$VERSION"
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-arm"
      sha256 "$ARM_SHA"
    else
      url "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-64"
      sha256 "$INTEL_SHA"
    end
  end

  def install
    binary = Dir["magicadventure-*"]
      .reject { |path| File.directory?(path) }
      .first
    raise "magicadventure binary not found in buildpath" unless binary

    bin.install binary => "magicadventure"
  end

  test do
    assert shell_output("#{bin}/magicadventure").include?("Magic Adventure")
  end
end
EOF

echo "Formula created. Pushing to tap..."

# Commit to tap repo
cd /tmp
rm -rf homebrew-tap-deploy
git clone git@github.com:MaizerGomes/homebrew-magicadvanture.git homebrew-tap-deploy
cd homebrew-tap-deploy
cp /tmp/magicadventure.rb .
git add magicadventure.rb
git commit -m "Update to v$VERSION"
git push origin master

echo "✓ Done! Homebrew updated to v$VERSION"
echo ""
echo "Users can now run: brew upgrade magicadventure"
