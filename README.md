# cefingo
This is experimental go binding for CEF.

## Supported Environment
* Windows 10 64bit
* msys2/mingw64

## How to Use
1. Download Spotify's Autobuild Image (windos 64bit)

    https://cef-builds.spotifycdn.com/index.html#windows64

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

1. go install gencefingo . (in msys2 mingw64 environment)

    ```
    $ go install github.com/turutcrane/gencefingo@latest
    ```

1. make cefingo dir and go mod init (in msys2 mingw64 environment)

    ```
    $ cd cefingo
    $ go mod init gtihub.com/turutcrane/cefingo
    ```


1. generate cefingo package (in msys2 mingw64 environment)

   ```
   $ cd cefingo-dir
   $ gencefingo -pkgdir .
   ```
## Example
  https://github.com/turutcrane/cefingo-sample

## Caution

Some functions and methods of cef has thread constraint. Any functions and methods in this package are not disable goroutine preemption. So thread error may be produced.
## License
This project is licensed under the MIT License.

This project quotes header files of the following third party libraries:
* Chromium Embedded Framework licensed under the BSD 3-clause
  license. Website: https://bitbucket.org/chromiumembedded/cef

Thanks to https://github.com/cztomczak/cefcapi .
