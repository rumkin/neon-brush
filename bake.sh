GODIR=.go
GOPATH=/usr/src/golang
PACKFILE=gopkg
if [ -f $PACKFILE ]; then
    OUT=`cat ${PACKFILE}`
else
    OUT=main
fi
GODEPS=godeps
GOVER=gover

function __init {
    # Project name
    if [ $# -lt 1 ]
    then
        echo -n "Enter project name: "
        read NAME
    else
        NAME=$1
    fi

    # Update FS layout
    if [ ! -f $PACKFILE ]; then
        echo $NAME > $PACKFILE;
    fi

    if [ ! -f "$GODEPS" ]; then
        touch $GODEPS
    fi

    if [ ! -d "$GODIR" ]; then
        mkdir "$GODIR"
    fi

    # Update .gitignore file
    if [ ! -f ".gitignore" ]; then
        touch .gitignore
    fi

    if [ -z "$(cat '.gitignore' | grep "$GODIR")" ]; then
        echo $GODIR >> .gitignore
    fi

    if [ -z "$(cat '.gitignore' | grep "build")" ]; then
        echo "build" >> .gitignore
    fi

    if [ ! -f "$NAME.go" ]; then
      echo 'package main' > "$NAME.go"
    fi

    if [ ! -d ".git" ]; then
      git init
    fi
};

function __build {
    build $OUT $@
}

# Build binary from sources.
# The first argument is a name of binary.
function build {
    OUT=$1
    shift 1

    docker run --rm \
		-e GOPATH=${GOPATH}:/ \
		-v $PWD/$GODIR:$GOPATH \
		-v $PWD:/src \
		-w /src \
		golang go build -o build/$OUT $@
}

# Install package from repository or file system
# Note that installation is running inside docker container
# and you should to add directory manually
function __install {
    docker run --rm \
		-e GOPATH=${GOPATH} \
		-v $PWD/$GODIR:${GOPATH} \
		-w $GOPATH/ \
		golang go get "$1" || exit 1

    if [ -z "$(cat $GODEPS | grep "$1")" ]; then
        echo $1 >> $GODEPS
    fi
}

# Remove package from list of dependencies
function __uninstall {
    SRC=$GODIR/src/$1

    if [ -d "$SRC" ]; then
        rm -rf $SRC
        rm -rf $GODIR/pkg/*/$1.a

        cat $GODEPS | grep -v "$1" > $GODEPS
    fi
}

# Build and run
function __run {
    __build $OUT $@
    ./build/$OUT
}

# Clean build directory
function __clean {
    rm -rf build/*
}
