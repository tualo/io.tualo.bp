
export CGO_ENABLED=1
export CGO-CFLAGS="-fembed-bitcode"

export GOOS=ios
export GOARCH=arm64
export SDK=iphoneos
# export SDK=iphonesimulator
go build -buildmode c-archive -o outputfilename.a /path/to/gofile/or/folder

go build -v -a -ldflags="-w -s" \
    -gcflags=-trimpath=/opt/homebrew/opt/opencv/lib/ \
    -asmflags=-trimpath=/opt/homebrew/opt/opencv/lib/ \
    -o ./fooapi spikes/mongoapi.go

# $ export SDK_PATH=`xcrun --sdk $SDK --show-sdk# https://gaitatzis.medium.com/compile-golang-as-a-mobile-library-243e38590f23
-path`
# $ export CLANG=`xcrun --sdk $SDK --find clang`
# $ export CARCH="x86_64"  # if compiling for iPhone simulator
# $ export CARCH="arm64"  # if compiling for iPhone
# $ exec $CLANG -arch $CARCH -isysroot $SDK_PATH -mios-version-min=10.0 "$@"
