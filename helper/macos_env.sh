export CGO_CXXFLAGS="--std=c++11"
export CGO_ENABLED=1
export GOARCH="arm64"
export LIBRARY_PATH="/opt/homebrew/lib"
export CPATH="/opt/homebrew/include"



export CGO_CXXFLAGS="--std=c++11"
export CGO_CPPFLAGS="-I/opt/homebrew/include/opencv4"
export CGO_LDFLAGS="-L/opt/homebrew/lib/opencv4 -lopencv_stitching -lopencv_superres -lopencv_videostab -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_dnn_objdetect -lopencv_dpm -lopencv_face -lopencv_photo -lopencv_fuzzy -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_optflow -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_stereo -lopencv_structured_light -lopencv_phase_unwrapping -lopencv_surface_matching -lopencv_tracking -lopencv_datasets -lopencv_dnn -lopencv_plot -lopencv_xfeatures2d -lopencv_shape -lopencv_video -lopencv_ml -lopencv_ximgproc -lopencv_calib3d -lopencv_features2d -lopencv_highgui -lopencv_videoio -lopencv_flann -lopencv_xobjdetect -lopencv_imgcodecs -lopencv_objdetect -lopencv_xphoto -lopencv_imgproc -lopencv_core"

export PATH=$PATH:~/go/bin
export PKG_CONFIG_PATH="$(brew --prefix opencv)/lib/pkgconfig:$PKG_CONFIG_PATH"

# Setze DYLD_LIBRARY_PATH für dynamische Bibliotheken
export DYLD_LIBRARY_PATH="$(brew --prefix opencv)/lib:$(brew --prefix)/lib:$DYLD_LIBRARY_PATH"

# CGO Flags für Go - verwende pkg-config für korrekte Flags
export CGO_CPPFLAGS="$(pkg-config --cflags opencv4)"
export CGO_LDFLAGS="$(pkg-config --libs opencv4)"

echo "✓ Umgebungsvariablen für OpenCV gesetzt"
echo "PKG_CONFIG_PATH: $PKG_CONFIG_PATH"
echo "CGO_CPPFLAGS: $CGO_CPPFLAGS"
echo ""

# /opt/homebrew/Library/Taps/homebrew/homebrew-core/Formula/i/io.tualo.bp.rb
# otool -L io.tualo.bp | grep homebrew | awk '{print $1}'
# go build -buildmode=c-shared -trimpath .
