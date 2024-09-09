# io.tualo.bp


````

    fyne package -os darwin --release 

    otool -L io.tualo.bp | grep homebrew | awk '{print $1}'

````