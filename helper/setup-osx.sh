xcode-select --install
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
(echo; echo 'eval "$(/opt/homebrew/bin/brew shellenv)"') >> ~/.zprofile
    eval "$(/opt/homebrew/bin/brew shellenv)"
brew install tesseract
brew install tesseract-lang
brew install leptonica
brew install zbar
brew install opencv


# xattr -d com.apple.quarantine /path/to/file

