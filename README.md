## Converter

This is a simple converter, that converts images in the given directory to jpeg and png, webp.

## How to use

1. Clone this repository
2. Install requirements
3. Run the script

```bash
git clone
cd converter
## Install libvips for linux
sudo apt-get install libvips
## Install libvips for mac
brew install vips
go get -u github.com/davidbyttow/govips/v2/vips
go get -u golang.org/x/image/webp
go mod tidy
## example
go run main.go -input=jpeg -output=webp ./test/jpeg
```

