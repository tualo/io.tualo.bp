#/bin/bash

entry_dir=`pwd`
train_data_list="deu eng oci osd"
train_data_dir="/opt/homebrew/share/tessdata/"
working_dir="./Scanner.app/Contents/MacOS/"
executable="io.tualo.bp"

fyne package -os darwin --release
cd $working_dir
# This script is used to package the application for OSX
list=`otool -L $executable | grep homebrew | awk '{print $1}'`
for i in $list; do
    echo "Fixing $i"
    install_name_tool -change $i "@executable_path/"`basename $i` $executable
    cp $i `basename $i`
    chmod +x *.dylib
done

#	<key>NSCameraUsageDescription</key>
#	<string>Zugriff auf die Kamera um Stimmzettel zu fotografieren.</string>

list=`ls /opt/homebrew/Cellar/opencv/4.10.0_3/lib/*.dylib`
for i in $list; do
    echo "Fixing $i"
    install_name_tool -change $i "@executable_path/"`basename $i` $executable
    cp $i `basename $i`
    chmod +x *.dylib
done

for i in $train_data_list; do
    cp $train_data_dir$i.traineddata .
done
cd $entry_dir