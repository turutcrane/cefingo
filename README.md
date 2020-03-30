# cefingo
This is experimental go binding for CEF.

## Supported Environment
* Windows 10 64bit (msys2 environment is recomended)

## How to Start
1. Download Spotify's Autobuild Image http://opensource.spotify.com/cefbuilds/cef_binary_80.1.2%2Bg9d2a31e%2Bchromium-80.0.3987.149_windows64.tar.bz2

1. Expand it.

1. Copy library files and Resouce files to a Directory in PATH Envrironment Variable (eg. $GOPATH/bin).

    ```bat
    C:\> xcopy /e \path\to\expand_dir\Release \path\to\gopath\bin
    C:\> xcopy /e \path\to\expand_dir\Resources \path\to\gopath\bin
    ```

1. create cefingo.pc file on PKG_CONFIG_PATH

    ```.pc
    target=C:\\path\\to\\gopath\\bin
    libdir=${target}
    includedir=C:\\path\\to\\expand_dir

    Name: cefingo
    Version: 0.1
    Description: cefingo
    Cflags: -I${includedir}
    Libs: -L${libdir} -lcef
    ```

1. go get this packages.

    ```bat
    C:\> go get github.com/turutcrane/cefingo/...
    ```


## Example
  https://github.com/turutcrane/cefingo-sample

## License
This project is licensed under the MIT License.

This project quotes header files of the following third party libraries:
* Chromium Embedded Framework licensed under the BSD 3-clause
  license. Website: https://bitbucket.org/chromiumembedded/cef

Thanks to https://github.com/cztomczak/cefcapi .
