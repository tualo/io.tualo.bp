# M1 Problems

## OpenCV IOS

`sudo ln -s /opt/homebrew/bin/python3	/opt/homebrew/bin/python`

---> https://github.com/opencv/opencv/issues/23738
IPHONEOS_DEPLOYMENT_TARGET=14.0 python3 build_framework.py ios --contrib opencv_contrib --iphoneos_archs arm64 x86_64
IPHONEOS_DEPLOYMENT_TARGET=17.0 python3 build_framework.py ios --contrib opencv_contrib --iphoneos_archs arm64 --iphonesimulator_archs x86_64
