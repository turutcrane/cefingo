# cefingo
This is experimental go binding for CEF.

## Supported Environment
* Windows 10 64bit (msys2 environment is recomended)

## How to Start
1. Download Spotify's Autobuild Image http://opensource.spotify.com/cefbuilds/cef_binary_3.3497.1836.gb472a8d_windows64.tar.bz2

1. Expand it.

1. Copy library files and Resouce files to a Directory in PATH Envrironment Variables (eg. $GOPATH/bin).

    ```cmd
    C:\> xcopy /e \path\to\expand_dir\Release \path\to\gopath\bin
    C:\> xcopy /e \path\to\expand_dir\Resources \path\to\gopath\bin
    ```

1. Setup CGO Environment values.

    ```cmd
    CGO_LDFLAGS=-L\path\to\gopath\bin -lcef
    CGO_CFLAGS=-I\path\to\expand_dir
    ```

1. go get this package.

    ```cmd
    C:\> go get github.com/turutcrane/cefingo
    ```


## Example
  https://github.com/turutcrane/cefingo-sample

## License
This project is licensed under the MIT License.

But this project includes the following third party libraries:
* Chromium Embedded Framework licensed under the BSD 3-clause
  license. Website: https://bitbucket.org/chromiumembedded/cef

Thanks to https://github.com/cztomczak/cefcapi .
