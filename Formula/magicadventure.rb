class Magicadventure < Formula
  desc "A text adventure game"
  homepage "https://github.com/MaizerGomes/magicadvanture"
  version "6.0"
  url "https://github.com/MaizerGomes/magicadvanture/releases/download/v6.0/magicadventure-mac-arm"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/MaizerGomes/magicadvanture/releases/download/v6.0/magicadventure-mac-64"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    else
      url "https://github.com/MaizerGomes/magicadvanture/releases/download/v6.0/magicadventure-mac-arm"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
  end

  def install
    bin.install "magicadventure"
  end

  test do
    assert shell_output("#{bin}/magicadventure").include?("Magic Adventure")
  end
end