#!/bin/bash
# usage> $ source setup_msys.sh cef_expaned_dir
#
# local here=$(cd $(dirname $BASH_SOURCE) ; pwd)

gopath=($(go env GOPATH| tr -s ';' ' '))
target=$(cygpath -u ${gopath[0]})/bin

copy_cef_binary() {
        echo Copy Target Directory: $target
        cp -r $(cygpath -u $CEF_PACKAGE)/Release/* $target
        cp -r $(cygpath -u $CEF_PACKAGE)/Resources/* $target
}

CURRENT_PACKAGE=CURRENT_PACKAGE.txt
if [ -z "$1" ] ; then
  if [ -f $CURRENT_PACKAGE ] ; then
    CEF_PACKAGE=$(cat $CURRENT_PACKAGE)
  else
    echo Usage: source $BASH_SOURCE CEF_PACKAGE_DIR
    return 1
  fi
else
  CEF_PACKAGE=$(cd $1; pwd)
  echo $CEF_PACKAGE > $CURRENT_PACKAGE
  copy_cef_binary
fi

export CEF_BINARY="$CEF_PACKAGE"
export CGO_CFLAGS="-I$(cygpath -w $CEF_BINARY)"
export CGO_LDFLAGS="-L$(cygpath -w "$target") -lcef"

export GODEBUG=cgocheck=2

echo Build With: $CEF_BINARY
