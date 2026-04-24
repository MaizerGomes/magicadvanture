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

# Get checksums
echo "Getting checksums..."
ARM_SHA=$(curl -sL "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-arm" | shasum -a256 | cut -d' ' -f1)
INTEL_SHA=$(curl -sL "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-64" | shasum -a256 | cut -d' ' -f1)

echo "ARM SHA: $ARM_SHA"
echo "Intel SHA: $INTEL_SHA"

# Update formula
cat > /tmp/magicadventure.rb << EOF
class Magicadventure < Formula
  desc "A terminal RPG with slot-based saves, turn-based combat, and optional multiplayer"
  homepage "https://github.com/MaizerGomes/magicadvanture"
  version "$VERSION"
  url "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-arm"
  sha256 "$ARM_SHA"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/MaizerGomes/magicadvanture/releases/download/v$VERSION/magicadventure-mac-64"
      sha256 "$INTEL_SHA"
    end
  end

  def install
    bin.install "magicadventure-mac-arm" => "magicadventure"
  end

  test do
    assert shell_output("#{bin}/magicadventure").include?("Magic Adventure")
  end
end
EOF

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